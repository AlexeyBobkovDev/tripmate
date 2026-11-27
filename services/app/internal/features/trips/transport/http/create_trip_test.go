package trips_transport_http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	domain "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsTransportHTTP_CreateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	tripResponse := testTripResponse()

	mockTripsService.
		EXPECT().
		CreateTrip(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(&tripResponse, nil)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips/1",
		testTripRequestBodyWithDefaults(),
	)

	req.SetPathValue("ownerID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var trip TripResponse
	if err := json.NewDecoder(rec.Body).Decode(&trip); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testExpectedTripResponse().ID, trip.ID)
	assert.Equal(t, testExpectedTripResponse().OwnerID, trip.OwnerID)
	assert.Equal(t, testExpectedTripResponse().Version, trip.Version)
	assert.Equal(t, testExpectedTripResponse().Name, trip.Name)
	assert.Equal(t, testExpectedTripResponse().PlanningArea, trip.PlanningArea)
	assert.Equal(t, testExpectedTripResponse().StartCoordinates, trip.StartCoordinates)
	assert.True(t, testExpectedTripResponse().StartDate.Equal(trip.StartDate))
	assert.True(t, testExpectedTripResponse().EndDate.Equal(trip.EndDate))
	assert.Equal(t, testExpectedTripResponse().Budget, trip.Budget)
	assert.Equal(t, testExpectedTripResponse().Interests, trip.Interests)
	assert.Equal(t, testExpectedTripResponse().TimeZone, trip.TimeZone)
	assert.True(t, testExpectedTripResponse().CreatedAt.Equal(trip.CreatedAt))
	assert.True(t, testExpectedTripResponse().UpdatedAt.Equal(trip.UpdatedAt))
	assert.Equal(t, testExpectedTripResponse().DeletedAt, trip.DeletedAt)
}

