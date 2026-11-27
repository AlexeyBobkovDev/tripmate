package routes_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestRoutesRepository_CreateRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	tripVersion := 1
	status := domain.RouteStatus(domain.RouteStatusGenerated)

	expectedRoute := &domain.Route{
		ID:          1,
		TripID:      1,
		TripVersion: 1,
		Version:     1,
		Status:      status,
		Days:        nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.ID
			*destinations[1].(*int) = expectedRoute.TripID
			*destinations[2].(*int) = expectedRoute.TripVersion
			*destinations[3].(*int) = expectedRoute.Version
			*destinations[4].(*domain.RouteStatus) = expectedRoute.Status
			*destinations[5].(*time.Time) = expectedRoute.CreatedAt
			*destinations[6].(*time.Time) = expectedRoute.UpdatedAt
			return nil
		})

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteQuery,
			tripID,
			tripVersion,
			status,
		).
		Return(mockRow)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	createdRoute, err := routesRepository.CreateRoute(
		context.Background(),
		1,
		&domain.Trip{ID: tripID, Version: tripVersion},
		&domain.RouteCreate{Status: status},
	)
	require.NoError(t, err)
	assert.Equal(t, expectedRoute, createdRoute)
}

func TestRoutesRepository_CreateRoute_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	tripVersion := 1
	status := domain.RouteStatus(domain.RouteStatusGenerated)
	scanError := errors.New("scan error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(scanError)

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteQuery,
			tripID,
			tripVersion,
			status,
		).
		Return(mockRow)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	createdRoute, err := routesRepository.CreateRoute(
		context.Background(),
		1,
		&domain.Trip{ID: tripID, Version: tripVersion},
		&domain.RouteCreate{Status: status},
	)
	require.ErrorIs(t, err, scanError)
	assert.Nil(t, createdRoute)
}

func TestRoutesRepository_CreateRoute_RouteDayError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	tripVersion := 1
	status := domain.RouteStatus(domain.RouteStatusGenerated)
	routeDayError := errors.New("route day error")

	// First, the route is created successfully
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1
			*destinations[1].(*int) = tripID
			*destinations[2].(*int) = tripVersion
			*destinations[3].(*int) = 1
			*destinations[4].(*domain.RouteStatus) = status
			*destinations[5].(*time.Time) = time.Now()
			*destinations[6].(*time.Time) = time.Now()
			return nil
		})

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteQuery,
			tripID,
			tripVersion,
			status,
		).
		Return(mockRow)

	// Then, when creating route day, it fails
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteDayQuery,
			gomock.Any(), // routeResponse.ID
			gomock.Any(), // day.Date
			gomock.Any(), // day.DayNumber
		).
		Return(mockRouteDayRow)

	mockRouteDayRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(routeDayError)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	createdRoute, err := routesRepository.CreateRoute(
		context.Background(),
		1,
		&domain.Trip{ID: tripID, Version: tripVersion},
		&domain.RouteCreate{
			Status: status,
			Days: []domain.RouteDay{
				{
					Date:      time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
					DayNumber: 1,
					Stops:     []domain.RouteStop{},
					Legs:      []domain.RouteLeg{},
				},
			},
		},
	)
	require.ErrorIs(t, err, routeDayError)
	assert.Nil(t, createdRoute)
}

