package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestLocation_Validate(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		wantErr  error
	}{
		{
			name: "valid location",
			location: Location{
				Name:      "Nizhny Novgorod",
				Latitude:  56.2965,
				Longitude: 43.9361,
			},
		},
		{
			name: "valid minimum latitude",
			location: Location{
				Name:      "Location",
				Latitude:  -90,
				Longitude: 0,
			},
		},
		{
			name: "valid maximum latitude",
			location: Location{
				Name:      "Location",
				Latitude:  90,
				Longitude: 0,
			},
		},
		{
			name: "valid minimum longitude",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: -180,
			},
		},
		{
			name: "valid maximum longitude",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: 180,
			},
		},
		{
			name: "valid zero coordinates",
			location: Location{
				Name:      "Null Island",
				Latitude:  0,
				Longitude: 0,
			},
		},
		{
			name: "valid coordinates near upper boundaries",
			location: Location{
				Name:      "Location",
				Latitude:  89.999999,
				Longitude: 179.999999,
			},
		},
		{
			name: "name with surrounding whitespace",
			location: Location{
				Name:      "  Nizhny Novgorod  ",
				Latitude:  56.3,
				Longitude: 43.9,
			},
		},
		{
			name: "empty name",
			location: Location{
				Name:      "",
				Latitude:  56.3,
				Longitude: 43.9,
			},
			wantErr: core_errors.ErrInvalidLocationName,
		},
		{
			name: "name containing only spaces",
			location: Location{
				Name:      "     ",
				Latitude:  56.3,
				Longitude: 43.9,
			},
			wantErr: core_errors.ErrInvalidLocationName,
		},
		{
			name: "name containing tabs and newlines",
			location: Location{
				Name:      "\t\n\r ",
				Latitude:  56.3,
				Longitude: 43.9,
			},
			wantErr: core_errors.ErrInvalidLocationName,
		},
		{
			name: "name containing Unicode whitespace",
			location: Location{
				Name:      "\u00A0\u2003",
				Latitude:  56.3,
				Longitude: 43.9,
			},
			wantErr: core_errors.ErrInvalidLocationName,
		},
		{
			name: "latitude below minimum",
			location: Location{
				Name:      "Location",
				Latitude:  -90.000001,
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude above maximum",
			location: Location{
				Name:      "Location",
				Latitude:  90.000001,
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude negative infinity",
			location: Location{
				Name:      "Location",
				Latitude:  math.Inf(-1),
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude positive infinity",
			location: Location{
				Name:      "Location",
				Latitude:  math.Inf(1),
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "latitude NaN",
			location: Location{
				Name:      "Location",
				Latitude:  math.NaN(),
				Longitude: 0,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "longitude below minimum",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: -180.000001,
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude above maximum",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: 180.000001,
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude negative infinity",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: math.Inf(-1),
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude positive infinity",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: math.Inf(1),
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "longitude NaN",
			location: Location{
				Name:      "Location",
				Latitude:  0,
				Longitude: math.NaN(),
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "invalid name takes precedence",
			location: Location{
				Name:      " ",
				Latitude:  100,
				Longitude: 200,
			},
			wantErr: core_errors.ErrInvalidLocationName,
		},
		{
			name: "invalid latitude takes precedence",
			location: Location{
				Name:      "Location",
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
