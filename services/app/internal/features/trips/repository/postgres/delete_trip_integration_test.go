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

func TestTripsRepository_DeleteTrip_Integration_Success(t *testing.T) {
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
		_, err := tripsRepository.pool.Exec(ctx, `DELETE FROM app.trips WHERE owner_id=$1;`, ownerID)
		if err != nil {
			return err
		}

		_, err = tripsRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, ownerID)
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

	err = tripsRepository.DeleteTrip(ctx, ownerID, tripID)
	require.NoError(t, err)

	var count int
	err = tripsRepository.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app.trips WHERE id=$1;`, tripID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestTripsRepository_DeleteTrip_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	err := tripsRepository.DeleteTrip(ctx, 1, 999999)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}