func TestRoutesRepository_CreateRoute_RouteStopError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRow := mocks.NewMockRow(ctrl)
	mockRouteStopRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	tripVersion := 1
	status := domain.RouteStatus(domain.RouteStatusGenerated)
	routeStopError := errors.New("route stop error")

	// First, the route is created successfully
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1
			*destinations[1].(*int) = tripID
			*destinations[2].(*int) = tripVersion
			*destinations[3].(*int) = 1
			*destinations[4].(*domain.RouteStatus) = status
			*destinations[5].(*time.Time) = time.Now()
			*destinations[6].(*time.Time) = time.Now()
			return nil
		})

	// Then, the route day is created successfully
	mockRouteDayRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1 // route day ID
			*destinations[1].(*time.Time) = time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC)
			*destinations[2].(*int) = 1 // day number
			return nil
		})

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteQuery,
			tripID,
			tripVersion,
			status,
		).
		Return(mockRow)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteDayQuery,
			gomock.Any(), // routeResponse.ID
			gomock.Any(), // day.Date
			gomock.Any(), // day.DayNumber
		).
		Return(mockRouteDayRow)

	// Then, when creating route stop, it fails
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteStopQuery,
			gomock.Any(), // routeDayID
			gomock.Any(), // stop.PlaceID
			gomock.Any(), // stop.Position
			gomock.Any(), // stop.Name
			gomock.Any(), // stop.Description
			gomock.Any(), // stop.Categories
			gomock.Any(), // stop.Interests
			gomock.Any(), // stop.Address
			gomock.Any(), // stop.WebsiteURL
			gomock.Any(), // stop.Phone
			gomock.Any(), // stop.VisitDurationMinutes
			gomock.Any(), // stop.Coordinates.Latitude
			gomock.Any(), // stop.Coordinates.Longitude
			gomock.Any(), // stop.ArrivalTime
			gomock.Any(), // stop.DepartureTime
		).
		Return(mockRouteStopRow)

	mockRouteStopRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(routeStopError)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	createdRoute, err := routesRepository.CreateRoute(
		context.Background(),
		1,
		&domain.Trip{ID: tripID, Version: tripVersion},
		&domain.RouteCreate{
			Status: status,
			Days: []domain.RouteDay{
				{
					Date:      time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
					DayNumber: 1,
					Stops: []domain.RouteStop{
						{
							PlaceID:       1,
							Position:      1,
							Name:          "Test Stop",
							Description:   "Test Description",
							Coordinates:   domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
							ArrivalTime:   time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
						},
					},
					Legs: []domain.RouteLeg{},
				},
			},
		},
	)
	require.ErrorIs(t, err, routeStopError)
	assert.Nil(t, createdRoute)
}

func TestRoutesRepository_CreateRoute_RouteLegError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRow := mocks.NewMockRow(ctrl)
	mockRouteStopRow := mocks.NewMockRow(ctrl)
	mockRouteLegRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	tripVersion := 1
	status := domain.RouteStatus(domain.RouteStatusGenerated)
	routeLegError := errors.New("route leg error")

	// First, the route is created successfully
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1
			*destinations[1].(*int) = tripID
			*destinations[2].(*int) = tripVersion
			*destinations[3].(*int) = 1
			*destinations[4].(*domain.RouteStatus) = status
			*destinations[5].(*time.Time) = time.Now()
			*destinations[6].(*time.Time) = time.Now()
			return nil
		})

	// Then, the route day is created successfully
	mockRouteDayRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1 // route day ID
			*destinations[1].(*time.Time) = time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC)
			*destinations[2].(*int) = 1 // day number
			return nil
		})

	// Then, the route stop is created successfully
	mockRouteStopRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1 // stop ID
			*destinations[1].(*int) = 1 // place ID
			*destinations[2].(*int) = 1 // position
			*destinations[3].(*string) = "Test Stop"
			*destinations[4].(*string) = "Test Description"
			*destinations[5].(*[]domain.Category) = []domain.Category{}
			*destinations[6].(*[]domain.Interest) = []domain.Interest{}
			*destinations[7].(*string) = "Test Address"
			*destinations[8].(**string) = nil
			*destinations[9].(**string) = nil
			*destinations[10].(*int) = 60          // visit duration in minutes
			*destinations[11].(*float64) = 55.7539 // latitude
			*destinations[12].(*float64) = 37.6208 // longitude
			*destinations[13].(*time.Time) = time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC)
			*destinations[14].(*time.Time) = time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC)
			return nil
		})

	// Then, when creating route leg, it fails
	mockRouteLegRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(routeLegError)

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteQuery,
			tripID,
			tripVersion,
			status,
		).
		Return(mockRow)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteDayQuery,
			gomock.Any(), // routeResponse.ID
			gomock.Any(), // day.Date
			gomock.Any(), // day.DayNumber
		).
		Return(mockRouteDayRow)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteStopQuery,
			gomock.Any(), // routeDayID
			gomock.Any(), // stop.PlaceID
			gomock.Any(), // stop.Position
			gomock.Any(), // stop.Name
			gomock.Any(), // stop.Description
			gomock.Any(), // stop.Categories
			gomock.Any(), // stop.Interests
			gomock.Any(), // stop.Address
			gomock.Any(), // stop.WebsiteURL
			gomock.Any(), // stop.Phone
			gomock.Any(), // stop.VisitDurationMinutes
			gomock.Any(), // stop.Coordinates.Latitude
			gomock.Any(), // stop.Coordinates.Longitude
			gomock.Any(), // stop.ArrivalTime
			gomock.Any(), // stop.DepartureTime
		).
		Return(mockRouteStopRow)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteLegQuery,
			gomock.Any(), // routeDayID
			gomock.Any(), // mapStopIDInTripToStopIDInDB[leg.FromStopID]
			gomock.Any(), // mapStopIDInTripToStopIDInDB[leg.ToStopID]
			gomock.Any(), // leg.TravelMode
			gomock.Any(), // leg.Distance
			gomock.Any(), // durationToSeconds(leg.Duration)
		).
		Return(mockRouteLegRow)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	createdRoute, err := routesRepository.CreateRoute(
		context.Background(),
		1,
		&domain.Trip{ID: tripID, Version: tripVersion},
		&domain.RouteCreate{
			Status: status,
			Days: []domain.RouteDay{
				{
					Date:      time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
					DayNumber: 1,
					Stops: []domain.RouteStop{
						{
							PlaceID:       1,
							Position:      1,
							Name:          "Test Stop",
							Description:   "Test Description",
							Coordinates:   domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
							ArrivalTime:   time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
						},
					},
					Legs: []domain.RouteLeg{
						{
							FromStopID: 1,
							ToStopID:   2,
							TravelMode: domain.TravelModeWalking,
							Distance:   1000,
							Duration:   time.Minute,
						},
					},
				},
			},
		},
	)
	require.ErrorIs(t, err, routeLegError)
	assert.Nil(t, createdRoute)
}

