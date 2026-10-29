package trips_service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service/mocks"
)

func TestTripsService_GetTrip_InvalidIDs(t *testing.T) {
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

			mockRepo.EXPECT().GetTrip(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			res, err := srv.GetTrip(context.Background(), tt.ownerID, tt.tripID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, res)
		})
	}
}

func TestTripsService_GetTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	expectedTrip := domain.Trip{ID: 1, OwnerID: 1, Name: "Moscow Trip"}

	mockRepo.EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(&expectedTrip, nil)

	res, err := srv.GetTrip(context.Background(), 1, 1)
	require.NoError(t, err)
	assert.Equal(t, &expectedTrip, res)
}

func TestTripsService_GetTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	mockRepo.EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(nil, errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	res, err := srv.GetTrip(context.Background(), 1, 1)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, res)
}
