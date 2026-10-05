package domain

import (
	"errors"
	"net/mail"
	"strings"
	"testing"
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateID(t *testing.T) {
	maxInt := int(^uint(0) >> 1)

	tests := []struct {
		name string
		id   int
		want error
	}{
		{
			name: "minimum valid",
			id:   1,
			want: nil,
		},
		{
			name: "normal positive",
			id:   123,
			want: nil,
		},
		{
			name: "maximum int",
			id:   maxInt,
			want: nil,
		},
		{
			name: "zero",
			id:   0,
			want: core_errors.ErrInvalidUserID,
		},
		{
			name: "negative one",
			id:   -1,
			want: core_errors.ErrInvalidUserID,
		},
		{
			name: "minimum int",
			id:   -maxInt - 1,
			want: core_errors.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id)

			if !errors.Is(err, tt.want) {
				t.Fatalf("validateID(%d) = %v, want %v", tt.id, err, tt.want)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	maxInt := int(^uint(0) >> 1)

	tests := []struct {
		name    string
		version int
		want    error
	}{
		{
			name:    "minimum valid",
			version: 1,
			want:    nil,
		},
		{
			name:    "normal version",
			version: 42,
			want:    nil,
		},
		{
			name:    "maximum int",
			version: maxInt,
			want:    nil,
		},
		{
			name:    "zero",
			version: 0,
			want:    core_errors.ErrInvalidUserVersion,
		},
		{
			name:    "negative one",
			version: -1,
			want:    core_errors.ErrInvalidUserVersion,
		},
		{
			name:    "minimum int",
			version: -maxInt - 1,
			want:    core_errors.ErrInvalidUserVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersion(tt.version)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateVersion(%d) = %v, want %v",
					tt.version,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "single letter",
			input: "A",
			want:  nil,
		},
		{
			name:  "normal name",
			input: "Alexey",
			want:  nil,
		},
		{
			name:  "Cyrillic",
			input: "Алексей",
			want:  nil,
		},
		{
			name:  "Unicode letters",
			input: "José",
			want:  nil,
		},
		{
			name:  "exactly 50 runes",
			input: strings.Repeat("A", 50),
			want:  nil,
		},
		{
			name:  "51 runes",
			input: strings.Repeat("A", 51),
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "empty",
			input: "",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "spaces only",
			input: "   ",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "tabs only",
			input: "\t\t",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "newline only",
			input: "\n",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "leading space",
			input: " Alex",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "trailing space",
			input: "Alex ",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "internal space",
			input: "Alex Bob",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "hyphen",
			input: "Alex-Bob",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "apostrophe",
			input: "O'Connor",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "digit",
			input: "Alex1",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "punctuation",
			input: "Alex!",
			want:  core_errors.ErrInvalidUserName,
		},
		{
			name:  "emoji",
			input: "Alex😀",
			want:  core_errors.ErrInvalidUserName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.input)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateName(%q) = %v, want %v",
					tt.input,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateSurname(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "single letter",
			input: "A",
			want:  nil,
		},
		{
			name:  "normal surname",
			input: "Bobkov",
			want:  nil,
		},
		{
			name:  "Cyrillic",
			input: "Бобков",
			want:  nil,
		},
		{
			name:  "Unicode letters",
			input: "García",
			want:  nil,
		},
		{
			name:  "exactly 50 runes",
			input: strings.Repeat("A", 50),
			want:  nil,
		},
		{
			name:  "51 runes",
			input: strings.Repeat("A", 51),
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "empty",
			input: "",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "spaces only",
			input: "   ",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "leading space",
			input: " Bobkov",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "trailing space",
			input: "Bobkov ",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "internal space",
			input: "Van Dyke",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "hyphen",
			input: "Smith-Jones",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "apostrophe",
			input: "O'Connor",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "digit",
			input: "Bobkov1",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "punctuation",
			input: "Bobkov!",
			want:  core_errors.ErrInvalidUserSurname,
		},
		{
			name:  "emoji",
			input: "Bobkov😀",
			want:  core_errors.ErrInvalidUserSurname,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSurname(tt.input)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateSurname(%q) = %v, want %v",
					tt.input,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateUserName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "single character",
			input: "a",
			want:  nil,
		},
		{
			name:  "normal username",
			input: "alexeybobkov",
			want:  nil,
		},
		{
			name:  "digits",
			input: "alexey123",
			want:  nil,
		},
		{
			name:  "underscore",
			input: "alex_bobkov",
			want:  nil,
		},
		{
			name:  "hyphen",
			input: "alex-bobkov",
			want:  nil,
		},
		{
			name:  "space",
			input: "alex bobkov",
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "spaces",
			input: "alex bobkov username",
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "more than 1 space sequently",
			input: "alex   bobkov",
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "Unicode",
			input: "алексей",
			want:  nil,
		},
		{
			name:  "exactly 100 runes",
			input: strings.Repeat("a", 100),
			want:  nil,
		},
		{
			name:  "101 runes",
			input: strings.Repeat("a", 101),
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "empty",
			input: "",
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "spaces only",
			input: "   ",
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "tabs only",
			input: "\t\t",
			want:  core_errors.ErrInvalidUserUsername,
		},
		{
			name:  "newlines only",
			input: "\n\n",
			want:  core_errors.ErrInvalidUserUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserName(tt.input)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateUserName(%q) = %v, want %v",
					tt.input,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateBirthDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		birthDate time.Time
		want      error
	}{
		{
			name:      "exactly 14 years ago",
			birthDate: now.AddDate(-14, 0, 0),
			want:      nil,
		},
		{
			name:      "20 years ago",
			birthDate: now.AddDate(-20, 0, 0),
			want:      nil,
		},
		{
			name:      "exactly 150 years ago",
			birthDate: now.AddDate(-150, 0, 0),
			want:      nil,
		},
		{
			name:      "younger than 14",
			birthDate: now.AddDate(-13, 0, 0),
			want:      core_errors.ErrInvalidUserBirthdate,
		},
		{
			name:      "older than 150",
			birthDate: now.AddDate(-151, 0, 0),
			want:      core_errors.ErrInvalidUserBirthdate,
		},
		{
			name:      "future",
			birthDate: now.Add(time.Hour),
			want:      core_errors.ErrInvalidUserBirthdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBirthDate(tt.birthDate)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateBirthDate(%v) = %v, want %v",
					tt.birthDate,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        error
	}{
		{
			name:        "empty",
			description: "",
			want:        nil,
		},
		{
			name:        "normal description",
			description: "Software engineer interested in distributed systems.",
			want:        nil,
		},
		{
			name:        "exactly 1000 ASCII runes",
			description: strings.Repeat("a", 1000),
			want:        nil,
		},
		{
			name:        "1001 ASCII runes",
			description: strings.Repeat("a", 1001),
			want:        core_errors.ErrInvalidUserDescription,
		},
		{
			name:        "exactly 1000 Unicode runes",
			description: strings.Repeat("Ж", 1000),
			want:        nil,
		},
		{
			name:        "1001 Unicode runes",
			description: strings.Repeat("Ж", 1001),
			want:        core_errors.ErrInvalidUserDescription,
		},
		{
			name:        "newline",
			description: "hello\nworld",
			want:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateDescription() = %v, want %v",
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  error
	}{
		{
			name:  "normal email",
			email: "user@example.com",
			want:  nil,
		},
		{
			name:  "subdomain",
			email: "user@mail.example.com",
			want:  nil,
		},
		{
			name:  "plus tag",
			email: "user+tag@example.com",
			want:  nil,
		},
		{
			name:  "hyphen",
			email: "first-last@example.com",
			want:  nil,
		},
		{
			name:  "underscore",
			email: "first_last@example.com",
			want:  nil,
		},
		{
			name:  "empty",
			email: "",
			want:  mail.ErrHeaderNotPresent,
		},
		{
			name:  "missing at",
			email: "userexample.com",
			want:  mail.ErrHeaderNotPresent,
		},
		{
			name:  "missing local part",
			email: "@example.com",
			want:  mail.ErrHeaderNotPresent,
		},
		{
			name:  "missing domain",
			email: "user@",
			want:  mail.ErrHeaderNotPresent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)

			if tt.want == nil {
				if err != nil {
					t.Fatalf(
						"validateEmail(%q) = %v, want nil",
						tt.email,
						err,
					)
				}
				return
			}

			if err == nil {
				t.Fatalf(
					"validateEmail(%q) = nil, want an error",
					tt.email,
				)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		want        error
	}{
		{
			name:        "valid US number",
			phoneNumber: "+14155552671",
			want:        nil,
		},
		{
			name:        "valid UK number",
			phoneNumber: "+442071838750",
			want:        nil,
		},
		{
			name:        "valid French number",
			phoneNumber: "+33142345678",
			want:        nil,
		},
		{
			name:        "empty",
			phoneNumber: "",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name:        "missing plus",
			phoneNumber: "14155552671",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name:        "contains spaces",
			phoneNumber: "+1 415 555 2671",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name:        "contains hyphens",
			phoneNumber: "+1-415-555-2671",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name:        "contains letters",
			phoneNumber: "+1415555267a",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name:        "invalid country code",
			phoneNumber: "+999123456789",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
		{
			name:        "too short",
			phoneNumber: "+1",
			want:        core_errors.ErrInvalidUserPhoneNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePhoneNumber(tt.phoneNumber)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validatePhoneNumber(%q) = %v, want %v",
					tt.phoneNumber,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     error
	}{
		{
			name:     "minimum valid password",
			password: "Aa1!aaaa",
			want:     nil,
		},
		{
			name:     "long valid password",
			password: "VerySecurePassword123!",
			want:     nil,
		},
		{
			name:     "exactly 8 characters",
			password: "Aa12345!",
			want:     nil,
		},
		{
			name:     "7 characters",
			password: "Aa1234!",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "no digit",
			password: "Password!",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "no letter",
			password: "12345678!",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "no punctuation",
			password: "Password123",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "no uppercase",
			password: "password123!",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "no lowercase",
			password: "PASSWORD123!",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "empty",
			password: "",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "spaces only",
			password: "        ",
			want:     core_errors.ErrInvalidUserPassword,
		},
		{
			name:     "valid with leading space",
			password: " Password1!",
			want:     nil,
		},
		{
			name:     "valid with trailing space",
			password: "Password1! ",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validatePassword(%q) = %v, want %v",
					tt.password,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateCreatedAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		createdAt time.Time
		want      error
	}{
		{
			name:      "past",
			createdAt: now.Add(-time.Hour),
			want:      nil,
		},
		{
			name:      "now",
			createdAt: now,
			want:      nil,
		},
		{
			name:      "future",
			createdAt: now.Add(time.Hour),
			want:      core_errors.ErrInvalidUserCreatedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreatedAt(tt.createdAt)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateCreatedAt(%v) = %v, want %v",
					tt.createdAt,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateUpdatedAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		updatedAt time.Time
		want      error
	}{
		{
			name:      "past",
			updatedAt: now.Add(-time.Hour),
			want:      nil,
		},
		{
			name:      "now",
			updatedAt: now,
			want:      nil,
		},
		{
			name:      "future",
			updatedAt: now.Add(time.Hour),
			want:      core_errors.ErrInvalidUserUpdatedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdatedAt(tt.updatedAt)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateUpdatedAt(%v) = %v, want %v",
					tt.updatedAt,
					err,
					tt.want,
				)
			}
		})
	}
}

func TestValidateDeletedAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		deletedAt *time.Time
		want      error
	}{
		{
			name:      "nil",
			deletedAt: nil,
			want:      nil,
		},
		{
			name: "past",
			deletedAt: func() *time.Time {
				t := now.Add(-time.Hour)
				return &t
			}(),
			want: nil,
		},
		{
			name: "now",
			deletedAt: func() *time.Time {
				t := now
				return &t
			}(),
			want: nil,
		},
		{
			name: "future",
			deletedAt: func() *time.Time {
				t := now.Add(time.Hour)
				return &t
			}(),
			want: core_errors.ErrInvalidUserDeletedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeletedAt(tt.deletedAt)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateDeletedAt(%v) = %v, want %v",
					tt.deletedAt,
					err,
					tt.want,
				)
			}
		})
	}
}

// -----------------------------------------------------------------------------
//
// Date ordering
// -----------------------------------------------------------------------------
//
// Contract:
//
// birthDate < createdAt <= updatedAt
//
// If deletedAt != nil:
//
// birthDate < createdAt <= updatedAt <= deletedAt
//
// -----------------------------------------------------------------------------

func TestValidateDatesOrder(t *testing.T) {
	birthDate := time.Date(
		2000,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	createdAt := time.Date(
		2025,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	updatedAt := time.Date(
		2026,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	deletedAt := time.Date(
		2027,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	tests := []struct {
		name      string
		birthDate time.Time
		createdAt time.Time
		updatedAt time.Time
		deletedAt *time.Time
		want      error
	}{
		{
			name:      "valid without deletion",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: nil,
			want:      nil,
		},
		{
			name:      "valid with deletion",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: &deletedAt,
			want:      nil,
		},
		{
			name:      "created one nanosecond after birth",
			birthDate: birthDate,
			createdAt: birthDate.Add(time.Nanosecond),
			updatedAt: birthDate.Add(2 * time.Nanosecond),
			deletedAt: nil,
			want:      nil,
		},
		{
			name:      "created equal birth",
			birthDate: birthDate,
			createdAt: birthDate,
			updatedAt: updatedAt,
			deletedAt: nil,
			want:      core_errors.ErrInvalidUserCreatedAt,
		},
		{
			name:      "created before birth",
			birthDate: birthDate,
			createdAt: birthDate.Add(-time.Nanosecond),
			updatedAt: updatedAt,
			deletedAt: nil,
			want:      core_errors.ErrInvalidUserCreatedAt,
		},
		{
			name:      "updated after created",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: createdAt.Add(time.Nanosecond),
			deletedAt: nil,
			want:      nil,
		},
		{
			name:      "updated equal created",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: createdAt,
			deletedAt: nil,
			want:      nil,
		},
		{
			name:      "updated before created",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: createdAt.Add(-time.Nanosecond),
			deletedAt: nil,
			want:      core_errors.ErrInvalidUserUpdatedAt,
		},
		{
			name:      "deleted after updated",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: func() *time.Time {
				t := updatedAt.Add(time.Nanosecond)
				return &t
			}(),
			want: nil,
		},
		{
			name:      "deleted equal updated",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: func() *time.Time {
				t := updatedAt
				return &t
			}(),
			want: nil,
		},
		{
			name:      "deleted before updated",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: func() *time.Time {
				t := updatedAt.Add(-time.Nanosecond)
				return &t
			}(),
			want: core_errors.ErrInvalidUserDeletedAt,
		},
		{
			name:      "deleted equal created with later updated",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: func() *time.Time {
				t := createdAt
				return &t
			}(),
			want: core_errors.ErrInvalidUserDeletedAt,
		},
		{
			name:      "created equal updated equal deleted",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: createdAt,
			deletedAt: func() *time.Time {
				t := createdAt
				return &t
			}(),
			want: nil,
		},
		{
			name:      "nil deletedAt",
			birthDate: birthDate,
			createdAt: createdAt,
			updatedAt: updatedAt,
			deletedAt: nil,
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDatesOrder(
				tt.birthDate,
				tt.createdAt,
				tt.updatedAt,
				tt.deletedAt,
			)

			if !errors.Is(err, tt.want) {
				t.Fatalf(
					"validateDatesOrder() = %v, want %v",
					err,
					tt.want,
				)
			}
		})
	}
}
