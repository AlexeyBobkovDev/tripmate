package trips_service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service/mocks"
)

func TestTripsService_GetTrips_InvalidOwnerID(t *testing.T) {
	tests := []struct {
		name    string
		ownerID int
	}{
		{
			name:    "zero ownerID",
			ownerID: 0,
		},
		{
			name:    "negative ownerID",
			ownerID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockTripsRepository(ctrl)
			srv := NewTripsService(mockRepo)

			mockRepo.EXPECT().GetTrips(gomock.Any(), gomock.Any()).Times(0)

			res, err := srv.GetTrips(context.Background(), tt.ownerID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, res)
		})
	}
}

func TestTripsService_GetTrips_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	expectedTrips := []domain.Trip{
		{ID: 1, OwnerID: 1, Name: "Moscow Trip"},
	}

	mockRepo.EXPECT().
		GetTrips(context.Background(), 1).
		Return(&expectedTrips, nil)

	res, err := srv.GetTrips(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, &expectedTrips, res)
}

func TestTripsService_GetTrips_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	repoErr := errors.New("database error")
	mockRepo.EXPECT().
		GetTrips(context.Background(), 1).
		Return(nil, repoErr)

	res, err := srv.GetTrips(context.Background(), 1)
	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, res)
}
