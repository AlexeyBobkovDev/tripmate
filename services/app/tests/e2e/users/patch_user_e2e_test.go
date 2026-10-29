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

func TestUsersTransportHTTP_PatchUser_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	userPatch := testutils.UserPatchTest()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// Prepare the data (create user)

	createdUser, createUserErr := createUser(t, ctx, 15*time.Second)
	require.Empty(t, createUserErr)

	userID := createdUser.ID

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	originalVersion := createdUser.Version

	patchUserRequest, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/%d", baseURL, userID),
		testPatchRequestBody(),
	)
	require.NoError(t, err)

	response, err := client.Do(patchUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	var retrievedUser users_transport_http.UserResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&retrievedUser))
	require.NotNil(t, retrievedUser)
	assert.Equal(t, originalVersion+1, retrievedUser.Version)

	patchedUser := userDomainFromDTO(retrievedUser)

	assert.Equal(t, createdUser.ID, patchedUser.ID)
	assert.Equal(t, createdUser.Version+1, patchedUser.Version)
	assert.Equal(t, createdUser.Name, patchedUser.Name)
	assert.Equal(t, createdUser.Surname, patchedUser.Surname)
	assert.Equal(t, createdUser.Username, patchedUser.Username)
	assert.Equal(t, createdUser.BirthDate, patchedUser.BirthDate)
	assert.Equal(t, createdUser.Description, patchedUser.Description)
	assert.Equal(t, createdUser.Email, patchedUser.Email)
	assert.Equal(t, createdUser.PhoneNumber, patchedUser.PhoneNumber)
	assert.True(t, createdUser.CreatedAt.Equal(patchedUser.CreatedAt))
	assert.True(t, createdUser.UpdatedAt.Equal(patchedUser.UpdatedAt))
	assert.Nil(t, createdUser.DeletedAt)
	assert.Equal(t, originalVersion+1, patchedUser.Version)
}

func TestUsersTransportHTTP_PatchUser_E2E_Conflict_MismatchingUserVersions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	createdUser, createUserErr := createUser(t, ctx, 15*time.Second)
	require.Empty(t, createUserErr)

	userID := createdUser.ID

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	// An intentional change of user's version
	userPatch.Version++

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	patchUserRequest, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/%d", baseURL, userID),
		testPatchRequestBodyWithIncreasedByOneVersion(createdUser.Version),
	)
	require.NoError(t, err)

	response, err := client.Do(patchUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusConflict, response.StatusCode)

	var errResponse core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&errResponse))
	require.Contains(t, errResponse.Error, core_errors.ErrConflict.Error())
}

func TestUsersTransportHTTP_PatchUser_E2E_NotFound_InvalidUserID(t *testing.T) {
	userID := -1

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	patchUserRequest, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/%d", baseURL, userID),
		testPatchRequestBody(),
	)
	require.NoError(t, err)

	response, err := client.Do(patchUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	var errResponse core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&errResponse))

	require.Contains(t, errResponse.Error, core_errors.ErrBadRequest.Error())
}

func TestUsersTransportHTTP_PatchUser_E2E_NotFound_AlreadyDeletedUser(t *testing.T) {
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

	patchUserRequest, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/%d", baseURL, userID),
		testPatchRequestBody(),
	)
	require.NoError(t, err)

	response, err = client.Do(patchUserRequest)
	require.NoError(t, err)

	require.Equal(t, http.StatusNotFound, response.StatusCode)

	var errResponse core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&errResponse))
	require.Contains(t, errResponse.Error, core_errors.ErrNotFound.Error())
}
