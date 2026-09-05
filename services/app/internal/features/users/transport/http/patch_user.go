package users_transport_http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	core_http_types "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/types"
	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

type PatchUserRequest struct {
	Version     int                              `json:"version"        validate:"min=1"                          example:"5"`
	Name        core_http_types.Nullable[string] `json:"name"                                                     example:"Name"`
	Surname     core_http_types.Nullable[string] `json:"surname"                                                  example:"Surname"`
	Username    core_http_types.Nullable[string] `json:"username"                                                 example:"Username"`
	Description core_http_types.Nullable[string] `json:"description"                                              example:"Some kind of description"`
	BirthDate   core_http_types.Nullable[string] `json:"birth_date"                                               example:"2006-01-02"`
	Email       core_http_types.Nullable[string] `json:"email"                                                    example:"checkemail@gmail.com"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"                                             example:"+79990746978"`
}

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, logger)

	userID, err := strconv.Atoi(core_utils.GetPathParam(r, UserIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"decode user patch",
		)
		return
	}

	userPatch, err := patchUserDTOToDomain(request)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"patch is not valid",
		)
		return
	}

	user, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := userDTOFromDomain(user)

	responseHandler.JSONResponse(http.StatusOK, response)
}

func patchUserDTOToDomain(request PatchUserRequest) (*domain.UserPatch, error) {
	var birthDate *time.Time
	if request.BirthDate.Set {
		birth, err := time.Parse(layout, *request.BirthDate.Value)
		birthDate = &birth
		if err != nil {
			return nil, fmt.Errorf("parse birthdate")
		}
	}
	return &domain.UserPatch{
		Version:     request.Version,
		Name:        request.Name.ToDomain(),
		Surname:     request.Surname.ToDomain(),
		Username:    request.Username.ToDomain(),
		Description: request.Description.ToDomain(),
		BirthDate:   domain.Nullable[time.Time]{Value: birthDate, Set: request.BirthDate.Set},
		Email:       request.Email.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}, nil
}
