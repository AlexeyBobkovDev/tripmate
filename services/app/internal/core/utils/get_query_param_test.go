package core_utils_test

import (
	"net/http"
	"testing"

	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

func TestGetQueryParam(t *testing.T) {
	params := []struct {
		name     string
		expected string
	}{
		{
			name:     "name",
			expected: "Alexey",
		},
		{
			name:     "age",
			expected: "15",
		},
		{
			name:     "city",
			expected: "Dunkirk",
		},
		{
			name:     "empty",
			expected: "",
		},
		{
			name:     "missing",
			expected: "",
		},
	}

	for _, param := range params {
		r := newRequestWithQueryParam(param.name, param.expected)
		value := core_utils.GetQueryParam(r, param.name)
		if value != param.expected {
			t.Errorf("expected(%s) != value(%s)", param.expected, value)
		}
	}
}

func newRequestWithQueryParam(name string, value string) *http.Request {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com",
		nil,
	)
	if err != nil {
		panic("failed to create new request")
	}

	query := r.URL.Query()
	query.Set(name, value)
	r.URL.RawQuery = query.Encode()

	return r
}

func BenchmarkGetQueryParam_OneParam(b *testing.B) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com?name=Alexey",
		nil,
	)
	if err != nil {
		panic("failed to create new request")
	}

	b.ResetTimer()

	for b.Loop() {
		core_utils.GetQueryParam(r, "name")
	}
}

func BenchmarkGetQueryParam_ManyParams(b *testing.B) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com?name=Alexey&age=15&city=Dunkirk&empty=&query=param",
		nil,
	)
	if err != nil {
		panic("failed to create new request")
	}

	b.ResetTimer()

	for b.Loop() {
		core_utils.GetQueryParam(r, "name")
	}
}
