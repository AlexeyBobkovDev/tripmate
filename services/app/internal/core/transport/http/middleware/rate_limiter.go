package core_middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Window struct {
	start     time.Time
	end       time.Time
	reqAmount int
}

func NewWindow(windowSize time.Duration) *Window {
	start := time.Now()
	return &Window{
		start:     start,
		end:       start.Add(windowSize),
		reqAmount: 0,
	}
}

type RateLimiter struct {
	prevWindow   *Window
	curWindow    *Window
	WindowSize   time.Duration
	MaxReqAmount int
}

func NewRateLimiter(
	maxReqAmount int,
	windowSize time.Duration,
) *RateLimiter {
	return &RateLimiter{
		prevWindow:   NewWindow(windowSize),
		curWindow:    NewWindow(windowSize),
		MaxReqAmount: maxReqAmount,
		WindowSize:   windowSize,
	}
}

func slidingWindowCounterAlgorithm(
	prevWindow *Window,
	curWindow *Window,
) float64 {
	windowSize := prevWindow.end.Sub(prevWindow.start).Seconds()
	elapsed := time.Since(curWindow.start).Seconds()
	weight := 1 - elapsed/windowSize

	if weight < 0 {
		weight = 0
	}

	fmt.Println("Requests since last 10s", float64(curWindow.reqAmount)+float64(prevWindow.reqAmount)*weight)
	return float64(curWindow.reqAmount) + float64(prevWindow.reqAmount)*weight
}

func (l *RateLimiter) IsAllowed() bool {
	if time.Now().After(l.curWindow.end) {
		l.Reset()
	}
	if slidingWindowCounterAlgorithm(
		l.prevWindow,
		l.curWindow,
	)+1 <= float64(l.MaxReqAmount) {
		return true
	}

	return false
}

func (l *RateLimiter) Add() {
	l.curWindow.reqAmount++
	fmt.Println("l.curWindow.reqAmount after increasing", l.curWindow.reqAmount)
}

func (l *RateLimiter) Reset() {
	l.prevWindow = l.curWindow
	l.curWindow = NewWindow(l.WindowSize)
	fmt.Printf("prev: %+v\n", *l.prevWindow)
	fmt.Printf("cur:  %+v\n", *l.curWindow)
}

type RateLimiterStore struct {
	mu           sync.RWMutex
	users        map[string]*RateLimiter
	WindowSize   time.Duration
	MaxReqAmount int
}

func (s *RateLimiterStore) addUser(user string) {
	if _, ok := s.users[user]; !ok {
		s.users[user] = NewRateLimiter(s.MaxReqAmount, s.WindowSize)
	}
}

func (s *RateLimiterStore) increase(user string) {
	s.users[user].Add()
}

func (s *RateLimiterStore) isAllowed(user string) bool {
	return s.users[user].IsAllowed()
}

func (s *RateLimiterStore) Allow(user string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.addUser(user)

	if s.isAllowed(user) {
		s.increase(user)
		return true
	}
	return false
}

func NewRateLimiterStore(
	maxReqAmount int,
	windowSize time.Duration,
) *RateLimiterStore {
	return &RateLimiterStore{
		users:        make(map[string]*RateLimiter),
		MaxReqAmount: maxReqAmount,
		WindowSize:   windowSize,
	}
}

func getUserIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func getUserIDFromToken(token string) (string, error) {
	// TODO: get user id from token
	// Ths is a stub
	return "123", nil
}

func RateLimiterMiddleware(
	maxReqAmount int,
	windowSize time.Duration,
) Middleware {
	fmt.Println(maxReqAmount, windowSize)
	rateLimiterStore := NewRateLimiterStore(
		maxReqAmount,
		windowSize,
	)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(rw, r)
			}

			var (
				rateLimitKey string
				err          error
			)
			if token := r.Header.Get("Authorization"); token != "" {
				rateLimitKey, err = getUserIDFromToken(token)
			} else {
				rateLimitKey, err = getUserIP(r)
			}
			if err != nil {
				http.Error(
					rw,
					"internal server error",
					http.StatusInternalServerError,
				)
				return
			}

			if !rateLimiterStore.Allow(rateLimitKey) {
				http.Error(
					rw,
					"Too many requests. Please wait for a while",
					http.StatusTooManyRequests,
				)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
