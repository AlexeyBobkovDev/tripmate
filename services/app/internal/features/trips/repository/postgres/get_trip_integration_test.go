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

func TestTripsRepository_GetTrip_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	// First create a user to satisfy foreign key constraint on owner_id
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

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, tripID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assert.Equal(t, createdTrip.ID, retrievedTrip.ID)
	assert.Equal(t, createdTrip.OwnerID, retrievedTrip.OwnerID)
	assert.Equal(t, createdTrip.Version, retrievedTrip.Version)
	assert.Equal(t, createdTrip.Name, retrievedTrip.Name)
	assert.Equal(t, createdTrip.PlanningArea, retrievedTrip.PlanningArea)
	assert.Equal(t, createdTrip.StartLocation, retrievedTrip.StartLocation)
	assert.Equal(t, createdTrip.Budget, retrievedTrip.Budget)
	assert.Equal(t, createdTrip.Interests, retrievedTrip.Interests)
	assert.Equal(t, createdTrip.TimeZone, retrievedTrip.TimeZone)
}

func TestTripsRepository_GetTrip_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, 1, 999999)
	require.Nil(t, retrievedTrip)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}
