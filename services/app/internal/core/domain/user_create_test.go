package domain

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestUserCreateValidate(t *testing.T) {
	now := time.Now()

	validUserCreate := UserCreate{
		Name:         "Alexey",
		Surname:      "Bobkov",
		Username:     "alexey",
		BirthDate:    now.AddDate(-20, 0, 0),
		Description:  "Software engineer",
		Email:        "alexey@example.com",
		PhoneNumber:  "+14155552671",
		PasswordHash: "some-password-hash",
	}

	tests := []struct {
		name    string
		mutate  func(*UserCreate)
		wantErr error
	}{
		{
			name:    "valid user create",
			mutate:  func(u *UserCreate) {},
			wantErr: nil,
		},
		{
			name: "invalid name",
			mutate: func(u *UserCreate) {
				u.Name = ""
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "invalid surname",
			mutate: func(u *UserCreate) {
				u.Surname = ""
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "invalid username",
			mutate: func(u *UserCreate) {
				u.Username = ""
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "invalid birth date",
			mutate: func(u *UserCreate) {
				u.BirthDate = now.AddDate(-10, 0, 0)
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid description",
			mutate: func(u *UserCreate) {
				u.Description = strings.Repeat("a", 1001)
			},
			wantErr: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "invalid email",
			mutate: func(u *UserCreate) {
				u.Email = "invalid-email"
			},
			wantErr: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "invalid phone number",
			mutate: func(u *UserCreate) {
				u.PhoneNumber = "+123"
			},
			wantErr: core_errors.ErrInvalidUserPhoneNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCreate := validUserCreate
			tt.mutate(&userCreate)

			err := userCreate.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"UserCreate.Validate() = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestUserCreateValidateAllowsNilDeletedAt(t *testing.T) {
	now := time.Now()

	userCreate := UserCreate{
		Name:         "Alexey",
		Surname:      "Bobkov",
		Username:     "alexey",
		BirthDate:    now.AddDate(-20, 0, 0),
		Description:  "Software engineer",
		Email:        "alexey@example.com",
		PhoneNumber:  "+14155552671",
		PasswordHash: "some-password-hash",
	}

	if err := userCreate.Validate(); err != nil {
		t.Fatalf("UserCreate.Validate() = %v, want nil", err)
	}
}

func TestNewUserCreate(t *testing.T) {
	birthDate := time.Date(
		2005,
		time.January,
		15,
		12,
		30,
		0,
		0,
		time.UTC,
	)

	got := NewUserCreate(
		"Alexey",
		"Bobkov",
		"alexey",
		birthDate,
		"Software engineer",
		"alexey@example.com",
		"+14155552671",
	)

	want := &UserCreate{
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   birthDate,
		Description: "Software engineer",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"NewUserCreate() = %#v, want %#v",
			got,
			want,
		)
	}
}

func TestNewUserCreateDoesNotSetPasswordHash(t *testing.T) {
	birthDate := time.Date(
		2005,
		time.January,
		15,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	got := NewUserCreate(
		"Alexey",
		"Bobkov",
		"alexey",
		birthDate,
		"Software engineer",
		"alexey@example.com",
		"+14155552671",
	)

	if got.PasswordHash != "" {
		t.Fatalf(
			"NewUserCreate().PasswordHash = %q, want empty string",
			got.PasswordHash,
		)
	}
}
