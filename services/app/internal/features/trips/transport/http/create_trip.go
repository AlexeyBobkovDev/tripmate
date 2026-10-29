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
	core_utils "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/utils"
)

type CreateTripRequest struct {
	Name string `json:"name"           validate:"required,min=1,max=100"       example:"Walk around Moscow"`

	PlanningArea  string          `json:"planning_area"  validate:"required,min=1,max=100"       example:"Moscow"`
	StartLocation domain.Location `json:"start_location" validate:"required"`

	StartDate string `json:"start_date"     validate:"required,datetime=2006-01-02" example:"2026-11-14"`
	EndDate   string `json:"end_date"       validate:"required,datetime=2006-01-02" example:"2026-11-16"`

	Budget    domain.Money      `json:"budget"         validate:"required"`
	Interests []domain.Interest `json:"interests"      validate:"required,min=1"`
	TimeZone  string            `json:"time_zone"      validate:"required"                     example:"Europe/Moscow"`
}

func (request *CreateTripRequest) Validate() error {
	if _, err := time.Parse(core_consts.Layout, request.StartDate); err != nil {
		return core_errors.ErrInvalidTripStartDate
	}
	if _, err := time.Parse(core_consts.Layout, request.EndDate); err != nil {
		return core_errors.ErrInvalidTripEndDate
	}

	return nil
}

// CreateTrip godoc
//
//	@Summary		Create Trip
//	@Description	Creates a new trip
//	@ID				create-trip
//	@Tags			trips
//	@Accept			json
//	@Produce		json
//	@Param			ownerID	path		int									true	"OwnerID ID"
//	@Param			request	body		CreateTripRequest					true	"Trip data"
//	@Success		201		{object}	TripResponse						"Created trip"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Invalid trip data"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/trips/{ownerID} [post]
func (h *TripsHTTPHandler) CreateTrip(rw http.ResponseWriter, r *http.Request) {
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

	var request CreateTripRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed decode request",
		)
		return
	}

	tripDomain, err := createTripDomainFromDTO(request)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to parse birthdate",
		)
		return
	}

	response, err := h.tripsService.CreateTrip(ctx, ownerID, tripDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create trip",
		)
		return
	}

	tripResponse := tripDTOFromDomain(response)
	responseHandler.JSONResponse(
		http.StatusCreated,
		tripResponse,
	)
}

func createTripDomainFromDTO(dto CreateTripRequest) (*domain.TripCreate, error) {
	startDate, err := time.Parse(core_consts.Layout, dto.StartDate)
	if err != nil {
		return nil, core_errors.ErrInvalidTripStartDate
	}
	endDate, err := time.Parse(core_consts.Layout, dto.EndDate)
	if err != nil {
		return nil, core_errors.ErrInvalidTripEndDate
	}
	return domain.NewTripCreate(
		dto.Name,
		dto.PlanningArea,
		dto.StartLocation,
		startDate,
		endDate,
		dto.Budget,
		dto.Interests,
		dto.TimeZone,
	), nil
}
