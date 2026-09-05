package users_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, logger)

	userID, err := strconv.Atoi(core_utils.GetPathParam(r, UserIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%v: %w", err, core_errors.ErrBadRequest),
			"Bad Request",
		)
		return
	}

	response, err := h.usersService.GetUser(
		ctx,
		userID,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	userDTO := userDTOFromDomain(response)

	responseHandler.JSONResponse(http.StatusOK, userDTO)
}
