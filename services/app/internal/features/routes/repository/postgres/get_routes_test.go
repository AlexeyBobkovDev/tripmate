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

func TestRoutesRepository_GetRoutes_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	createdAt := time.Now()

	expectedRoutes := []*domain.Route{
		{
			ID:          1,
			TripID:      1,
			TripVersion: 1,
			Version:     1,
			Status:      domain.RouteStatusGenerated,
			Days:        nil,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
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
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].ID
			*destinations[1].(*int) = expectedRoutes[0].TripID
			*destinations[2].(*int) = expectedRoutes[0].TripVersion
			*destinations[3].(*int) = expectedRoutes[0].Version
			*destinations[4].(*domain.RouteStatus) = expectedRoutes[0].Status
			*destinations[5].(*time.Time) = expectedRoutes[0].CreatedAt
			*destinations[6].(*time.Time) = expectedRoutes[0].UpdatedAt
			return nil
		}),
		mockRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteDaysQuery,
			expectedRoutes[0].ID,
		).
		Return(mockRouteDaysRows, nil)

	mockRouteDaysRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, expectedRoutes, routes)
}

func TestRoutesRepository_GetRoutes_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	queryError := errors.New("query error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(nil, queryError)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.ErrorIs(t, err, queryError)
	assert.Nil(t, routes)
}
