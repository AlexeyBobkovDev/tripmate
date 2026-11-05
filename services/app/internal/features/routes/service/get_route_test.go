package routes_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/service/mocks"
)

func TestRoutesService_GetRoute_InvalidIDs(t *testing.T) {
	tests := []struct {
		name    string
		tripID  int
		routeID int
	}{
		{
			name:    "zero tripID",
			tripID:  0,
			routeID: 1,
		},
		{
			name:    "negative tripID",
			tripID:  -1,
			routeID: 1,
		},
		{
			name:    "zero routeID",
			tripID:  1,
			routeID: 0,
		},
		{
			name:    "negative routeID",
			tripID:  1,
			routeID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
			routesService := NewRoutesService(mockRoutesRepository, nil, nil)

			mockRoutesRepository.
				EXPECT().
				GetRoute(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			res, err := routesService.GetRoute(context.Background(), tt.tripID, tt.routeID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, res)
		})
	}
}

func TestRoutesService_GetRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	routesService := NewRoutesService(mockRoutesRepository, nil, nil)

	expectedRoute := &domain.Route{
		ID:          1,
		TripID:      1,
		TripVersion: 1,
		Version:     1,
		Status:      domain.RouteStatusGenerated,
		Days:        []domain.RouteDay{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRoutesRepository.
		EXPECT().
		GetRoute(context.Background(), 1, 1).
		Return(expectedRoute, nil)

	res, err := routesService.GetRoute(context.Background(), 1, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedRoute, res)
}

func TestRoutesService_GetRoute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	routesService := NewRoutesService(mockRoutesRepository, nil, nil)

	mockRoutesRepository.
		EXPECT().
		GetRoute(context.Background(), 1, 1).
		Return(nil, errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	res, err := routesService.GetRoute(context.Background(), 1, 1)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, res)
}
