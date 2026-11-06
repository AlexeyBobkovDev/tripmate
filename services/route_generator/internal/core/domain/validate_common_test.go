package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want error
	}{
		{name: "minimum invalid", id: 0, want: core_errors.ErrInvalidID},
		{name: "negative", id: -1, want: core_errors.ErrInvalidID},
		{name: "minimum integer", id: math.MinInt, want: core_errors.ErrInvalidID},
		{name: "one", id: 1},
		{name: "positive", id: 42},
		{name: "maximum integer", id: math.MaxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
