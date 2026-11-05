package trips_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service/mocks"
)

func TestTripsService_DeleteTrip_InvalidIDs(t *testing.T) {
	tests := []struct {
		name    string
		ownerID int
		tripID  int
	}{
		{
			name:    "zero ownerID",
			ownerID: 0,
			tripID:  1,
		},
		{
			name:    "negative ownerID",
			ownerID: -1,
			tripID:  1,
		},
		{
			name:    "zero tripID",
			ownerID: 1,
			tripID:  0,
		},
		{
			name:    "negative tripID",
			ownerID: 1,
			tripID:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockTripsRepository(ctrl)
			srv := NewTripsService(mockRepo)

			mockRepo.EXPECT().DeleteTrip(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			err := srv.DeleteTrip(context.Background(), tt.ownerID, tt.tripID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
		})
	}
}

func TestTripsService_DeleteTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	mockRepo.EXPECT().
		DeleteTrip(context.Background(), 1, 1).
		Return(nil)

	err := srv.DeleteTrip(context.Background(), 1, 1)
	require.NoError(t, err)
}

func TestTripsService_DeleteTrip_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	mockRepo.EXPECT().
		DeleteTrip(context.Background(), 1, 1).
		Return(core_postgres_pool_errors.ErrViolatedForeignKey)

	err := srv.DeleteTrip(context.Background(), 1, 1)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
}
