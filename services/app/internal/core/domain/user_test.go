package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestUser_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name: "valid user",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			name: "invalid ID",
			user: User{
				ID:          -1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid version",
			user: User{
				ID:          1,
				Version:     0,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidVersion,
		},
		{
			name: "invalid name",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "invalid surname",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "invalid username",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "invalid birth date - future",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Now().AddDate(1, 0, 0),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid birth date - too young",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Now().AddDate(-13, 0, 0),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid birth date - too old",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Now().AddDate(-151, 0, 0),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid description - too long",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: strings.Repeat("a", 1001),
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "invalid email",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "invalid-email",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "invalid phone number",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "invalid",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name: "invalid created at",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   time.Time{},
				UpdatedAt:   now,
			},
			wantErr: core_errors.ErrInvalidCreatedAt,
		},
		{
			name: "invalid updated at",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   time.Time{},
			},
			wantErr: core_errors.ErrInvalidUpdatedAt,
		},
		{
			name: "invalid deleted at - future",
			user: User{
				ID:          1,
				Version:     1,
				Name:        "John",
				Surname:     "Doe",
				Username:    "johndoe",
				BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description: "Test user",
				Email:       "john@example.com",
				PhoneNumber: "+79001234567",
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   &[]time.Time{now.AddDate(1, 0, 0)}[0],
			},
			wantErr: core_errors.ErrInvalidDeletedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUser_ApplyPatch(t *testing.T) {
	now := time.Now()
	user := User{
		ID:          1,
		Version:     1,
		Name:        "John",
		Surname:     "Doe",
		Username:    "johndoe",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Description: "Test user",
		Email:       "john@example.com",
		PhoneNumber: "+79001234567",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name    string
		patch   *UserPatch
		wantErr error
	}{
		{
			name: "valid patch",
			patch: &UserPatch{
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{"JohnUpdated"}[0], Set: true},
			},
		},
		{
			name: "invalid version",
			patch: &UserPatch{
				Version: 2,
			},
			wantErr: core_errors.ErrDifferentVersions,
		},
		{
			name: "invalid name",
			patch: &UserPatch{
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{"Invalid Name!"}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "patch with invalid surname",
			patch: &UserPatch{
				Version: 1,
				Surname: Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "patch with invalid username",
			patch: &UserPatch{
				Version:  1,
				Username: Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "patch with invalid birth date",
			patch: &UserPatch{
				Version:   1,
				BirthDate: Nullable[time.Time]{Value: &[]time.Time{time.Now()}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "patch with invalid description",
			patch: &UserPatch{
				Version:     1,
				Description: Nullable[string]{Value: &[]string{strings.Repeat("a", 1001)}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "patch with invalid email",
			patch: &UserPatch{
				Version: 1,
				Email:   Nullable[string]{Value: &[]string{"invalid-email"}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "patch with invalid phone number",
			patch: &UserPatch{
				Version:     1,
				PhoneNumber: Nullable[string]{Value: &[]string{"invalid"}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidUserPhoneNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current := user
			err := current.ApplyPatch(tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewUser(t *testing.T) {
	createdAt := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	var deletedAt *time.Time
	user := NewUser(
		7,
		2,
		"John",
		"Doe",
		"johndoe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"Test user",
		"john@example.com",
		"+79001234567",
		createdAt,
		createdAt,
		deletedAt,
	)

	require.NotNil(t, user)
	assert.Equal(t, 7, user.ID)
	assert.Equal(t, 2, user.Version)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "Doe", user.Surname)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, 1990, user.BirthDate.Year())
	assert.Equal(t, "Test user", user.Description)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "+79001234567", user.PhoneNumber)
	assert.Equal(t, createdAt, user.CreatedAt)
	assert.Equal(t, createdAt, user.UpdatedAt)
	assert.Nil(t, user.DeletedAt)
}

func TestUser_Validate_DatesOrderError(t *testing.T) {
	user := User{
		ID:          1,
		Version:     1,
		Name:        "John",
		Surname:     "Doe",
		Username:    "johndoe",
		BirthDate:   time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		Email:       "john@example.com",
		PhoneNumber: "+79001234567",
	}

	err := user.Validate()
	require.ErrorIs(t, err, core_errors.ErrInvalidCreatedAt)
}

func TestUser_ApplyPatch_WithAllFieldsSet(t *testing.T) {
	now := time.Now()
	user := User{
		ID:          1,
		Version:     1,
		Name:        "John",
		Surname:     "Doe",
		Username:    "johndoe",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Description: "Test user",
		Email:       "john@example.com",
		PhoneNumber: "+79001234567",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	patch := &UserPatch{
		Version:     1,
		Surname:     Nullable[string]{Value: &[]string{"DoeUpdated"}[0], Set: true},
		Username:    Nullable[string]{Value: &[]string{"johndoe2"}[0], Set: true},
		BirthDate:   Nullable[time.Time]{Value: &[]time.Time{time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)}[0], Set: true},
		Description: Nullable[string]{Value: &[]string{"Updated description"}[0], Set: true},
		Email:       Nullable[string]{Value: &[]string{"john2@example.com"}[0], Set: true},
		PhoneNumber: Nullable[string]{Value: &[]string{"+79009876543"}[0], Set: true},
	}

	err := user.ApplyPatch(patch)
	require.NoError(t, err)
	require.Equal(t, "DoeUpdated", user.Surname)
	require.Equal(t, "johndoe2", user.Username)
	require.Equal(t, "Updated description", user.Description)
	require.Equal(t, "john2@example.com", user.Email)
	require.Equal(t, "+79009876543", user.PhoneNumber)
}

func TestUser_ApplyPatch_InvalidPatchReturnsErrorWithoutMutation(t *testing.T) {
	now := time.Now()
	user := User{
		ID:          1,
		Version:     1,
		Name:        "John",
		Surname:     "Doe",
		Username:    "johndoe",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Description: "Test user",
		Email:       "",
		PhoneNumber: "+79001234567",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	patch := &UserPatch{
		Version: 1,
		Surname: Nullable[string]{Value: &[]string{"DoeUpdated"}[0], Set: true},
	}

	err := user.ApplyPatch(patch)
	require.Error(t, err)
	assert.NotNil(t, err)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "Doe", user.Surname)
	assert.Equal(t, "", user.Email)
}
