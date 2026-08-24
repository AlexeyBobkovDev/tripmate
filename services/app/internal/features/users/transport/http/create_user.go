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
	Name        string `json:"name"           validate:"required,min=1,max=100"         example:"Name"`
	Surname     string `json:"surname"        validate:"required,min=1,max=100"         example:"Surname"`
	Username    string `json:"username"       validate:"required,min=1,max=100"         example:"Username"`
	Description string `json:"description"    validate:"omitempty,min=1,max=1000"       example:"Some kind of description"`
	BirthDate   string `json:"birth_date"     validate:"required,datetime=2006-01-02"   example:"2006-01-02"`
	Email       string `json:"email"          validate:"required,email"                 example:"checkemail@gmail.com"`
	PhoneNumber string `json:"phone_number"   validate:"required,e164"                  example:"+79990746978"`
	Password    string `json:"password"       validate:"required,min=8,max=100"         example:"some-random-password"`
}

type CreateUserResponse struct {
	ID          int        `json:"id"             example:"1"`
	Name        string     `json:"name"           example:"Name"`
	Surname     string     `json:"surname"        example:"Surname"`
	Username    string     `json:"username"       example:"Username"`
	Description string     `json:"description"    example:"Some kind of description"`
	BirthDate   string     `json:"birth_date"     example:"2006-01-02"`
	Email       string     `json:"email"          example:"checkemail@gmail.com"`
	PhoneNumber string     `json:"phone_number"   example:"+79990746978"`
	CreatedAt   time.Time  `json:"created_at"     example:"2006-01-02T15:06:07.292454Z"`
	UpdatedAt   time.Time  `json:"updated_at"     example:"2006-01-02T15:06:07.292454Z"`
	DeletedAt   *time.Time `json:"deleted_at"     example:"2006-01-02T15:06:07.292454Z"`
}

// CreateUser godoc
//
//	@Summary		Create User
//	@Description	This endpoint creates a user
//	@ID				create-user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateUserRequest					true	"The data which this endpoint will use to create user"
//	@Success		201		{object}	CreateUserResponse					"Created user"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid user data"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/users [post]
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
		http.StatusCreated,
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
