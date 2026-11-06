package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestCoordinates_Validate(t *testing.T) {
	tests := []struct {
		name     string
		location Coordinates
		wantErr  error
	}{
		{
			name: "valid location",
			location: Coordinates{
				Latitude:  56.2965,
				Longitude: 43.9361,
			},
		},
		{
			name: "valid minimum latitude",
			location: Coordinates{
				Latitude:  -90,
				Longitude: 0,
			},
		},
		{
			name: "valid maximum latitude",
			location: Coordinates{
				Latitude:  90,
				Longitude: 0,
			},
		},
		{
			name: "valid minimum longitude",
			location: Coordinates{
				Latitude:  0,
				Longitude: -180,
			},
		},
		{
			name: "valid maximum longitude",
			location: Coordinates{
				Latitude:  0,
				Longitude: 180,
			},
		},
		{
			name: "valid zero coordinates",
			location: Coordinates{
				Latitude:  0,
				Longitude: 0,
			},
		},
		{
			name: "valid coordinates near upper boundaries",
			location: Coordinates{
				Latitude:  89.999999,
				Longitude: 179.999999,
			},
		},
		{
			name: "name with surrounding whitespace",
			location: Coordinates{
				Latitude:  56.3,
				Longitude: 43.9,
			},
		},
		{
			name: "latitude below minimum",
			location: Coordinates{
				Latitude:  -90.000001,
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude above maximum",
			location: Coordinates{
				Latitude:  90.000001,
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude negative infinity",
			location: Coordinates{
				Latitude:  math.Inf(-1),
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude positive infinity",
			location: Coordinates{
				Latitude:  math.Inf(1),
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude NaN",
			location: Coordinates{
				Latitude:  math.NaN(),
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "longitude below minimum",
			location: Coordinates{
				Latitude:  0,
				Longitude: -180.000001,
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude above maximum",
			location: Coordinates{
				Latitude:  0,
				Longitude: 180.000001,
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude negative infinity",
			location: Coordinates{
				Latitude:  0,
				Longitude: math.Inf(-1),
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude positive infinity",
			location: Coordinates{
				Latitude:  0,
				Longitude: math.Inf(1),
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude NaN",
			location: Coordinates{
				Latitude:  0,
				Longitude: math.NaN(),
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "invalid latitude takes precedence",
			location: Coordinates{
				Latitude:  100,
				Longitude: 200,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.location.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
