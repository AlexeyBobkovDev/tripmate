//go:build integration

package trips_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsRepository_PatchTrip_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	createUserQuery := `
		INSERT INTO app.users (name, surname, username, birth_date, description, email, phone_number, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`
	userCreate := testutils.UserCreateTest()
	var ownerID int
	err := poolForIntegrationTesting.QueryRow(
		ctx,
		createUserQuery,
		userCreate.Name,
		userCreate.Surname,
		userCreate.Username,
		userCreate.BirthDate,
		userCreate.Description,
		userCreate.Email,
		userCreate.PhoneNumber,
		userCreate.PasswordHash,
	).Scan(&ownerID)
	require.NoError(t, err)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := tripsRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, ownerID)
		return err
	})

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	tripCreate := domain.TripCreate{
		OwnerID:       ownerID,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, &tripCreate)
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	tripID := createdTrip.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := tripsRepository.pool.Exec(ctx, `DELETE FROM app.trips WHERE id=$1;`, tripID)
		return err
	})

	newName := "Updated Trip Name"
	createdTrip.Name = newName
	originalVersion := createdTrip.Version

	patchedTrip, err := tripsRepository.PatchTrip(ctx, ownerID, tripID, createdTrip)
	require.NoError(t, err)
	require.NotNil(t, patchedTrip)

	assert.Equal(t, tripID, patchedTrip.ID)
	assert.Equal(t, ownerID, patchedTrip.OwnerID)
	assert.Equal(t, originalVersion+1, patchedTrip.Version)
	assert.Equal(t, newName, patchedTrip.Name)
}

func TestTripsRepository_PatchTrip_Integration_Conflict(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	createUserQuery := `
		INSERT INTO app.users (name, surname, username, birth_date, description, email, phone_number, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`
	userCreate := testutils.UserCreateTest()
	var ownerID int
	err := poolForIntegrationTesting.QueryRow(
		ctx,
		createUserQuery,
		userCreate.Name,
		userCreate.Surname,
		userCreate.Username,
		userCreate.BirthDate,
		userCreate.Description,
		userCreate.Email,
		userCreate.PhoneNumber,
		userCreate.PasswordHash,
	).Scan(&ownerID)
	require.NoError(t, err)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := tripsRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, ownerID)
		return err
	})

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	tripCreate := domain.TripCreate{
		OwnerID:       ownerID,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, &tripCreate)
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	tripID := createdTrip.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := tripsRepository.pool.Exec(ctx, `DELETE FROM app.trips WHERE id=$1;`, tripID)
		return err
	})

	// Intentional version mismatch
	createdTrip.Version = 999
	patchedTrip, err := tripsRepository.PatchTrip(ctx, ownerID, tripID, createdTrip)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	assert.Nil(t, patchedTrip)
}

func TestTripsRepository_PatchTrip_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	trip := &domain.Trip{
		ID:            999999,
		OwnerID:       1,
		Version:       1,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	patchedTrip, err := tripsRepository.PatchTrip(ctx, 1, 999999, trip)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, patchedTrip)
}
