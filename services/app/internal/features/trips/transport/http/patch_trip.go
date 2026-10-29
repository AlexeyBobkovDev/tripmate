package trips_transport_http

import (
	"errors"
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

type PatchTripRequest struct {
	Version int `json:"version"        validate:"min=1"                  example:"5"`

	Name          core_http_types.Nullable[string]          `json:"name"              swaggertype:"string"    example:"Walk around NYC"`
	PlanningArea  core_http_types.Nullable[string]          `json:"planning_area"     swaggertype:"string"    example:"NYC"`
	StartLocation core_http_types.Nullable[domain.Location] `json:"start_location"`

	StartDate core_http_types.Nullable[string] `json:"start_date"   swaggertype:"string"   example:"2026-11-14"`
	EndDate   core_http_types.Nullable[string] `json:"end_date"     swaggertype:"string"   example:"2026-11-16"`

	Budget    core_http_types.Nullable[domain.Money]      `json:"budget"`
	Interests core_http_types.Nullable[[]domain.Interest] `json:"interests"`
	TimeZone  core_http_types.Nullable[string]            `json:"time_zone"        swaggertype:"string"        example:"Europe/Moscow"`
}

func (request *PatchTripRequest) Validate() error {
	if request.StartDate.Set {
		if request.StartDate.Value == nil {
			return core_errors.ErrInvalidTripStartDate
		}

		if _, err := time.Parse(core_consts.Layout, *request.StartDate.Value); err != nil {
			return core_errors.ErrInvalidTripStartDate
		}
	}

	if request.EndDate.Set {
		if request.EndDate.Value == nil {
			return core_errors.ErrInvalidTripEndDate
		}

		if _, err := time.Parse(core_consts.Layout, *request.EndDate.Value); err != nil {
			return core_errors.ErrInvalidTripEndDate
		}
	}

	return nil
}

// PatchTrip godoc
//
//	@Summary		Patch Trip
//	@Description	Partially updates a trip by ID
//	@Description	If a field is not sent, the field won't be changed
//	@Description	If a field is sent with the value of null, the field will be changed to null
//	@Description	If a field is sent with an ordinary value, the field will be changed to this value
//	@Description	If Birthdate is set as null, ErrBadRequest(400 status code) is returned
//	@ID				patch-trip
//	@Tags			trips
//	@Accept			json
//	@Produce		json
//	@Param			request	body		PatchTripRequest					true	"Trip data"
//	@Param			ownerID	path		int									true	"OwnerID ID"
//	@Param			tripID	path		int									true	"Trip ID"
//	@Success		200		{object}	TripResponse						"Patched trip"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid trip data"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Trip not found"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/trips/{ownerID}/{tripID} [patch]
func (h *TripsHTTPHandler) PatchTrip(w http.ResponseWriter, r *http.Request) {
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
			"failed to get trip id",
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

	var request PatchTripRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"decode trip patch",
		)
		return
	}

	tripPatch, err := patchTripDomainFromDTO(ownerID, request)
	if err != nil {
		responseHandler.ErrorResponse(
			errors.Join(err, core_errors.ErrBadRequest),
			"patch is not valid",
		)
		return
	}

	trip, err := h.tripsService.PatchTrip(ctx, ownerID, tripID, tripPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch trip",
		)
		return
	}

	response := tripDTOFromDomain(trip)

	responseHandler.JSONResponse(http.StatusOK, response)
}

func patchTripDomainFromDTO(
	ownerID int,
	dto PatchTripRequest,
) (*domain.TripPatch, error) {
	var startDate domain.Nullable[time.Time]
	if dto.StartDate.Set {
		if dto.StartDate.Value == nil {
			return nil, core_errors.ErrInvalidTripStartDate
		}

		parsed, err := time.Parse(
			core_consts.Layout,
			*dto.StartDate.Value,
		)
		if err != nil {
			return nil, core_errors.ErrInvalidTripStartDate
		}

		startDate = domain.Nullable[time.Time]{
			Value: &parsed,
			Set:   true,
		}
	}

	var endDate domain.Nullable[time.Time]
	if dto.EndDate.Set {
		if dto.EndDate.Value == nil {
			return nil, core_errors.ErrInvalidTripEndDate
		}

		parsed, err := time.Parse(
			core_consts.Layout,
			*dto.EndDate.Value,
		)
		if err != nil {
			return nil, core_errors.ErrInvalidTripEndDate
		}

		endDate = domain.Nullable[time.Time]{
			Value: &parsed,
			Set:   true,
		}
	}

	return domain.NewTripPatch(
		ownerID,
		dto.Version,
		dto.Name.ToDomain(),
		dto.PlanningArea.ToDomain(),
		dto.StartLocation.ToDomain(),
		startDate,
		endDate,
		dto.Budget.ToDomain(),
		dto.Interests.ToDomain(),
		dto.TimeZone.ToDomain(),
	), nil
}
