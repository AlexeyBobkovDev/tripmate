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
