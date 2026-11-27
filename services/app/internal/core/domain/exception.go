package domain

import "time"

type Exception struct {
	ExceptionDate time.Time
	OpeningHours
}
