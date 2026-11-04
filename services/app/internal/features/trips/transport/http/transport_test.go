package trips_transport_http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
)

func TestTripsHTTPHandler_Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(service)

	routes := handler.Routes()

	assert.Len(t, routes, 5)

	expected := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/trips/{ownerID}"},
		{http.MethodGet, "/trips/{ownerID}/{tripID}"},
		{http.MethodGet, "/trips/{ownerID}"},
		{http.MethodDelete, "/trips/{ownerID}/{tripID}"},
		{http.MethodPatch, "/trips/{ownerID}/{tripID}"},
	}

	for i, want := range expected {
		assert.Equal(t, want.method, routes[i].Method)
		assert.Equal(t, want.path, routes[i].Path)
		assert.NotNil(t, routes[i].Handler)
	}
}
