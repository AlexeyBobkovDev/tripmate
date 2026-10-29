package users_transport_http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
)

func TestUsersHTTPHandler_Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(service)

	routes := handler.Routes()

	assert.Len(t, routes, 4)

	expected := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/users"},
		{http.MethodGet, "/users/{userID}"},
		{http.MethodDelete, "/users/{userID}"},
		{http.MethodPatch, "/users/{userID}"},
	}

	for i, want := range expected {
		assert.Equal(t, want.method, routes[i].Method)
		assert.Equal(t, want.path, routes[i].Path)
		assert.NotNil(t, routes[i].Handler)
	}
}
