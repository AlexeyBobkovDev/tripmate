package core_utils_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

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
		{
			name:     "special-chars",
			expected: "user-123",
		},
		{
			name:     "unicode",
			expected: "测试",
		},
		{
			name:     "numbers",
			expected: "12345",
		},
		{
			name:     "mixed-chars",
			expected: "user_name123",
		},
		{
			name:     "url-encoded",
			expected: "hello%20world",
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

func TestGetQueryParam_WithMultipleParams(t *testing.T) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com?name=Alexey&age=15&city=Dunkirk",
		nil,
	)
	if err != nil {
		t.Fatal("failed to create new request")
	}

	// Test getting each parameter
	assert.Equal(t, "Alexey", core_utils.GetQueryParam(r, "name"))
	assert.Equal(t, "15", core_utils.GetQueryParam(r, "age"))
	assert.Equal(t, "Dunkirk", core_utils.GetQueryParam(r, "city"))
	assert.Equal(t, "", core_utils.GetQueryParam(r, "nonexistent"))
}

func TestGetQueryParam_WithEncodedValues(t *testing.T) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com?search=hello%20world&query=test%26example",
		nil,
	)
	if err != nil {
		t.Fatal("failed to create new request")
	}

	assert.Equal(t, "hello world", core_utils.GetQueryParam(r, "search"))
	assert.Equal(t, "test&example", core_utils.GetQueryParam(r, "query"))
}

func TestGetQueryParam_WithEmptyRequest(t *testing.T) {
	// Test with nil request
	var r *http.Request = nil
	assert.Panics(t, func() {
		core_utils.GetQueryParam(r, "any")
	}, "Expected panic for nil request")
}

func TestGetQueryParam_WithComplexQuery(t *testing.T) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com?filter=active&sort=asc&page=1&limit=20&search=test+query",
		nil,
	)
	if err != nil {
		t.Fatal("failed to create new request")
	}

	assert.Equal(t, "active", core_utils.GetQueryParam(r, "filter"))
	assert.Equal(t, "asc", core_utils.GetQueryParam(r, "sort"))
	assert.Equal(t, "1", core_utils.GetQueryParam(r, "page"))
	assert.Equal(t, "20", core_utils.GetQueryParam(r, "limit"))
	assert.Equal(t, "test query", core_utils.GetQueryParam(r, "search"))
	assert.Equal(t, "", core_utils.GetQueryParam(r, "nonexistent"))
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
