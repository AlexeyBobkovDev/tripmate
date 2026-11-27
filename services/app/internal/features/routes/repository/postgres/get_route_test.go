package routes_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestRoutesRepository_GetRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	mockRouteStopRows := mocks.NewMockRows(ctrl)
	mockOpeningHoursRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	expectedRoute := &domain.Route{
		ID:          1,
		TripID:      1,
		TripVersion: 1,
		Version:     1,
		Status:      domain.RouteStatusGenerated,
		Days: []domain.RouteDay{
			{
				ID:        1,
				Date:      time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
				DayNumber: 1,
				Stops: []domain.RouteStop{
					{
						ID:                   1,
						PlaceID:              1,
						Position:             1,
						Name:                 "Test Stop",
						Description:          "Test Description",
						Address:              "Test Address",
						WebsiteURL:           nil,
						Phone:                nil,
						VisitDurationMinutes: 60,
						OpeningHours:         nil,
						Coordinates:          domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:          time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
						DepartureTime:        time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
					},
				},
				Legs: nil,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).
		Return(mockRow)
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
		Query(gomock.Any(), getRouteDayQuery, expectedRoute.ID).
		Return(mockRouteDayRows, nil)
	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().
		Scan(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.Days[0].ID
			*destinations[1].(*time.Time) = expectedRoute.Days[0].Date
			*destinations[2].(*int) = expectedRoute.Days[0].DayNumber
			return nil
		})
	mockRouteDayRows.EXPECT().Next().Return(false)
	mockRouteDayRows.EXPECT().Err().Return(nil)
	mockRouteDayRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteStopQuery, expectedRoute.Days[0].ID).
		Return(mockRouteStopRows, nil)
	mockRouteStopRows.EXPECT().Next().Return(true)
	mockRouteStopRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			stop := expectedRoute.Days[0].Stops[0]
			*destinations[0].(*int) = stop.ID
			*destinations[1].(*int) = stop.PlaceID
			*destinations[2].(*int) = stop.Position
			*destinations[3].(*string) = stop.Name
			*destinations[4].(*string) = stop.Description
			*destinations[5].(*[]domain.Category) = stop.Categories
			*destinations[6].(*[]domain.Interest) = stop.Interests
			*destinations[7].(*string) = stop.Address
			*destinations[8].(**string) = stop.WebsiteURL
			*destinations[9].(**string) = stop.Phone
			*destinations[10].(*int) = stop.VisitDurationMinutes
			*destinations[11].(*float64) = stop.Coordinates.Latitude
			*destinations[12].(*float64) = stop.Coordinates.Longitude
			*destinations[13].(*time.Time) = stop.ArrivalTime
			*destinations[14].(*time.Time) = stop.DepartureTime
			return nil
		})

	mockPool.EXPECT().
		Query(gomock.Any(), getOpeningHoursQuery, expectedRoute.Days[0].Stops[0].ID).
		Return(mockOpeningHoursRows, nil)
	mockOpeningHoursRows.EXPECT().Next().Return(false)
	mockOpeningHoursRows.EXPECT().Err().Return(nil)
	mockOpeningHoursRows.EXPECT().Close()

	mockRouteStopRows.EXPECT().Next().Return(false)
	mockRouteStopRows.EXPECT().Err().Return(nil)
	mockRouteStopRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteLegQuery, expectedRoute.Days[0].ID).
		Return(mockRouteLegRows, nil)
	mockRouteLegRows.EXPECT().Next().Return(false)
	mockRouteLegRows.EXPECT().Err().Return(nil)
	mockRouteLegRows.EXPECT().Close()

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.NoError(t, err)
	require.NotNil(t, route)
	assert.Equal(t, expectedRoute, route)
	require.Len(t, route.Days, 1)
	require.Len(t, route.Days[0].Stops, 1)
	require.Len(t, route.Days[0].Legs, 0)
}

func TestRoutesRepository_GetRoute_RouteLegScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	mockRouteStopRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	legScanError := errors.New("route leg scan error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteDayQuery, gomock.Any()).
		Return(mockRouteDayRows, nil)
	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(false)
	mockRouteDayRows.EXPECT().Err().Return(nil)
	mockRouteDayRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteStopQuery, gomock.Any()).
		Return(mockRouteStopRows, nil)
	mockRouteStopRows.EXPECT().Next().Return(false)
	mockRouteStopRows.EXPECT().Err().Return(nil)
	mockRouteStopRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteLegQuery, gomock.Any()).
		Return(mockRouteLegRows, nil)
	mockRouteLegRows.EXPECT().Next().Return(true)
	mockRouteLegRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(legScanError)
	mockRouteLegRows.EXPECT().Close()

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.ErrorIs(t, err, legScanError)
	assert.Nil(t, route)
}
