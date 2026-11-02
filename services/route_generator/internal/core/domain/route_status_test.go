package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteStatus_Constants(t *testing.T) {
	tests := []struct {
		name        string
		routeStatus RouteStatus
		want        string
	}{
		{name: "generated status", routeStatus: RouteStatusGenerated, want: "generated"},
		{name: "edited status", routeStatus: RouteStatusEdited, want: "edited"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.routeStatus))
		})
	}
}
