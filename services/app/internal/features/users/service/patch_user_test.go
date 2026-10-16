package users_service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestPatchUser_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{
			name: "negative ID",
			id:   -1,
		},
		{
			name: "zero ID",
			id:   0,
		},
	}
	userPatch := testhelpers.UserPatchTest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
			mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
			usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

			mockUsersRepository.EXPECT().
				PatchUser(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)
			mockUsersRepository.
				EXPECT().
				GetUser(gomock.Any(), gomock.Any()).
				Times(0)

			retrievedUser, err := usersService.PatchUser(context.Background(), tt.id, &userPatch)
			assert.Nil(t, retrievedUser)

			require.ErrorIs(t, err, core_errors.ErrBadRequest)
		},
		)
	}
}

func TestPatchUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1
	userPatch := testhelpers.UserPatchTest()
	user := testhelpers.UserTest()
	patchedUser := testhelpers.UserTest()
	require.NoError(t, patchedUser.ApplyPatch(&userPatch))

	mockUsersRepository.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(&user, nil)

	mockUsersRepository.
		EXPECT().
		PatchUser(gomock.Any(), userID, &patchedUser).
		Return(&patchedUser, nil)

	retrievedUser, err := usersService.PatchUser(context.Background(), userID, &userPatch)
	require.NoError(t, err)
	assert.Equal(t, patchedUser, *retrievedUser)
}

func TestPatchUser_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	userPatch := testhelpers.UserPatchTest()
	getUserError := errors.New("get user error")

	mockUsersRepository.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(nil, getUserError)

	retrievedUser, err := usersService.PatchUser(context.Background(), userID, &userPatch)
	require.ErrorIs(t, err, getUserError)
	assert.Nil(t, retrievedUser)
}

func TestPatchUser_ApplyPatchError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	userPatch := testhelpers.UserPatchTest()
	// Make userPatch invalid
	tooLongName := strings.Repeat("a", 150)
	userPatch.Name.Value = &tooLongName
	user := testhelpers.UserTest()

	mockUsersRepository.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(&user, nil)

	mockUsersRepository.
		EXPECT().
		PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	retrievedUser, err := usersService.PatchUser(context.Background(), userID, &userPatch)
	require.ErrorIs(t, err, userPatch.Validate())
	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	assert.Nil(t, retrievedUser)
}

func TestPatchUser_PatchUserError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	patchUserRepositoryError := errors.New("patch user repository error")
	userPatch := testhelpers.UserPatchTest()
	user := testhelpers.UserTest()
	patchedUser := testhelpers.UserTest()
	require.NoError(t, patchedUser.ApplyPatch(&userPatch))

	mockUsersRepository.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(&user, nil)

	mockUsersRepository.
		EXPECT().
		PatchUser(gomock.Any(), userID, &patchedUser).
		Return(nil, patchUserRepositoryError)

	retrievedUser, err := usersService.PatchUser(context.Background(), userID, &userPatch)
	require.ErrorIs(t, err, patchUserRepositoryError)
	require.Nil(t, retrievedUser)
}

func TestPatchUser_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	patchUserRepositoryError := errors.New("patch user repository error")

	userPatch := testhelpers.UserPatchTest()
	user := testhelpers.UserTest()

	user.Version = userPatch.Version

	patchedUser := testhelpers.UserTest()
	require.NoError(t, patchedUser.ApplyPatch(&userPatch))

	mockUsersRepository.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(&user, nil)

	mockUsersRepository.
		EXPECT().
		PatchUser(gomock.Any(), userID, &patchedUser).
		Return(nil, errors.Join(patchUserRepositoryError, core_errors.ErrConflict))

	retrievedUser, err := usersService.PatchUser(context.Background(), userID, &userPatch)
	require.ErrorIs(t, err, patchUserRepositoryError)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, retrievedUser)
}

func TestPatchUser_VersionMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(
		mockUsersRepository,
		mockPasswordManager,
	)

	userID := 1
	userPatch := testhelpers.UserPatchTest()
	user := testhelpers.UserTest()
	user.Version = userPatch.Version - 1

	mockUsersRepository.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(&user, nil)

	mockUsersRepository.
		EXPECT().
		PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	retrievedUser, err := usersService.PatchUser(
		context.Background(),
		userID,
		&userPatch,
	)

	require.Error(t, err)
	require.ErrorContains(t, err, "expected version")
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, retrievedUser)
}
