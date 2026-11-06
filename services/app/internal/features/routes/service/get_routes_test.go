package routes_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/service/mocks"
)

func TestRoutesService_GetRoutes_InvalidTripID(t *testing.T) {
	tests := []struct {
		name   string
		tripID int
	}{
		{
			name:   "zero tripID",
			tripID: 0,
		},
		{
			name:   "negative tripID",
			tripID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
			routesService := NewRoutesService(mockRoutesRepository, nil, nil)

			mockRoutesRepository.
				EXPECT().
				GetRoutes(gomock.Any(), gomock.Any()).
				Times(0)

			res, err := routesService.GetRoutes(context.Background(), tt.tripID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, res)
		})
	}
}

func TestRoutesService_GetRoutes_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	routesService := NewRoutesService(mockRoutesRepository, nil, nil)

	createdAt := time.Now()
	expectedRoutes := []*domain.Route{
		{
			ID:          1,
			TripID:      1,
			TripVersion: 1,
			Version:     1,
			Status:      domain.RouteStatusGenerated,
			Days:        []domain.RouteDay{},
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
	}

	mockRoutesRepository.
		EXPECT().
		GetRoutes(context.Background(), 1).
		Return(expectedRoutes, nil)

	res, err := routesService.GetRoutes(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, expectedRoutes, res)
}

func TestRoutesService_GetRoutes_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	routesService := NewRoutesService(mockRoutesRepository, nil, nil)

	repoError := errors.New("database error")
	mockRoutesRepository.
		EXPECT().
		GetRoutes(context.Background(), 1).
		Return(nil, repoError)

	res, err := routesService.GetRoutes(context.Background(), 1)
	require.ErrorIs(t, err, repoError)
	assert.Nil(t, res)
}
