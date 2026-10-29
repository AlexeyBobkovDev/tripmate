package trips_service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service/mocks"
)

func TestTripsService_PatchTrip_InvalidIDs(t *testing.T) {
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

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	tripPatch := domain.NewTripPatch(
		1,
		1,
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[domain.Location]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: &startDate, Set: true},
		domain.Nullable[time.Time]{Value: &endDate, Set: true},
		domain.Nullable[domain.Money]{Value: nil, Set: false},
		domain.Nullable[[]domain.Interest]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockTripsRepository(ctrl)
			srv := NewTripsService(mockRepo)

			mockRepo.EXPECT().GetTrip(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			mockRepo.EXPECT().PatchTrip(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			res, err := srv.PatchTrip(context.Background(), tt.ownerID, tt.tripID, tripPatch)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, res)
		})
	}
}

func TestTripsService_PatchTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	existingTrip := &domain.Trip{
		ID:            1,
		OwnerID:       1,
		Version:       1,
		Name:          "Old Trip",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
		CreatedAt:     time.Now().Add(-72 * time.Hour),
		UpdatedAt:     time.Now().Add(-72 * time.Hour),
	}

	newName := "Updated Trip"
	tripPatch := domain.NewTripPatch(
		1,
		1,
		domain.Nullable[string]{Value: &newName, Set: true},
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[domain.Location]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[domain.Money]{Value: nil, Set: false},
		domain.Nullable[[]domain.Interest]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
	)

	patchedTrip := *existingTrip
	patchedTrip.Version = 2
	patchedTrip.Name = newName

	mockRepo.EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(existingTrip, nil)

	mockRepo.EXPECT().
		PatchTrip(context.Background(), 1, 1, gomock.Any()).
		Return(&patchedTrip, nil)

	res, err := srv.PatchTrip(context.Background(), 1, 1, tripPatch)
	require.NoError(t, err)
	assert.Equal(t, &patchedTrip, res)
}

func TestTripsService_PatchTrip_GetTripError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	getErr := errors.New("not found")
	tripPatch := domain.NewTripPatch(
		1,
		1,
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[domain.Location]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[domain.Money]{Value: nil, Set: false},
		domain.Nullable[[]domain.Interest]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
	)

	mockRepo.EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(nil, getErr)

	res, err := srv.PatchTrip(context.Background(), 1, 1, tripPatch)
	require.ErrorIs(t, err, getErr)
	assert.Nil(t, res)
}

func TestTripsService_PatchTrip_VersionMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	existingTrip := &domain.Trip{
		ID:      1,
		OwnerID: 1,
		Version: 2,
	}

	tripPatch := domain.NewTripPatch(
		1,
		1, // Mismatched version
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[domain.Location]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[domain.Money]{Value: nil, Set: false},
		domain.Nullable[[]domain.Interest]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
	)

	mockRepo.EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(existingTrip, nil)

	res, err := srv.PatchTrip(context.Background(), 1, 1, tripPatch)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	assert.Nil(t, res)
}

func TestTripsService_PatchTrip_ApplyPatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTripsRepository(ctrl)
	srv := NewTripsService(mockRepo)

	existingTrip := &domain.Trip{
		ID:      1,
		OwnerID: 1,
		Version: 1,
	}

	tooLongName := strings.Repeat("a", 200)
	tripPatch := domain.NewTripPatch(
		1,
		1,
		domain.Nullable[string]{Value: &tooLongName, Set: true},
		domain.Nullable[string]{Value: nil, Set: false},
		domain.Nullable[domain.Location]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[time.Time]{Value: nil, Set: false},
		domain.Nullable[domain.Money]{Value: nil, Set: false},
		domain.Nullable[[]domain.Interest]{Value: nil, Set: false},
		domain.Nullable[string]{Value: nil, Set: false},
	)

	mockRepo.EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(existingTrip, nil)

	res, err := srv.PatchTrip(context.Background(), 1, 1, tripPatch)
	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	assert.Nil(t, res)
}
