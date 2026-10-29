//go:build integration

package trips_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsRepository_CreateTrip_Integration_Success(t *testing.T) {
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

	actualTrip, err := tripsRepository.CreateTrip(ctx, ownerID, &tripCreate)
	require.NoError(t, err)
	require.NotNil(t, actualTrip)

	tripID := actualTrip.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := tripsRepository.pool.Exec(ctx, `DELETE FROM app.trips WHERE id=$1;`, tripID)
		return err
	})

	assert.NotZero(t, actualTrip.ID)
	assert.Equal(t, ownerID, actualTrip.OwnerID)
	assert.Equal(t, tripCreate.Name, actualTrip.Name)
	assert.Equal(t, tripCreate.PlanningArea, actualTrip.PlanningArea)
	assert.Equal(t, tripCreate.StartLocation, actualTrip.StartLocation)
	assert.Equal(t, tripCreate.Budget, actualTrip.Budget)
	assert.Equal(t, tripCreate.Interests, actualTrip.Interests)
	assert.Equal(t, tripCreate.TimeZone, actualTrip.TimeZone)
	assert.NotZero(t, actualTrip.CreatedAt)
	assert.NotZero(t, actualTrip.UpdatedAt)
	assert.Nil(t, actualTrip.DeletedAt)
}

func TestTripsRepository_CreateTrip_Integration_ForeignKeyViolation(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	tripCreate := domain.TripCreate{
		OwnerID:       999999, // Non-existent owner
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	actualTrip, err := tripsRepository.CreateTrip(ctx, tripCreate.OwnerID, &tripCreate)
	require.Error(t, err)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
	assert.Nil(t, actualTrip)
}
