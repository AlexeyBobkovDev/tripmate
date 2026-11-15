package users_transport_http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_http_types "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/types"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	t.Run("returns nil when birth date is valid", func(t *testing.T) {
		request := CreateUserRequest{BirthDate: "2006-01-02"}

		require.NoError(t, request.Validate())
	})

	t.Run("returns invalid birth date format error when birth date is invalid", func(t *testing.T) {
		request := CreateUserRequest{BirthDate: "02-01-2006"}

		err := request.Validate()

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidBirthDateFormat)
	})
}

func TestCreateUserDomainFromDTO(t *testing.T) {
	t.Run("returns user create with parsed birth date and password", func(t *testing.T) {
		birthDate := time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)
		userCreate, password, err := createUserDomainFromDTO(CreateUserRequest{
			Name:        "Alex",
			Surname:     "Bobkov",
			Username:    "alex",
			BirthDate:   birthDate.Format("2006-01-02"),
			Description: "Some kind of description",
			Email:       "alex@gmail.com",
			PhoneNumber: "+79998887766",
			Password:    "robust-password-Ab!(*&^1324",
		})

		require.NoError(t, err)
		require.NotNil(t, userCreate)
		assert.Equal(t, "Alex", userCreate.Name)
		assert.Equal(t, "Bobkov", userCreate.Surname)
		assert.Equal(t, "alex", userCreate.Username)
		assert.True(t, userCreate.BirthDate.Equal(birthDate))
		assert.Equal(t, "Some kind of description", userCreate.Description)
		assert.Equal(t, "alex@gmail.com", userCreate.Email)
		assert.Equal(t, "+79998887766", userCreate.PhoneNumber)
		assert.Equal(t, "robust-password-Ab!(*&^1324", password)
	})

	t.Run("returns invalid user birthdate error when birth date cannot be parsed", func(t *testing.T) {
		userCreate, password, err := createUserDomainFromDTO(CreateUserRequest{BirthDate: "not-a-date"})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidUserBirthdate)
		assert.Nil(t, userCreate)
		assert.Empty(t, password)
	})
}

func TestPatchUserDTOToDomain(t *testing.T) {
	t.Run("returns unset user patch when no field is set", func(t *testing.T) {
		userPatch, err := patchUserDTOToDomain(PatchUserRequest{})

		require.NoError(t, err)
		require.NotNil(t, userPatch)
		assert.Equal(t, 0, userPatch.Version)
		assert.False(t, userPatch.Name.Set)
		assert.False(t, userPatch.Surname.Set)
		assert.False(t, userPatch.Username.Set)
		assert.False(t, userPatch.Description.Set)
		assert.False(t, userPatch.BirthDate.Set)
		assert.False(t, userPatch.Email.Set)
		assert.False(t, userPatch.PhoneNumber.Set)
	})

	t.Run("parses set birth date", func(t *testing.T) {
		birthDate := "2006-01-02"
		userPatch, err := patchUserDTOToDomain(PatchUserRequest{
			Version:   1,
			BirthDate: core_http_types.Nullable[string]{Value: &birthDate, Set: true},
		})

		require.NoError(t, err)
		require.NotNil(t, userPatch)
		require.NotNil(t, userPatch.BirthDate.Value)
		assert.True(t, userPatch.BirthDate.Value.Equal(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)))
		assert.True(t, userPatch.BirthDate.Set)
	})

	t.Run("returns bad request error when birth date is set to null", func(t *testing.T) {
		userPatch, err := patchUserDTOToDomain(PatchUserRequest{
			Version:   1,
			BirthDate: core_http_types.Nullable[string]{Value: nil, Set: true},
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrBadRequest)
		assert.Nil(t, userPatch)
	})

	t.Run("returns bad request error when birth date cannot be parsed", func(t *testing.T) {
		birthDate := "not-a-date"
		userPatch, err := patchUserDTOToDomain(PatchUserRequest{
			Version:   1,
			BirthDate: core_http_types.Nullable[string]{Value: &birthDate, Set: true},
		})

		require.Error(t, err)
		assert.Nil(t, userPatch)
	})
}
