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
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_GetUser_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	createUserRequest := testutils.UserCreateTest()

	createdUser, createUserErr := createUser(t, ctx, 15*time.Second)
	require.Empty(t, createUserErr)

	userID := createdUser.ID

	getUserRequest, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%d", baseURL, userID),
		nil,
	)
	require.NoError(t, err)

	response, err := client.Do(getUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	var retrievedUser users_transport_http.UserResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&retrievedUser))

	gotUser := userDomainFromDTO(retrievedUser)

	require.NotNil(t, retrievedUser)

	assert.NotZero(t, createdUser.ID)
	assert.NotZero(t, createdUser.Version)
	assert.NotZero(t, createdUser.CreatedAt)
	assert.NotZero(t, createdUser.UpdatedAt)

	assert.Nil(t, createdUser.DeletedAt)

	assert.Equal(t, createUserRequest.Name, createdUser.Name)
	assert.Equal(t, createUserRequest.Surname, createdUser.Surname)
	assert.Equal(t, createUserRequest.Username, createdUser.Username)
	assert.Equal(t, createUserRequest.BirthDate, createdUser.BirthDate)
	assert.Equal(t, createUserRequest.Description, createdUser.Description)
	assert.Equal(t, createUserRequest.Email, createdUser.Email)
	assert.Equal(t, createUserRequest.PhoneNumber, createdUser.PhoneNumber)

	assert.Equal(t, createdUser.ID, gotUser.ID)
	assert.Equal(t, createdUser.Version, gotUser.Version)
	assert.Equal(t, createdUser.Name, gotUser.Name)
	assert.Equal(t, createdUser.Surname, gotUser.Surname)
	assert.Equal(t, createdUser.Username, gotUser.Username)
	assert.Equal(t, createdUser.BirthDate, gotUser.BirthDate)
	assert.Equal(t, createdUser.Description, gotUser.Description)
	assert.Equal(t, createdUser.Email, gotUser.Email)
	assert.Equal(t, createdUser.PhoneNumber, gotUser.PhoneNumber)
	assert.True(t, createdUser.CreatedAt.Equal(gotUser.CreatedAt))
	assert.True(t, createdUser.UpdatedAt.Equal(gotUser.UpdatedAt))
	assert.Nil(t, createdUser.DeletedAt)
}

func TestUsersTransportHTTP_GetUser_E2E_BadRequest(t *testing.T) {
	userID := -1

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	getUserRequest, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%d", baseURL, userID),
		nil,
	)
	require.NoError(t, err)

	response, err := client.Do(getUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	var errResponse core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&errResponse))

	require.Contains(t, errResponse.Error, core_errors.ErrBadRequest.Error())
}

func TestUsersTransportHTTP_GetUser_E2E_AlreadyDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdUser, createUserErr := createUser(t, ctx, 15*time.Second)
	require.Empty(t, createUserErr)

	userID := createdUser.ID

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	deleteUserRequest, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/%d", baseURL, userID),
		nil,
	)
	require.NoError(t, err)

	response, err := client.Do(deleteUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusNoContent, response.StatusCode)

	getUserRequest, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%d", baseURL, userID),
		nil,
	)
	require.NoError(t, err)

	response, err = client.Do(getUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusNotFound, response.StatusCode)

	var errResponse core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&errResponse))
	require.Contains(t, errResponse.Error, core_errors.ErrNotFound.Error())
}
