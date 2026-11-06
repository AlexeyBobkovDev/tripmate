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

// CreateRoute godoc
//
//	@Summary		Create Route
//	@Description	Creates a new route
//	@ID				create-route
//	@Tags			routes
//	@Accept			json
//	@Produce		json
//	@Param			ownerID	path		int									true	"Owner ID"
//	@Param			tripID	path		int									true	"Trip ID"
//	@Success		201		{object}	RouteResponse						"Created route"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid owner/trip ID"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/routes/{ownerID}/{tripID} [post]
func (h *RoutesHTTPHandler) CreateRoute(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(rw, logger)

	ownerID, err := strconv.Atoi(core_utils.GetPathParam(r, OwnerIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"Bad Request",
		)
		return
	}
	if ownerID <= 0 {
		responseHandler.ErrorResponse(
			core_errors.ErrBadRequest,
			"ownerID must be greater than zero",
		)
		return
	}

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

	response, err := h.routesService.CreateRoute(ctx, ownerID, tripID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create route",
		)
		return
	}

	routeResponse := routeDTOFromDomain(response)
	responseHandler.JSONResponse(
		http.StatusCreated,
		routeResponse,
	)
}
