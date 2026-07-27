package users_transport_http

import (
	"context"
	"net/http"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/server"
)

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user *domain.User,
		password []byte,
	) (*domain.User, error)
}

type UsersHTTPHandler struct {
	usersService UsersService
}

func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []*core_server.Route {
	return []*core_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