func TestTripsTransportHTTP_CreateTrip_ServiceError(t *testing.T) {
	serviceError := core_postgres_pool_errors.ErrViolatedForeignKey

	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	mockTripsService.
		EXPECT().
		CreateTrip(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, serviceError)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips/1",
		testTripRequestBodyWithDefaults(),
	)

	req.SetPathValue("ownerID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var returnedServiceError core_http_response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&returnedServiceError); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "failed to create trip", returnedServiceError.Message)
	assert.Equal(t, core_postgres_pool_errors.ErrViolatedForeignKey.Error(), returnedServiceError.Error)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestTripsTransportHTTP_CreateTrip_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	mockTripsService.
		EXPECT().
		CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips/1",
		strings.NewReader(`invalid json`),
	)
	req.SetPathValue("ownerID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestTripsTransportHTTP_CreateTrip_ValidationError(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty name",
			body: `{
				"owner_id": 1,
				"name": "",
				"planning_area": "Moscow",
				"start_location": {"latitude": 55.7539, "longitude": 37.6208},
				"start_date": "2026-11-14",
				"end_date": "2026-11-14",
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": ["history"],
				"time_zone": "Europe/Moscow"
			}`,
		},
		{
			name: "invalid owner ID",
			body: `{
				"owner_id": 0,
				"name": "Trip",
				"planning_area": "Moscow",
				"start_location": {"latitude": 55.7539, "longitude": 37.6208},
				"start_date": "2026-11-14",
				"end_date": "2026-11-14",
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": ["history"],
				"time_zone": "Europe/Moscow"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTripsService := mocks.NewMockTripsService(ctrl)
			tripsTransport := NewTripsHTTPHandler(mockTripsService)

			mockTripsService.
				EXPECT().
				CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(
				t,
				http.MethodPost,
				"/trips",
				strings.NewReader(tt.body),
			)

			logger := testutils.SetupLogger(
				t,
				testutils.TestLevel,
				t.TempDir(),
			)

			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			tripsTransport.CreateTrip(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}

func TestTripsTransportHTTP_CreateTrip_OwnerIDNotPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	mockTripsService.
		EXPECT().
		CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	rec, req := testutils.SetupHTTP(t, http.MethodPost, "/trips/0", testTripRequestBodyWithDefaults())
	req.SetPathValue("ownerID", "0")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTripsTransportHTTP_CreateTrip_OwnerIDInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	mockTripsService.
		EXPECT().
		CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	rec, req := testutils.SetupHTTP(t, http.MethodPost, "/trips/abc", testTripRequestBodyWithDefaults())
	req.SetPathValue("ownerID", "abc")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTripRequest_Validate(t *testing.T) {
	t.Run("returns nil for valid dates", func(t *testing.T) {
		request := CreateTripRequest{
			Name:             "Walk around Moscow",
			PlanningArea:     "Moscow",
			StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        "2035-03-10",
			EndDate:          "2035-03-11",
			DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			MealRequirements: []MealRequirementRequestDTO{},
			Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []domain.Interest{"history"},
			TimeZone:         "Europe/Moscow",
		}

		require.NoError(t, request.Validate())
	})

	t.Run("returns invalid start date error when start date is invalid", func(t *testing.T) {
		request := CreateTripRequest{
			Name:             "Walk around Moscow",
			PlanningArea:     "Moscow",
			StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        "not-a-date",
			EndDate:          "2035-03-11",
			DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			MealRequirements: []MealRequirementRequestDTO{},
			Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []domain.Interest{"history"},
			TimeZone:         "Europe/Moscow",
		}

		err := request.Validate()

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripStartDate)
	})

	t.Run("returns invalid end date error when end date is invalid", func(t *testing.T) {
		request := CreateTripRequest{
			Name:             "Walk around Moscow",
			PlanningArea:     "Moscow",
			StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        "2035-03-10",
			EndDate:          "not-a-date",
			DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			MealRequirements: []MealRequirementRequestDTO{},
			Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []domain.Interest{"history"},
			TimeZone:         "Europe/Moscow",
		}

		err := request.Validate()

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndDate)
	})
}

func TestCreateTripDomainFromDTO(t *testing.T) {
	t.Run("returns nil error and valid trip when dates are valid", func(t *testing.T) {
		trip, err := createTripDomainFromDTO(CreateTripRequest{
			Name:             "Trip to Moscow",
			PlanningArea:     "Moscow",
			StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        "2035-03-10",
			EndDate:          "2035-03-11",
			DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			MealRequirements: []MealRequirementRequestDTO{},
			Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []domain.Interest{"history"},
			TimeZone:         "Europe/Moscow",
		})

		require.NoError(t, err)
		require.NotNil(t, trip)
		require.True(t, trip.StartDate.Equal(time.Date(2035, 3, 10, 0, 0, 0, 0, time.UTC)))
		require.True(t, trip.EndDate.Equal(time.Date(2035, 3, 11, 0, 0, 0, 0, time.UTC)))
	})

	t.Run("returns invalid start date error when start date is invalid", func(t *testing.T) {
		trip, err := createTripDomainFromDTO(CreateTripRequest{
			Name:             "Trip to Moscow",
			PlanningArea:     "Moscow",
			StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        "not-a-date",
			EndDate:          "2035-03-11",
			DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			MealRequirements: []MealRequirementRequestDTO{},
			Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []domain.Interest{"history"},
			TimeZone:         "Europe/Moscow",
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripStartDate)
		assert.Nil(t, trip)
	})

	t.Run("returns invalid end date error when end date is invalid", func(t *testing.T) {
		trip, err := createTripDomainFromDTO(CreateTripRequest{
			Name:             "Trip to Moscow",
			PlanningArea:     "Moscow",
			StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        "2035-03-10",
			EndDate:          "not-a-date",
			DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			MealRequirements: []MealRequirementRequestDTO{},
			Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []domain.Interest{"history"},
			TimeZone:         "Europe/Moscow",
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndDate)
		assert.Nil(t, trip)
	})
}

func TestTripsTransportHTTP_CreateTrip_InvalidTripDates(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
	}{
		{name: "invalid start date", startDate: "not-a-date", endDate: "2035-03-10"},
		{name: "invalid end date", startDate: "2035-03-10", endDate: "not-a-date"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTripsService := mocks.NewMockTripsService(ctrl)
			tripsTransport := NewTripsHTTPHandler(mockTripsService)

			mockTripsService.
				EXPECT().
				CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			body := `{
				"owner_id": 1,
				"name": "Trip to Moscow",
				"planning_area": "Moscow",
				"start_location": {"latitude": 55.7539, "longitude": 37.6208},
				"start_date": "` + tt.startDate + `",
				"end_date": "` + tt.endDate + `",
				"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
				"meal_requirements": [],
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": ["history"],
				"time_zone": "Europe/Moscow"
			}`

			rec, req := testutils.SetupHTTP(t, http.MethodPost, "/trips/1", strings.NewReader(body))
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			tripsTransport.CreateTrip(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestTripsTransportHTTP_CreateTrip_InvalidBirthdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockService)

	body := `{
		"owner_id": 1,
		"name": "Trip to Moscow",
		"planning_area": "Moscow",
		"start_location": {"latitude": 55.7539, "longitude": 37.6208},
		"start_date": "2026-12-31",
		"end_date": "2026-12-31",
		"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
		"meal_requirements": [],
		"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
		"interests": ["history"],
		"time_zone": "Europe/Moscow",
		"birthdate": "invalid-date"
	}`

	rec, req := testutils.SetupHTTP(t, http.MethodPost, "/trips/1", strings.NewReader(body))
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTripsTransportHTTP_CreateTrip_TransactionCommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockService)

	commitError := errors.New("commit error")

	mockService.EXPECT().
		CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, commitError)

	body := `{
		"owner_id": 1,
		"name": "Trip to Moscow",
		"planning_area": "Moscow",
		"start_location": {"latitude": 55.7539, "longitude": 37.6208},
		"start_date": "2026-12-31",
		"end_date": "2026-12-31",
		"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
		"meal_requirements": [],
		"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
		"interests": ["history"],
		"time_zone": "Europe/Moscow",
		"birthdate": "1990-01-01"
	}`

	rec, req := testutils.SetupHTTP(t, http.MethodPost, "/trips/1", strings.NewReader(body))
	logger := testutils.SetupLogger(
		t,
		testutils.TestLevel,
		t.TempDir(),
	)
	ctx := core_logger.ToContext(req.Context(), logger)
	req.SetPathValue("ownerID", "1")
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestTripsTransportHTTP_CreateTrip_ParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockService)

	body := `{
		"owner_id": 1,
		"name": "Trip to Moscow",
		"planning_area": "Moscow",
		"start_location": {"latitude": 55.7539, "longitude": 37.6208},
		"start_date": "2026-12-31",
		"end_date": "2026-12-31",
		"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
		"meal_requirements": [],
		"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
		"interests": ["history"],
		"time_zone": "Europe/Moscow",
		"birthdate": "invalid-date"
	}`

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips/1",
		strings.NewReader(body),
	)
	logger := testutils.SetupLogger(
		t,
		testutils.TestLevel,
		t.TempDir(),
	)
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
