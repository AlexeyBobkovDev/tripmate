package users_transport_http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
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

func (request *PatchUserRequest) Validate() error {
	if request.BirthDate.Value == nil {
		return nil
	}

	if _, err := time.Parse(core_consts.Layout, *request.BirthDate.Value); err != nil {
		return core_errors.ErrInvalidBirthDateFormat
	}

	return nil
}

// PatchUser godoc
//
//	@Summary		Patch User
//	@Description	Partially updates a user by ID
//	@Description	If a field is not sent, the field won't be changed
//	@Description	If a field is sent with the value of null, the field will be changed to null
//	@Description	If a field is sent with an ordinary value, the field will be changed to this value
//	@Description	If Birthdate is set as null, ErrBadRequest(400 status code) is returned
//	@ID				patch-user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		PatchUserRequest					true	"User data"
//	@Param			userID	path		int									true	"User ID"
//	@Success		200		{object}	UserResponse						"Patched user"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid user data"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"User not found"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/users/{userID} [patch]
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, logger)

	userID, err := strconv.Atoi(core_utils.GetPathParam(r, UserIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"failed to get user id",
		)
		return
	}
	if userID <= 0 {
		responseHandler.ErrorResponse(
			core_errors.ErrBadRequest,
			"userID must be greater than zero",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"decode user patch",
		)
		return
	}

	userPatch, err := patchUserDTOToDomain(request)
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
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
		if request.BirthDate.Value == nil {
			return nil, fmt.Errorf("birthdate cannot be nil: %w", core_errors.ErrBadRequest)
		}
		birth, err := time.Parse(core_consts.Layout, *request.BirthDate.Value)
		if err != nil {
			return nil, fmt.Errorf("parse birthdate: %w", err)
		}
		birthDate = &birth
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