func TestRoutesRepository_CreateRoute_WithMultipleDaysAndStops(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := make([]*mocks.MockRow, 2)
	mockRouteStopRows := make([]*mocks.MockRow, 4)
	mockRouteLegRows := make([]*mocks.MockRow, 2)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	tripVersion := 1
	status := domain.RouteStatus(domain.RouteStatusGenerated)

	// First, the route is created successfully
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 1
			*destinations[1].(*int) = tripID
			*destinations[2].(*int) = tripVersion
			*destinations[3].(*int) = 1
			*destinations[4].(*domain.RouteStatus) = status
			*destinations[5].(*time.Time) = time.Now()
			*destinations[6].(*time.Time) = time.Now()
			return nil
		})

	// Mock route day rows for 2 days
	for i := 0; i < 2; i++ {
		mockRouteDayRows[i] = mocks.NewMockRow(ctrl)
		mockRouteDayRows[i].EXPECT().
			Scan(
				gomock.Any(), gomock.Any(), gomock.Any(),
			).
			DoAndReturn(func(destinations ...any) error {
				*destinations[0].(*int) = i + 1 // route day ID
				*destinations[1].(*time.Time) = time.Date(2026, time.November, 14+i, 0, 0, 0, 0, time.UTC)
				*destinations[2].(*int) = i + 1 // day number
				return nil
			})
	}

	// Mock route stop rows for 4 stops (2 per day)
	for i := 0; i < 4; i++ {
		mockRouteStopRows[i] = mocks.NewMockRow(ctrl)
		mockRouteStopRows[i].EXPECT().
			Scan(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any(), gomock.Any(),
			).
			DoAndReturn(func(destinations ...any) error {
				*destinations[0].(*int) = i + 1   // stop ID
				*destinations[1].(*int) = i + 1   // place ID
				*destinations[2].(*int) = i%2 + 1 // position within the day
				*destinations[3].(*string) = fmt.Sprintf("Test Stop %d", i+1)
				*destinations[4].(*string) = fmt.Sprintf("Test Description %d", i+1)
				*destinations[5].(*[]domain.Category) = []domain.Category{}
				*destinations[6].(*[]domain.Interest) = []domain.Interest{}
				*destinations[7].(*string) = fmt.Sprintf("Test Address %d", i+1)
				*destinations[8].(**string) = nil
				*destinations[9].(**string) = nil
				*destinations[10].(*int) = 60                           // visit duration in minutes
				*destinations[11].(*float64) = 55.7539 + float64(i)*0.1 // latitude
				*destinations[12].(*float64) = 37.6208 + float64(i)*0.1 // longitude
				*destinations[13].(*time.Time) = time.Date(2026, time.November, 14+i/2, 9, 0, 0, 0, time.UTC)
				*destinations[14].(*time.Time) = time.Date(2026, time.November, 14+i/2, 10, 0, 0, 0, time.UTC)
				return nil
			})
	}

	// Mock route leg rows for 2 legs
	for i := 0; i < 2; i++ {
		mockRouteLegRows[i] = mocks.NewMockRow(ctrl)
		mockRouteLegRows[i].EXPECT().
			Scan(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).
			DoAndReturn(func(destinations ...any) error {
				*destinations[0].(*int) = i + 1 // leg ID
				*destinations[1].(*int) = i + 1 // from stop ID
				*destinations[2].(*int) = i + 2 // to stop ID
				*destinations[3].(*domain.TravelMode) = domain.TravelModeWalking
				*destinations[4].(*int) = 1000 + i*500       // distance
				*destinations[5].(*int64) = int64(60 + i*30) // duration in seconds
				return nil
			})
	}

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createRouteQuery,
			tripID,
			tripVersion,
			status,
		).
		Return(mockRow)

	// Expect 2 route day queries
	for i := 0; i < 2; i++ {
		mockTx.EXPECT().
			QueryRow(
				gomock.Any(),
				createRouteDayQuery,
				gomock.Any(), // routeResponse.ID
				gomock.Any(), // day.Date
				gomock.Any(), // day.DayNumber
			).
			Return(mockRouteDayRows[i])
	}

	// Expect 4 route stop queries (2 per day)
	for i := 0; i < 4; i++ {
		mockTx.EXPECT().
			QueryRow(
				gomock.Any(),
				createRouteStopQuery,
				gomock.Any(), // routeDayID
				gomock.Any(), // stop.PlaceID
				gomock.Any(), // stop.Position
				gomock.Any(), // stop.Name
				gomock.Any(), // stop.Description
				gomock.Any(), // stop.Categories
				gomock.Any(), // stop.Interests
				gomock.Any(), // stop.Address
				gomock.Any(), // stop.WebsiteURL
				gomock.Any(), // stop.Phone
				gomock.Any(), // stop.VisitDurationMinutes
				gomock.Any(), // stop.Coordinates.Latitude
				gomock.Any(), // stop.Coordinates.Longitude
				gomock.Any(), // stop.ArrivalTime
				gomock.Any(), // stop.DepartureTime
			).
			Return(mockRouteStopRows[i])
	}

	// The first stop has one opening hours entry
	mockOpeningHoursRow := mocks.NewMockRow(ctrl)
	mockOpeningHoursRow.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 9
			*destinations[1].(*int) = 18
			return nil
		})
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createOpeningHoursQuery,
			gomock.Any(), // stop ID
			gomock.Any(), // opens at
			gomock.Any(), // closes at
		).
		Return(mockOpeningHoursRow)

	// The second stop also has one opening hours entry
	mockOpeningHoursRow2 := mocks.NewMockRow(ctrl)
	mockOpeningHoursRow2.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 10
			*destinations[1].(*int) = 20
			return nil
		})
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createOpeningHoursQuery,
			gomock.Any(), // stop ID
			gomock.Any(), // opens at
			gomock.Any(), // closes at
		).
		Return(mockOpeningHoursRow2)

	// The third stop also has one opening hours entry
	mockOpeningHoursRow3 := mocks.NewMockRow(ctrl)
	mockOpeningHoursRow3.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 9
			*destinations[1].(*int) = 19
			return nil
		})
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createOpeningHoursQuery,
			gomock.Any(), // stop ID
			gomock.Any(), // opens at
			gomock.Any(), // closes at
		).
		Return(mockOpeningHoursRow3)

	// The fourth stop also has one opening hours entry
	mockOpeningHoursRow4 := mocks.NewMockRow(ctrl)
	mockOpeningHoursRow4.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = 10
			*destinations[1].(*int) = 20
			return nil
		})
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createOpeningHoursQuery,
			gomock.Any(), // stop ID
			gomock.Any(), // opens at
			gomock.Any(), // closes at
		).
		Return(mockOpeningHoursRow4)

	// Expect 2 route leg queries (1 per day)
	for i := 0; i < 2; i++ {
		mockTx.EXPECT().
			QueryRow(
				gomock.Any(),
				createRouteLegQuery,
				gomock.Any(), // routeDayID
				gomock.Any(), // mapStopIDInTripToStopIDInDB[leg.FromStopID]
				gomock.Any(), // mapStopIDInTripToStopIDInDB[leg.ToStopID]
				gomock.Any(), // leg.TravelMode
				gomock.Any(), // leg.Distance
				gomock.Any(), // durationToSeconds(leg.Duration)
			).
			Return(mockRouteLegRows[i])
	}

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	createdRoute, err := routesRepository.CreateRoute(
		context.Background(),
		1,
		&domain.Trip{ID: tripID, Version: tripVersion},
		&domain.RouteCreate{
			Status: status,
			Days: []domain.RouteDay{
				{
					Date:      time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
					DayNumber: 1,
					Stops: []domain.RouteStop{
						{
							PlaceID:       1,
							Position:      1,
							Name:          "Test Stop 1",
							Description:   "Test Description 1",
							OpeningHours:  []domain.OpeningHours{{OpensAt: 9, ClosesAt: 18}},
							Coordinates:   domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
							ArrivalTime:   time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
						},
						{
							PlaceID:       2,
							Position:      2,
							Name:          "Test Stop 2",
							Description:   "Test Description 2",
							OpeningHours:  []domain.OpeningHours{{OpensAt: 10, ClosesAt: 20}},
							Coordinates:   domain.Coordinates{Latitude: 55.8539, Longitude: 37.7208},
							ArrivalTime:   time.Date(2026, time.November, 14, 11, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 14, 12, 0, 0, 0, time.UTC),
						},
					},
					Legs: []domain.RouteLeg{
						{
							FromStopID: 1,
							ToStopID:   2,
							TravelMode: domain.TravelModeWalking,
							Distance:   1000,
							Duration:   time.Minute,
						},
					},
				},
				{
					Date:      time.Date(2026, time.November, 15, 0, 0, 0, 0, time.UTC),
					DayNumber: 2,
					Stops: []domain.RouteStop{
						{
							PlaceID:       3,
							Position:      1,
							Name:          "Test Stop 3",
							Description:   "Test Description 3",
							OpeningHours:  []domain.OpeningHours{{OpensAt: 9, ClosesAt: 19}},
							Coordinates:   domain.Coordinates{Latitude: 55.9539, Longitude: 37.8208},
							ArrivalTime:   time.Date(2026, time.November, 15, 9, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 15, 10, 0, 0, 0, time.UTC),
						},
						{
							PlaceID:       4,
							Position:      2,
							Name:          "Test Stop 4",
							Description:   "Test Description 4",
							OpeningHours:  []domain.OpeningHours{{OpensAt: 10, ClosesAt: 20}},
							Coordinates:   domain.Coordinates{Latitude: 56.0539, Longitude: 37.9208},
							ArrivalTime:   time.Date(2026, time.November, 15, 11, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
						},
					},
					Legs: []domain.RouteLeg{
						{
							FromStopID: 1,
							ToStopID:   2,
							TravelMode: domain.TravelModeWalking,
							Distance:   1500,
							Duration:   time.Minute + 30*time.Second,
						},
					},
				},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, createdRoute)
	assert.Equal(t, 2, len(createdRoute.Days))
	assert.Equal(t, 2, len(createdRoute.Days[0].Stops))
	assert.Equal(t, 2, len(createdRoute.Days[1].Stops))
	assert.Equal(t, 1, len(createdRoute.Days[0].Legs))
	assert.Equal(t, 1, len(createdRoute.Days[1].Legs))
}

func TestDurationToSeconds(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected int64
	}{
		{"zero duration", 0, 0},
		{"1 second", time.Second, 1},
		{"1 minute", time.Minute, 60},
		{"1 hour", time.Hour, 3600},
		{"1 day", 24 * time.Hour, 86400},
		{"1 minute 30 seconds", time.Minute + 30*time.Second, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := durationToSeconds(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDurationFromSeconds(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected time.Duration
	}{
		{"zero seconds", 0, 0},
		{"1 second", 1, time.Second},
		{"60 seconds", 60, time.Minute},
		{"3600 seconds", 3600, time.Hour},
		{"86400 seconds", 86400, 24 * time.Hour},
		{"90 seconds", 90, time.Minute + 30*time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := durationFromSeconds(tt.seconds)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoutesRepository_CreateRoute_BeginError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	beginError := errors.New("begin error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, beginError)

	route, err := routesRepository.CreateRoute(context.Background(), 1, &domain.Trip{ID: tripID}, &domain.RouteCreate{
		Status: domain.RouteStatusGenerated,
	})

	require.ErrorIs(t, err, beginError)
	assert.Nil(t, route)
}
