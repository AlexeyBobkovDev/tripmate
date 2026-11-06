package routes_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/service/mocks"
)

func TestRoutesService_DeleteRoute_InvalidIDs(t *testing.T) {
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
				DeleteRoute(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			err := routesService.DeleteRoute(context.Background(), tt.tripID, tt.routeID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
		})
	}
}

func TestRoutesService_DeleteRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	routesService := NewRoutesService(mockRoutesRepository, nil, nil)

	mockRoutesRepository.
		EXPECT().
		DeleteRoute(context.Background(), 1, 1).
		Return(nil)

	err := routesService.DeleteRoute(context.Background(), 1, 1)
	require.NoError(t, err)
}

func TestRoutesService_DeleteRoute_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	routesService := NewRoutesService(mockRoutesRepository, nil, nil)

	mockRoutesRepository.
		EXPECT().
		DeleteRoute(context.Background(), 1, 1).
		Return(core_postgres_pool_errors.ErrViolatedForeignKey)

	err := routesService.DeleteRoute(context.Background(), 1, 1)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
}
