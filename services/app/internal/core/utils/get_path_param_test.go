package core_utils_test

import (
	"net/http"
	"testing"

	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

func TestGetPathParam(t *testing.T) {
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
		r := newRequestWithPathParam(param.name, param.expected)
		value := core_utils.GetPathParam(r, param.name)
		if value != param.expected {
			t.Errorf("expected(%s) != value(%s)", param.expected, value)
		}
	}
}

func newRequestWithPathParam(name string, value string) *http.Request {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com",
		nil,
	)
	if err != nil {
		panic("failed to create new request")
	}

	r.SetPathValue(name, value)

	return r
}

func BenchmarkGetPathParam_OneParam(b *testing.B) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com",
		nil,
	)
	if err != nil {
		panic("failed to create new request")
	}
	r.SetPathValue("a", "1")

	b.ResetTimer()

	for b.Loop() {
		core_utils.GetPathParam(r, "name")
	}
}

func BenchmarkGetPathParam_ManyParams(b *testing.B) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com",
		nil,
	)
	if err != nil {
		panic("failed to create new request")
	}
	r.SetPathValue("a", "1")
	r.SetPathValue("b", "2")
	r.SetPathValue("c", "3")
	r.SetPathValue("d", "4")

	b.ResetTimer()

	for b.Loop() {
		core_utils.GetPathParam(r, "name")
	}
}
