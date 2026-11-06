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

// GetRoutes godoc
//
//	@Summary		Gets Routes
//	@Description	Gets routes
//	@ID				get-routes
//	@Tags			routes
//	@Param			tripID	path		int									true	"TripID ID"
//	@Success		200		{object}	[]RouteResponse						"Successfully retrieved routes"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid TripID"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Not Found"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/routes/{tripID} [get]
func (h *RoutesHTTPHandler) GetRoutes(w http.ResponseWriter, r *http.Request) {
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

	responses, err := h.routesService.GetRoutes(
		ctx,
		tripID,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get routes",
		)
		return
	}

	routeDTOs := routeDTOsFromDomains(responses)

	responseHandler.JSONResponse(http.StatusOK, routeDTOs)
}
