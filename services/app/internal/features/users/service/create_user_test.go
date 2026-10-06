package users_service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	const (
		password     = "robust-password-Ab!(*&^1324"
		passwordHash = "$argon2id$v=19$m=65536,t=3,p=2$hashed-password"
	)

	testUserCreate := testhelpers.UserCreateTest()

	expectedUser := domain.NewUser(
		0,
		1,
		testUserCreate.Name,
		testUserCreate.Surname,
		testUserCreate.Username,
		testUserCreate.BirthDate,
		testUserCreate.Description,
		testUserCreate.Email,
		testUserCreate.PhoneNumber,
		testhelpers.CreatedAtTest(),
		testhelpers.CreatedAtTest(),
		nil,
	)

	mockPasswordManager.
		EXPECT().
		Hash(password).
		Return(passwordHash, nil)

	mockUsersRepository.
		EXPECT().
		CreateUser(context.Background(), &testUserCreate).
		Return(expectedUser, nil)

	returnedUser, err := usersService.CreateUser(context.Background(), &testUserCreate, password)
	require.NoError(t, err)
	assert.Equal(t, expectedUser, returnedUser)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	const (
		password     = "robust-password-Ab!(*&^1324"
		passwordHash = "$argon2id$v=19$m=65536,t=3,p=2$hashed-password"
	)

	testUserCreate := testhelpers.UserCreateTest()

	mockPasswordManager.
		EXPECT().
		Hash(password).
		Return(passwordHash, nil)

	mockUsersRepository.
		EXPECT().
		CreateUser(context.Background(), &testUserCreate).
		Return(nil, core_postgres_pool_errors.ErrViolatedForeignKey)

	returnedUser, err := usersService.CreateUser(context.Background(), &testUserCreate, password)
	require.Nil(t, returnedUser)
	assert.Error(t, err)
	assert.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
	assert.ErrorContains(t, err, "failed to create user")
}

func TestCreateUser_UserValidationError(t *testing.T) {
	const (
		password = "robust-password-Ab!(*&^1324"
	)

	testUserCreate := testhelpers.UserCreateTest()
	tests := []struct {
		name string
		fn   func(domain.UserCreate) domain.UserCreate
		err  error
	}{
		{
			name: "empty name",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Name = ""
				return user
			},
			err: core_errors.ErrInvalidUserName,
		},
		{
			name: "name too long",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Name = strings.Repeat("a", 51)
				return user
			},
			err: core_errors.ErrInvalidUserName,
		},
		{
			name: "name contains numbers",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Name = "Alex123"
				return user
			},
			err: core_errors.ErrInvalidUserName,
		},
		{
			name: "empty surname",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Surname = ""
				return user
			},
			err: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "surname too long",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Surname = strings.Repeat("a", 51)
				return user
			},
			err: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "empty username",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Username = ""
				return user
			},
			err: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "username contains spaces",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Username = "alex bob"
				return user
			},
			err: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "username too long",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Username = strings.Repeat("a", 101)
				return user
			},
			err: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "empty email",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Email = ""
				return user
			},
			err: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "invalid email format",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Email = "alex"
				return user
			},
			err: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "empty phone number",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.PhoneNumber = ""
				return user
			},
			err: core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name: "invalid phone number format",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.PhoneNumber = "123456"
				return user
			},
			err: core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name: "description too long",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.Description = strings.Repeat("a", 1001)
				return user
			},
			err: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "birth date too early",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.BirthDate = time.Date(
					1800,
					1,
					1,
					0,
					0,
					0,
					0,
					time.UTC,
				)
				return user
			},
			err: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "birth date in future",
			fn: func(user domain.UserCreate) domain.UserCreate {
				user.BirthDate = time.Now().AddDate(10, 0, 0)
				return user
			},
			err: core_errors.ErrInvalidUserBirthdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
			mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
			usersService := NewUsersService(mockUsersRepository, mockPasswordManager)
			user := tt.fn(testUserCreate)

			mockPasswordManager.
				EXPECT().
				Hash(gomock.Any()).
				Times(0)

			mockUsersRepository.
				EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)

			returnedUser, err := usersService.CreateUser(context.Background(), &user, password)
			require.Nil(t, returnedUser)

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.err)
		},
		)
	}
}

func TestCreateUser_PasswordValidationError(t *testing.T) {
	testUserCreate := testhelpers.UserCreateTest()
	tests := []struct {
		name     string
		password string
		err      error
	}{
		{
			name:     "empty password",
			password: "",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password too short",
			password: "Ab1!",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password exactly minimum length but invalid",
			password: "12345678",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password without uppercase letter",
			password: "robust-password-123!@#",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password without lowercase letter",
			password: "ROBUST-PASSWORD-123!@#",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password without digit",
			password: "Robust-password-!@#",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password without special character",
			password: "RobustPassword123",
			err:      core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "password too long",
			password: strings.Repeat("A", 101),
			err:      core_errors.ErrInvalidUserPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
			mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
			usersService := NewUsersService(mockUsersRepository, mockPasswordManager)
			mockPasswordManager.
				EXPECT().
				Hash(gomock.Any()).
				Times(0)

			mockUsersRepository.
				EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)

			returnedUser, err := usersService.CreateUser(context.Background(), &testUserCreate, tt.password)
			require.Nil(t, returnedUser)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.err)
		},
		)
	}
}
