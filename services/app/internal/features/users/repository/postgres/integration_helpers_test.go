//go:build integration

package users_postgres_repository

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func registerCleanup(t *testing.T, cleanup func(context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		if err := cleanup(cleanupCtx); err != nil {
			t.Errorf("integration test cleanup failed: %v", err)
		}
	})
}

// checkIfDatabaseIsDirty guards tests against leaked rows from a previous test.
// Integration tests in this package run sequentially, so the database must be
// empty before each test starts. Checking every table also verifies that
// cascading deletes and cleanups worked.
func checkIfDatabaseIsDirty(t *testing.T) {
	t.Helper()

	tables := []string{
		"app.users",
		"app.trips",
		"app.meal_requirements",
		"app.routes",
		"app.route_days",
		"app.route_stops",
		"app.route_legs",
	}

	for _, table := range tables {
		var counter int

		require.NoError(t, poolForIntegrationTesting.QueryRow(
			context.Background(),
			fmt.Sprintf(`SELECT COUNT(*) FROM %s;`, table),
		).Scan(&counter))
		require.Zero(t, counter, "dirty database state: %s is not empty", table)
	}
}

// uniqueUserCreate builds a valid user whose unique fields do not collide with
// other users created in the same test.
func uniqueUserCreate(t *testing.T) domain.UserCreate {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)

	userCreate := testutils.UserCreateTest()
	userCreate.Username = "users-" + suffix
	userCreate.Email = fmt.Sprintf("users-%s@example.com", suffix)
	userCreate.PhoneNumber = "+7999" + suffix[len(suffix)-7:]
	userCreate.PasswordHash = "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$c29tZWhhc2g"

	return userCreate
}

// userPatchFixture returns a patch with values that differ from
// testutils.UserCreateTest, so a successful patch is observable.
func userPatchFixture(version int) domain.UserPatch {
	name := "Patchedname"
	surname := "Patchedsurname"
	username := "patched-username"
	birthDate := time.Date(2000, time.January, 2, 0, 0, 0, 0, time.UTC)
	description := "Patched description"
	email := "patched@example.com"
	phoneNumber := "+79991234567"

	return domain.UserPatch{
		Version:     version,
		Name:        domain.Nullable[string]{Value: &name, Set: true},
		Surname:     domain.Nullable[string]{Value: &surname, Set: true},
		Username:    domain.Nullable[string]{Value: &username, Set: true},
		BirthDate:   domain.Nullable[time.Time]{Value: &birthDate, Set: true},
		Description: domain.Nullable[string]{Value: &description, Set: true},
		Email:       domain.Nullable[string]{Value: &email, Set: true},
		PhoneNumber: domain.Nullable[string]{Value: &phoneNumber, Set: true},
	}
}
