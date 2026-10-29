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

// DeleteTrip godoc
//
//	@Summary		Delete Trip
//	@Description	Deletes a trip
//	@ID				delete-trip
//	@Tags			trips
//	@Param			ownerID	path	int	true	"OwnerID ID"
//	@Param			tripID	path	int	true	"Trip ID"
//	@Success		204		"Deleted trip"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid trip ID"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Not Found"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/trips/{ownerID}/{tripID} [delete]
func (h *TripsHTTPHandler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
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

	tripID, err := strconv.Atoi(core_utils.GetPathParam(r, TripIDParam))
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"tripID must be of type int",
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

	err = h.tripsService.DeleteTrip(ctx, ownerID, tripID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete trip",
		)
		return
	}

	responseHandler.NoContent()
}
