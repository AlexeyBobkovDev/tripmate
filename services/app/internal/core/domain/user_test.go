package domain

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestUserValidate(t *testing.T) {
	now := time.Now()

	validUser := User{
		ID:          1,
		Version:     1,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   now.AddDate(-20, 0, 0),
		Description: "Software engineer",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
		DeletedAt:   nil,
	}

	tests := []struct {
		name    string
		mutate  func(*User)
		wantErr error
	}{
		{
			name:    "valid user",
			mutate:  func(u *User) {},
			wantErr: nil,
		},
		{
			name: "invalid ID",
			mutate: func(u *User) {
				u.ID = 0
			},
			wantErr: core_errors.ErrInvalidUserID,
		},
		{
			name: "invalid version",
			mutate: func(u *User) {
				u.Version = 0
			},
			wantErr: core_errors.ErrInvalidUserVersion,
		},
		{
			name: "invalid name",
			mutate: func(u *User) {
				u.Name = ""
			},
			wantErr: core_errors.ErrInvalidUserName,
		},
		{
			name: "invalid surname",
			mutate: func(u *User) {
				u.Surname = ""
			},
			wantErr: core_errors.ErrInvalidUserSurname,
		},
		{
			name: "invalid username",
			mutate: func(u *User) {
				u.Username = ""
			},
			wantErr: core_errors.ErrInvalidUserUsername,
		},
		{
			name: "invalid birth date",
			mutate: func(u *User) {
				u.BirthDate = now.AddDate(-10, 0, 0)
			},
			wantErr: core_errors.ErrInvalidUserBirthdate,
		},
		{
			name: "invalid description",
			mutate: func(u *User) {
				u.Description = strings.Repeat("a", 1001)
			},
			wantErr: core_errors.ErrInvalidUserDescription,
		},
		{
			name: "invalid email",
			mutate: func(u *User) {
				u.Email = "not-an-email"
			},
			wantErr: core_errors.ErrInvalidUserEmail,
		},
		{
			name: "invalid phone number",
			mutate: func(u *User) {
				u.PhoneNumber = "+123"
			},
			wantErr: core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name: "created at in the future",
			mutate: func(u *User) {
				u.CreatedAt = now.Add(time.Hour)
			},
			wantErr: core_errors.ErrInvalidUserCreatedAt,
		},
		{
			name: "updated at in the future",
			mutate: func(u *User) {
				u.UpdatedAt = now.Add(time.Hour)
			},
			wantErr: core_errors.ErrInvalidUserUpdatedAt,
		},
		{
			name: "deleted at in the future",
			mutate: func(u *User) {
				deletedAt := now.Add(time.Hour)
				u.DeletedAt = &deletedAt
			},
			wantErr: core_errors.ErrInvalidUserDeletedAt,
		},
		{
			name: "created before birth date",
			mutate: func(u *User) {
				u.CreatedAt = u.BirthDate
			},
			wantErr: core_errors.ErrInvalidUserCreatedAt,
		},
		{
			name: "updated before created",
			mutate: func(u *User) {
				u.UpdatedAt = u.CreatedAt.Add(-time.Second)
			},
			wantErr: core_errors.ErrInvalidUserUpdatedAt,
		},
		{
			name: "deleted before updated",
			mutate: func(u *User) {
				deletedAt := u.UpdatedAt.Add(-time.Second)
				u.DeletedAt = &deletedAt
			},
			wantErr: core_errors.ErrInvalidUserDeletedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := validUser
			tt.mutate(&user)

			err := user.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"User.Validate() = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestUserApplyPatch(t *testing.T) {
	now := time.Now()

	validUser := User{
		ID:          1,
		Version:     1,
		Name:        "User",
		Surname:     "User",
		Username:    "user123",
		BirthDate:   now.AddDate(-20, 0, 0),
		Description: "Original description",
		Email:       "user@gmail.com",
		PhoneNumber: "+79998887766",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
		DeletedAt:   nil,
	}

	tests := []struct {
		name     string
		patch    UserPatch
		wantErr  error
		wantUser User
	}{
		{
			name: "valid empty patch",
			patch: UserPatch{
				Version: 1,
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				return user
			}(),
		},
		{
			name: "different version",
			patch: UserPatch{
				Version: 2,
			},
			wantErr:  core_errors.ErrDifferentVersions,
			wantUser: validUser,
		},
		{
			name: "updates name",
			patch: UserPatch{
				Version: 1,
				Name: Nullable[string]{
					Set: true,
					Value: func() *string {
						v := "Alexey"
						return &v
					}(),
				},
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				user.Name = "Alexey"
				return user
			}(),
		},
		{
			name: "updates surname",
			patch: UserPatch{
				Version: 1,
				Surname: Nullable[string]{
					Set: true,
					Value: func() *string {
						v := "Bobkov"
						return &v
					}(),
				},
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				user.Surname = "Bobkov"
				return user
			}(),
		},
		{
			name: "updates username",
			patch: UserPatch{
				Version: 1,
				Username: Nullable[string]{
					Set: true,
					Value: func() *string {
						v := "alexey"
						return &v
					}(),
				},
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				user.Username = "alexey"
				return user
			}(),
		},
		{
			name: "updates description",
			patch: UserPatch{
				Version: 1,
				Description: Nullable[string]{
					Set: true,
					Value: func() *string {
						v := "New description"
						return &v
					}(),
				},
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				user.Description = "New description"
				return user
			}(),
		},
		{
			name: "updates email",
			patch: UserPatch{
				Version: 1,
				Email: Nullable[string]{
					Set: true,
					Value: func() *string {
						v := "alexey@example.com"
						return &v
					}(),
				},
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				user.Email = "alexey@example.com"
				return user
			}(),
		},
		{
			name: "updates phone number",
			patch: UserPatch{
				Version: 1,
				PhoneNumber: Nullable[string]{
					Set: true,
					Value: func() *string {
						v := "+14155552671"
						return &v
					}(),
				},
			},
			wantErr: nil,
			wantUser: func() User {
				user := validUser
				user.Version = 2
				user.PhoneNumber = "+14155552671"
				return user
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := validUser

			err := user.ApplyPatch(&tt.patch)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantUser, user)
		})
	}
}

func TestUserApplyPatchInvalidPatchDoesNotMutateUser(t *testing.T) {
	now := time.Now()

	user := User{
		ID:          1,
		Version:     1,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   now.AddDate(-20, 0, 0),
		Description: "Original description",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
	}

	original := User{
		ID:          1,
		Version:     1,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   now.AddDate(-20, 0, 0),
		Description: "Original description",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
	}

	patch := UserPatch{
		Version: 2,
		Name: Nullable[string]{
			Set:   true,
			Value: func() *string { v := ""; return &v }(),
		},
	}

	err := user.ApplyPatch(&patch)

	if !errors.Is(err, core_errors.ErrInvalidUserName) {
		t.Fatalf(
			"ApplyPatch() = %v, want %v",
			err,
			core_errors.ErrInvalidUserName,
		)
	}

	if !reflect.DeepEqual(user, original) {
		t.Fatalf("user was mutated after invalid patch")
	}

	if !strings.Contains(err.Error(), "validate user patch") {
		t.Fatalf("error = %q, want patch validation context", err)
	}
}

func TestUserApplyPatchInvalidResultDoesNotMutateUser(t *testing.T) {
	now := time.Now()

	createdAt := now.AddDate(-20, 0, 0)
	updatedAt := createdAt.Add(time.Hour)

	user := User{
		ID:          1,
		Version:     1,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   now.AddDate(-34, 0, 0),
		Description: "Original description",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	original := User{
		ID:          1,
		Version:     1,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   now.AddDate(-34, 0, 0),
		Description: "Original description",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// The patch itself is valid, but the resulting User is invalid:
	// the new birth date is after CreatedAt.
	newBirthDate := user.CreatedAt.Add(time.Hour)

	patch := UserPatch{
		Version: 1,
		BirthDate: Nullable[time.Time]{
			Set:   true,
			Value: &newBirthDate,
		},
	}

	err := user.ApplyPatch(&patch)

	if !errors.Is(err, core_errors.ErrInvalidUserCreatedAt) {
		t.Fatalf(
			"ApplyPatch() = %v, want %v",
			err,
			core_errors.ErrInvalidUserCreatedAt,
		)
	}

	if !strings.Contains(err.Error(), "validate patched user") {
		t.Fatalf("error = %q, want patched-user validation context", err)
	}

	if !reflect.DeepEqual(user, original) {
		t.Fatalf("user was mutated after invalid resulting state")
	}
}

func TestUserApplyPatchNilPatch(t *testing.T) {
	now := time.Now()

	user := User{
		ID:          1,
		Version:     1,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   now.AddDate(-20, 0, 0),
		Description: "Original description",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-time.Hour),
	}

	defer func() {
		if recover() == nil {
			t.Fatal("ApplyPatch(nil) did not panic")
		}
	}()

	_ = user.ApplyPatch(nil)
}

func TestNewUser(t *testing.T) {
	birthDate := time.Date(2005, time.January, 15, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)
	deletedAt := time.Date(2026, time.January, 3, 12, 0, 0, 0, time.UTC)

	got := NewUser(
		42,
		7,
		"Alexey",
		"Bobkov",
		"alexey",
		birthDate,
		"Software engineer",
		"alexey@example.com",
		"+14155552671",
		createdAt,
		updatedAt,
		&deletedAt,
	)

	want := &User{
		ID:          42,
		Version:     7,
		Name:        "Alexey",
		Surname:     "Bobkov",
		Username:    "alexey",
		BirthDate:   birthDate,
		Description: "Software engineer",
		Email:       "alexey@example.com",
		PhoneNumber: "+14155552671",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   &deletedAt,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"NewUser() = %#v, want %#v",
			got,
			want,
		)
	}
}
