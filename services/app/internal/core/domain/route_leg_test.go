package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteLeg_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*RouteLeg)
		wantErr error
	}{
		{name: "valid leg"},
		{name: "invalid leg ID", mutate: func(l *RouteLeg) { l.ID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid from stop ID", mutate: func(l *RouteLeg) { l.FromStopID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid to stop ID", mutate: func(l *RouteLeg) { l.ToStopID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "unsupported travel mode", mutate: func(l *RouteLeg) { l.TravelMode = "unknown" }, wantErr: core_errors.ErrUnsupportedRouteTravelMode},
		{name: "zero distance", mutate: func(l *RouteLeg) { l.Distance = 0 }, wantErr: core_errors.ErrInvalidRouteDistance},
		{name: "negative distance", mutate: func(l *RouteLeg) { l.Distance = -1 }, wantErr: core_errors.ErrInvalidRouteDistance},
		{name: "zero duration", mutate: func(l *RouteLeg) { l.Duration = 0 }, wantErr: core_errors.ErrInvalidRouteDuration},
		{name: "negative duration", mutate: func(l *RouteLeg) { l.Duration = -time.Minute }, wantErr: core_errors.ErrInvalidRouteDuration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leg := RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   5,
				Duration:   10 * time.Minute,
			}
			if tt.mutate != nil {
				tt.mutate(&leg)
			}
			err := leg.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
