package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTravelMode_Constants(t *testing.T) {
	tests := []struct {
		name       string
		travelMode TravelMode
		want       string
	}{
		{name: "walking", travelMode: TravelModeWalking, want: "walking"},
		{name: "cycling", travelMode: TravelModeCycling, want: "cycling"},
		{name: "driving", travelMode: TravelModeDriving, want: "driving"},
		{name: "transit", travelMode: TravelModeTransit, want: "transit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.travelMode))
		})
	}
}
