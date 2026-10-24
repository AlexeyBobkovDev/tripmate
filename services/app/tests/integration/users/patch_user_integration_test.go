//go:build integration

package users_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersService_PatchUser_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	createUserRequest := testutils.UserCreateTest()
	createdUser := createUser(t, ctx, &createUserRequest, testutils.Password)

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	originalVersion := createdUser.Version

	patchedUser, err := testDeps.usersService.PatchUser(ctx, createdUser.ID, &userPatch)
	require.NoError(t, err)
	require.NotNil(t, patchedUser)
	assert.Equal(t, originalVersion+1, patchedUser.Version)

	gotUser, err := testDeps.usersService.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)

	assert.Equal(t, patchedUser, gotUser)
	assert.Equal(t, originalVersion+1, gotUser.Version)
}

func TestUsersService_PatchUser_Integration_Conflict_MismatchingUserVersions(t *testing.T) {
	checkIfDatabaseIsDirty(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	createUserRequest := testutils.UserCreateTest()
	createdUser := createUser(t, ctx, &createUserRequest, testutils.Password)

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	// An intentional change of user's version
	userPatch.Version++

	retrievedUser, err := testDeps.usersService.PatchUser(ctx, createdUser.ID, &userPatch)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, retrievedUser)

	gotUser, err := testDeps.usersService.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)

	assert.Equal(t, createdUser.Version, gotUser.Version)
}

func TestUsersService_PatchUser_Integration_NotFound_InvalidUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	userPatch := testutils.UserPatchTest()

	retrievedUser, err := testDeps.usersService.PatchUser(ctx, -1, &userPatch)
	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, retrievedUser)
}

func TestUsersService_PatchUser_Integration_NotFound_AlreadyDeletedUser(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	createUserRequest := testutils.UserCreateTest()
	createdUser := createUser(t, ctx, &createUserRequest, testutils.Password)

	require.NoError(t, testDeps.usersService.DeleteUser(ctx, createdUser.ID))

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	retrievedUser, err := testDeps.usersService.PatchUser(ctx, createdUser.ID, &userPatch)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedUser)
}
