package routes_transport_http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http/mocks"
)

func TestRoutesHTTPHandler_Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockRoutesService(ctrl)
	handler := NewRoutesHTTPHandler(service)

	routes := handler.Routes()

	assert.Len(t, routes, 4)

	expected := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/routes/{ownerID}/{tripID}"},
		{http.MethodGet, "/routes/{tripID}/{routeID}"},
		{http.MethodGet, "/routes/{tripID}"},
		{http.MethodDelete, "/routes/{tripID}/{routeID}"},
	}

	for i, want := range expected {
		assert.Equal(t, want.method, routes[i].Method)
		assert.Equal(t, want.path, routes[i].Path)
		assert.NotNil(t, routes[i].Handler)
	}
}
