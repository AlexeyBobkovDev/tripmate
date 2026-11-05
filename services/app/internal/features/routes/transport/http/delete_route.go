package routes_transport_http

import (
	"errors"
	"net/http"
	"strconv"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

// DeleteRoute godoc
//
//	@Summary		Delete Route
//	@Description	Deletes a route
//	@ID				delete-route
//	@Tags			routes
//	@Param			tripID	path	int	true	"Trip ID"
//	@Param			routeID	path	int	true	"Route ID"
//	@Success		204		"Deleted route"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid trip/route ID"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Not Found"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/routes/{tripID}/{routeID} [delete]
func (h *RoutesHTTPHandler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, logger)

	tripID, err := strconv.Atoi(core_utils.GetPathParam(r, TripIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"Bad Request",
		)
		return
	}
	if tripID <= 0 {
		responseHandler.ErrorResponse(
			core_errors.ErrBadRequest,
			"tripID must be greater than zero",
		)
		return
	}

	routeID, err := strconv.Atoi(core_utils.GetPathParam(r, RouteIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"routeID must be of type int",
		)
		return
	}
	if routeID <= 0 {
		responseHandler.ErrorResponse(
			core_errors.ErrBadRequest,
			"routeID must be greater than zero",
		)
		return
	}

	err = h.routesService.DeleteRoute(ctx, tripID, routeID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete route",
		)
		return
	}

	responseHandler.NoContent()
}
