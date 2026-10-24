//go:build e2e

package users_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_DeleteUser_E2E_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdUser, errCreatedUser := createUser(t, ctx, 5*time.Second)
	require.Equal(t, "", errCreatedUser)
	require.Equal(t, testutils.TestName, createdUser.Name)
	require.Equal(t, testutils.TestSurname, createdUser.Surname)
	require.Equal(t, testutils.TestUsername, createdUser.Username)
	require.Equal(t, testutils.BirthDateTest(), createdUser.BirthDate)
	require.Equal(t, testutils.TestDescription, createdUser.Description)
	require.Equal(t, testutils.TestEmail, createdUser.Email)
	require.Equal(t, testutils.TestPhoneNumber, createdUser.PhoneNumber)

	createdUserAgain, errCreatedUser := createUser(t, ctx, 5*time.Second)

	require.Contains(t, errCreatedUser, core_errors.ErrConflict.Error())
	require.Nil(t, createdUserAgain)

	userID := createdUser.ID

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	deleteUserRequest, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://localhost:8080/api/v1/users/%d", createdUser.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := client.Do(deleteUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusNoContent, response.StatusCode)

	var deletedAt time.Time
	err = testDeps.pool.QueryRow(
		ctx,
		`SELECT deleted_at FROM app.users WHERE id=$1;`,
		userID,
	).Scan(&deletedAt)
	require.NoError(t, err)
	assert.False(t, deletedAt.IsZero())
}

func TestUsersTransportHTTP_DeleteUser_E2E_InvalidUserID(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	deleteUserRequest, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://localhost:8080/api/v1/users/%d", -1),
		nil,
	)
	require.NoError(t, err)

	response, err := client.Do(deleteUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	var deleteUserErr core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&deleteUserErr))

	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	require.Contains(t, deleteUserErr.Error, core_errors.ErrBadRequest.Error())
}
