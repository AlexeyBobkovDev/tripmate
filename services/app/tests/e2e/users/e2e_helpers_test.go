//go:build e2e

package users_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
)

func registerCleanup(t *testing.T, cleanup func(context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		if err := cleanup(cleanupCtx); err != nil {
			t.Errorf("E2E test cleanup failed: %v", err)
		}
	})
}

// In function createUsr an error is returned as a string,
// because HTTP does not support sending errors
func createUser(t *testing.T, ctx context.Context, timeout time.Duration) (*domain.User, string) {
	client := &http.Client{
		Timeout: timeout,
	}
	createUserRequest, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/v1/users",
		testRequestBodyWithDefaults(),
	)
	require.NoError(t, err)

	response, err := client.Do(createUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	var (
		createdUser users_transport_http.UserResponse
		errResponse core_http_response.ErrorResponse
	)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		require.NoError(t, json.NewDecoder(response.Body).Decode(&createdUser))
		registerCleanup(t, func(ctx context.Context) error {
			_, err := testDeps.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
			return err
		})
		require.Equal(t, http.StatusCreated, response.StatusCode)
		require.Equal(t, "application/json", response.Header.Get("Content-Type"))

		return userDomainFromDTO(createdUser), ""
	} else {
		require.NoError(t, json.NewDecoder(response.Body).Decode(&errResponse))
		return nil, errResponse.Error
	}
}

func userDomainFromDTO(dto users_transport_http.UserResponse) *domain.User {
	birthDate, err := time.Parse(core_consts.Layout, dto.BirthDate)
	if err != nil {
		panic(fmt.Errorf("failed to parse birthdate: %w", err))
	}
	return &domain.User{
		ID:          dto.ID,
		Version:     dto.Version,
		Name:        dto.Name,
		Surname:     dto.Surname,
		Username:    dto.Username,
		Description: dto.Description,
		BirthDate:   birthDate,
		Email:       dto.Email,
		PhoneNumber: dto.PhoneNumber,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		DeletedAt:   dto.DeletedAt,
	}
}

func checkIfDatabaseIsDirty(t *testing.T) {
	var counter int

	require.NoError(t, testDeps.pool.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM app.users;
		`,
	).Scan(&counter))
	if counter != 0 {
		panic("dirty database state!!!")
	}
}
