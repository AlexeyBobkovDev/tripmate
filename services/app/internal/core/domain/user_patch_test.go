package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestUserPatchValidate(t *testing.T) {
	now := time.Now()
	validBirthDate := now.AddDate(-20, 0, 0)

	validName := "Alexey"
	validSurname := "Bobkov"
	validUsername := "alexey"
	validDescription := "Valid description"
	validEmail := "alexey@example.com"
	validPhone := "+14155552671"

	invalidName := "123"
	invalidSurname := "123"
	invalidUsername := ""
	invalidDescription := string(make([]rune, 1001))
	invalidEmail := "not-an-email"
	invalidPhone := "+123"
	invalidBirthDate := now.AddDate(-10, 0, 0)

	tests := []struct {
		name    string
		patch   UserPatch
		wantErr error
	}{
		{
			name: "valid patch with no fields set",
			patch: UserPatch{
				Version: 1,
			},
			wantErr: nil,
		},
		{
			name: "invalid version",
			patch: UserPatch{
				Version: 0,
			},
			wantErr: core_errors.ErrInvalidVersion,
		},
		{
			name: "negative version",
			patch: UserPatch{
				Version: -1,
			},
			wantErr: core_errors.ErrInvalidVersion,
		},

		{
			name: "valid name",
			patch: UserPatch{
				Version: 1,
				Name: Nullable[string]{
					Set:   true,
					Value: &validName,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid name",
			patch: UserPatch{
				Version: 1,
				Name: Nullable[string]{
					Set:   true,
					Value: &invalidName,
				},
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "name not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				Name: Nullable[string]{
					Set:   false,
					Value: &invalidName,
				},
			},
			wantErr: nil,
		},

		{
			name: "valid surname",
			patch: UserPatch{
				Version: 1,
				Surname: Nullable[string]{
					Set:   true,
					Value: &validSurname,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid surname",
			patch: UserPatch{
				Version: 1,
				Surname: Nullable[string]{
					Set:   true,
					Value: &invalidSurname,
				},
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "surname not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				Surname: Nullable[string]{
					Set:   false,
					Value: &invalidSurname,
				},
			},
			wantErr: nil,
		},

		{
			name: "valid username",
			patch: UserPatch{
				Version: 1,
				Username: Nullable[string]{
					Set:   true,
					Value: &validUsername,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid username",
			patch: UserPatch{
				Version: 1,
				Username: Nullable[string]{
					Set:   true,
					Value: &invalidUsername,
				},
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "username not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				Username: Nullable[string]{
					Set:   false,
					Value: &invalidUsername,
				},
			},
			wantErr: nil,
		},

		{
			name: "valid birth date",
			patch: UserPatch{
				Version: 1,
				BirthDate: Nullable[time.Time]{
					Set:   true,
					Value: &validBirthDate,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid birth date",
			patch: UserPatch{
				Version: 1,
				BirthDate: Nullable[time.Time]{
					Set:   true,
					Value: &invalidBirthDate,
				},
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "birth date not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				BirthDate: Nullable[time.Time]{
					Set:   false,
					Value: &invalidBirthDate,
				},
			},
			wantErr: nil,
		},

		{
			name: "valid description",
			patch: UserPatch{
				Version: 1,
				Description: Nullable[string]{
					Set:   true,
					Value: &validDescription,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid description",
			patch: UserPatch{
				Version: 1,
				Description: Nullable[string]{
					Set:   true,
					Value: &invalidDescription,
				},
			},
			wantErr: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "description not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				Description: Nullable[string]{
					Set:   false,
					Value: &invalidDescription,
				},
			},
			wantErr: nil,
		},

		{
			name: "valid email",
			patch: UserPatch{
				Version: 1,
				Email: Nullable[string]{
					Set:   true,
					Value: &validEmail,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid email",
			patch: UserPatch{
				Version: 1,
				Email: Nullable[string]{
					Set:   true,
					Value: &invalidEmail,
				},
			},
			wantErr: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "email not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				Email: Nullable[string]{
					Set:   false,
					Value: &invalidEmail,
				},
			},
			wantErr: nil,
		},

		{
			name: "valid phone number",
			patch: UserPatch{
				Version: 1,
				PhoneNumber: Nullable[string]{
					Set:   true,
					Value: &validPhone,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid phone number",
			patch: UserPatch{
				Version: 1,
				PhoneNumber: Nullable[string]{
					Set:   true,
					Value: &invalidPhone,
				},
			},
			wantErr: core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name: "phone number not set ignores invalid value",
			patch: UserPatch{
				Version: 1,
				PhoneNumber: Nullable[string]{
					Set:   false,
					Value: &invalidPhone,
				},
			},
			wantErr: nil,
		},

		{
			name: "all fields valid",
			patch: UserPatch{
				Version: 1,
				Name: Nullable[string]{
					Set:   true,
					Value: &validName,
				},
				Surname: Nullable[string]{
					Set:   true,
					Value: &validSurname,
				},
				Username: Nullable[string]{
					Set:   true,
					Value: &validUsername,
				},
				BirthDate: Nullable[time.Time]{
					Set:   true,
					Value: &validBirthDate,
				},
				Description: Nullable[string]{
					Set:   true,
					Value: &validDescription,
				},
				Email: Nullable[string]{
					Set:   true,
					Value: &validEmail,
				},
				PhoneNumber: Nullable[string]{
					Set:   true,
					Value: &validPhone,
				},
			},
			wantErr: nil,
		},

		{
			name: "unset fields with invalid values are ignored",
			patch: UserPatch{
				Version: 1,

				Name: Nullable[string]{
					Set:   false,
					Value: &invalidName,
				},
				Surname: Nullable[string]{
					Set:   false,
					Value: &invalidSurname,
				},
				Username: Nullable[string]{
					Set:   false,
					Value: &invalidUsername,
				},
				BirthDate: Nullable[time.Time]{
					Set:   false,
					Value: &invalidBirthDate,
				},
				Description: Nullable[string]{
					Set:   false,
					Value: &invalidDescription,
				},
				Email: Nullable[string]{
					Set:   false,
					Value: &invalidEmail,
				},
				PhoneNumber: Nullable[string]{
					Set:   false,
					Value: &invalidPhone,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.patch.Validate()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewUserPatch(t *testing.T) {
	birthDate := time.Date(
		2005,
		time.May,
		10,
		12,
		30,
		0,
		0,
		time.UTC,
	)

	name := "Alexey"
	surname := "Bobkov"
	username := "alexey"
	description := "Description"
	email := "alexey@example.com"
	phone := "+14155552671"

	namePatch := Nullable[string]{
		Set:   true,
		Value: &name,
	}
	surnamePatch := Nullable[string]{
		Set:   true,
		Value: &surname,
	}
	usernamePatch := Nullable[string]{
		Set:   true,
		Value: &username,
	}
	birthDatePatch := Nullable[time.Time]{
		Set:   true,
		Value: &birthDate,
	}
	descriptionPatch := Nullable[string]{
		Set:   true,
		Value: &description,
	}
	emailPatch := Nullable[string]{
		Set:   true,
		Value: &email,
	}
	phonePatch := Nullable[string]{
		Set:   true,
		Value: &phone,
	}

	got := NewUserPatch(
		42,
		namePatch,
		surnamePatch,
		usernamePatch,
		birthDatePatch,
		descriptionPatch,
		emailPatch,
		phonePatch,
	)

	want := &UserPatch{
		Version:     42,
		Name:        namePatch,
		Surname:     surnamePatch,
		Username:    usernamePatch,
		BirthDate:   birthDatePatch,
		Description: descriptionPatch,
		Email:       emailPatch,
		PhoneNumber: phonePatch,
	}

	require.Equal(t, want, got)
}
