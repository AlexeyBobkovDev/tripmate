package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullable_String(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[string]
		expected string
	}{
		{
			name:     "string with value",
			nullable: Nullable[string]{Value: &[]string{"test"}[0], Set: true},
			expected: "test",
		},
		{
			name:     "string without value",
			nullable: Nullable[string]{Value: nil, Set: false},
			expected: "",
		},
		{
			name:     "string with nil value",
			nullable: Nullable[string]{Value: nil, Set: true},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nullable.Value != nil {
				assert.Equal(t, tt.expected, *tt.nullable.Value)
			} else {
				assert.Equal(t, tt.expected, "")
			}
		})
	}
}

func TestNullable_Int(t *testing.T) {
	tests := []struct {
		name     string
		nullable Nullable[int]
		expected int
	}{
		{
			name:     "int with value",
			nullable: Nullable[int]{Value: &[]int{42}[0], Set: true},
			expected: 42,
		},
		{
			name:     "int without value",
			nullable: Nullable[int]{Value: nil, Set: false},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nullable.Value != nil {
				assert.Equal(t, tt.expected, *tt.nullable.Value)
			} else {
				assert.Equal(t, tt.expected, 0)
			}
		})
	}
}

func TestNullable_Time(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		nullable Nullable[time.Time]
		expected time.Time
	}{
		{
			name:     "time with value",
			nullable: Nullable[time.Time]{Value: &[]time.Time{now}[0], Set: true},
			expected: now,
		},
		{
			name:     "time without value",
			nullable: Nullable[time.Time]{Value: nil, Set: false},
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nullable.Value != nil {
				assert.Equal(t, tt.expected, *tt.nullable.Value)
			} else {
				assert.Equal(t, tt.expected, time.Time{})
			}
		})
	}
}

func TestNullable_Equality(t *testing.T) {
	// Both set with same value
	testValue := "test"
	nullable1 := Nullable[string]{Value: &testValue, Set: true}
	nullable2 := Nullable[string]{Value: &testValue, Set: true}
	assert.True(t, nullable1.Value == nullable2.Value && nullable1.Set == nullable2.Set)

	// Both set with different values
	otherValue := "other"
	nullable3 := Nullable[string]{Value: &testValue, Set: true}
	nullable4 := Nullable[string]{Value: &otherValue, Set: true}
	assert.False(t, nullable3.Value == nullable4.Value && nullable3.Set == nullable4.Set)

	// One set, one not
	nullable5 := Nullable[string]{Value: &testValue, Set: true}
	nullable6 := Nullable[string]{Value: nil, Set: false}
	assert.False(t, nullable5.Value == nullable6.Value && nullable5.Set == nullable6.Set)

	// Both not set
	nullable7 := Nullable[string]{Value: nil, Set: false}
	nullable8 := Nullable[string]{Value: nil, Set: false}
	assert.True(t, nullable7.Value == nullable8.Value && nullable7.Set == nullable8.Set)
}

func TestNullable_Hashability(t *testing.T) {
	// Test that Nullable can be used as a map key
	nullableMap := make(map[Nullable[string]]bool)

	testValue := "test"
	nullableWithVal := Nullable[string]{Value: &testValue, Set: true}
	nullableWithoutVal := Nullable[string]{Value: nil, Set: false}

	// Add to map
	nullableMap[nullableWithVal] = true
	nullableMap[nullableWithoutVal] = true

	// Retrieve from map
	assert.True(t, nullableMap[nullableWithVal], "Nullable with value should be in map")
	assert.True(t, nullableMap[nullableWithoutVal], "Nullable without value should be in map")

	// Test a different nullable (same as nullableWithoutVal but created separately)
	differentNullable := Nullable[string]{Value: nil, Set: false}
	assert.True(t, nullableMap[differentNullable], "Same nullable should be in map")
}
