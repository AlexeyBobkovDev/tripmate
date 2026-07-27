package users_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Name        string `json:"name"           validate:"required,min=1,max=100"`
	Surname     string `json:"surname"        validate:"required,min=1,max=100"`
	Username    string `json:"username"       validate:"required,min=1,max=100"`
	Description string `json:"description"    validate:"omitempty,min=1,max=1000"`
	BirthDate   string `json:"birth_date"     validate:"required,datetime=2006-01-02"`
	Email       string `json:"email"          validate:"required,email"`
	PhoneNumber string `json:"phone_number"   validate:"required,e164"`
	Password    string `json:"password"       validate:"required,min=8,max=100"`
}

type CreateUserResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Surname     string     `json:"surname"`
	Username    string     `json:"username"`
	Description string     `json:"description"`
	BirthDate   string     `json:"birth_date"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, logger)
	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed decode request",
		)
	}

	userDomain, err := domainFromDTO(request)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid data",
		)
	}
	response, err := h.usersService.CreateUser(ctx, userDomain, []byte(request.Password))
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create user",
		)
	}

	userResponse := CreateUserResponse(dtoFromDomain(response))
	responseHandler.JSONResponse(
		http.StatusOK,
		userResponse,
	)
}

const layout = "2006-01-02"

func domainFromDTO(dto CreateUserRequest) (*domain.User, error) {
	birthDate, err := time.Parse(layout, dto.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("parse birthdate: %w", err)
	}
	return domain.NewUserUninitialized(
		dto.Name,
		dto.Surname,
		dto.Username,
		birthDate,
		dto.Description,
		dto.Email,
		dto.PhoneNumber,
	), nil
}

func dtoFromDomain(user *domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Username:    user.Username,
		Description: user.Description,
		BirthDate:   user.BirthDate.Format(layout),
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}
