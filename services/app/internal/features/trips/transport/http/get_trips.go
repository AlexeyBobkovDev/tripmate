package trips_transport_http

import (
	"errors"
	"net/http"
	"strconv"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

// GetTrip godoc
//
//	@Summary		Gets Trip
//	@Description	Gets trips
//	@ID				get-trips
//	@Tags			trips
//	@Param			ownerID	path		int									true	"OwnerID ID"
//	@Success		200		{object}	TripResponse						"Successfully retrieved trip"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid trip ID"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Not Found"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/trips/{ownerID} [get]
func (h *TripsHTTPHandler) GetTrips(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(w, logger)

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

	responses, err := h.tripsService.GetTrips(
		ctx,
		ownerID,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get trip",
		)
		return
	}

	tripDTOs := tripDTOsFromDomains(responses)

	responseHandler.JSONResponse(http.StatusOK, tripDTOs)
}
