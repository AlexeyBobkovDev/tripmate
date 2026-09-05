package users_transport_http

import (
	"net/http"
	"strconv"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, logger)

	userID, err := strconv.Atoi(core_utils.GetPathParam(r, UserIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"userID must be of type int",
		)
		return
	}

	err = h.usersService.DeleteUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.JSONResponse(
		http.StatusNoContent,
		"user is successfully deleted",
	)
}
