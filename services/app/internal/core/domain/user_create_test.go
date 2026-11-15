package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestUserCreate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		create  UserCreate
		wantErr error
	}{
		{
			name: "valid user create",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
		},
		{
			name: "invalid name - empty",
			create: UserCreate{
				Name:         "",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "invalid name - too long",
			create: UserCreate{
				Name:         "This is a very long name that exceeds the maximum allowed length for user names",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "invalid name - contains non-letters",
			create: UserCreate{
				Name:         "John123",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "invalid surname - empty",
			create: UserCreate{
				Name:         "John",
				Surname:      "",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "invalid surname - too long",
			create: UserCreate{
				Name:         "John",
				Surname:      "This is a very long surname that exceeds the maximum allowed length",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "invalid username - empty",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "invalid username - too long",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "This is a very long username that exceeds the maximum allowed length for user usernames",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "invalid username - contains spaces",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "john doe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "invalid birth date - future",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Now().AddDate(1, 0, 0),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid birth date - too young",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Now().AddDate(-13, 0, 0),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid birth date - too old",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Now().AddDate(-151, 0, 0),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid description - too long",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  strings.Repeat("This is a very long description that exceeds the maximum allowed length of 1000 characters. ", 100),
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "invalid email",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "invalid-email",
				PhoneNumber:  "+79001234567",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "invalid phone number",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "invalid",
				PasswordHash: "hashedpassword",
			},
			wantErr: core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			// The password hash is not part of UserCreate.Validate; plaintext
			// passwords are validated separately via ValidatePassword before
			// hashing.
			name: "empty password hash is not validated here",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "",
			},
		},
		{
			name: "short password hash is not validated here",
			create: UserCreate{
				Name:         "John",
				Surname:      "Doe",
				Username:     "johndoe",
				BirthDate:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:  "Test user",
				Email:        "john@example.com",
				PhoneNumber:  "+79001234567",
				PasswordHash: "short",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.create.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password",
			password: "Password1!",
			wantErr:  nil,
		},
		{
			name:     "valid password - longer",
			password: "MySuperSecurePassword123!",
			wantErr:  nil,
		},
		{
			name:     "invalid password - too short",
			password: "Pass1!",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - no digits",
			password: "Password!",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - no letters",
			password: "12345678!",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - no special characters",
			password: "Password123",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - no uppercase",
			password: "password1!",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - no lowercase",
			password: "PASSWORD1!",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - no digits",
			password: "Password!",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - only special characters",
			password: "!@#$%^&*()",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - only letters",
			password: "abcdefg",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "invalid password - only digits",
			password: "12345678",
			wantErr:  core_errors.ErrInvalidUserPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewUserCreate(t *testing.T) {
	create := NewUserCreate(
		"John",
		"Doe",
		"johndoe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"Test user",
		"john@example.com",
		"+79001234567",
	)

	// Verify all fields are set
	assert.Equal(t, "John", create.Name)
	assert.Equal(t, "Doe", create.Surname)
	assert.Equal(t, "johndoe", create.Username)
	assert.Equal(t, time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), create.BirthDate)
	assert.Equal(t, "Test user", create.Description)
	assert.Equal(t, "john@example.com", create.Email)
	assert.Equal(t, "+79001234567", create.PhoneNumber)
	assert.Equal(t, "", create.PasswordHash)

	// Verify validation passes
	require.NoError(t, create.Validate())
}
