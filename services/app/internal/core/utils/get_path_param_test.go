package core_utils_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

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
	}

	for _, param := range params {
		r := newRequestWithPathParam(param.name, param.expected)
		value := core_utils.GetPathParam(r, param.name)
		if value != param.expected {
			t.Errorf("expected(%s) != value(%s)", param.expected, value)
		}
	}
}

func TestGetPathParam_NonexistentParam(t *testing.T) {
	r := newRequestWithPathParam("existing", "value")

	// Test getting a parameter that doesn't exist
	value := core_utils.GetPathParam(r, "nonexistent")
	assert.Empty(t, value, "Expected empty string for nonexistent parameter")
}

func TestGetPathParam_EmptyRequest(t *testing.T) {
	// Test with nil request
	var r *http.Request = nil
	assert.Panics(t, func() {
		core_utils.GetPathParam(r, "any")
	}, "Expected panic for nil request")
}

func TestGetPathParam_WithMultipleParams(t *testing.T) {
	r, err := http.NewRequest(
		http.MethodGet,
		"http://example.com",
		nil,
	)
	if err != nil {
		t.Fatal("failed to create new request")
	}

	// Set multiple path parameters
	r.SetPathValue("id", "123")
	r.SetPathValue("name", "test")
	r.SetPathValue("category", "general")

	// Test getting each parameter
	assert.Equal(t, "123", core_utils.GetPathParam(r, "id"))
	assert.Equal(t, "test", core_utils.GetPathParam(r, "name"))
	assert.Equal(t, "general", core_utils.GetPathParam(r, "category"))
	assert.Equal(t, "", core_utils.GetPathParam(r, "nonexistent"))
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
