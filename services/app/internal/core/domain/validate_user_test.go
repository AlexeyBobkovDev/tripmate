package domain

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateUserName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "ASCII letters", input: "Alexey"},
		{name: "Unicode letters", input: "Алексей"},
		{name: "Chinese letters", input: "李"},
		{name: "one rune", input: "A"},
		{name: "exactly 50 runes", input: strings.Repeat("a", 50)},
		{name: "empty", input: "", want: core_errors.ErrInvalidUserName},
		{name: "whitespace only", input: " \t\n", want: core_errors.ErrInvalidUserName},
		{name: "51 runes", input: strings.Repeat("a", 51), want: core_errors.ErrInvalidUserName},
		{name: "digit", input: "Alex1", want: core_errors.ErrInvalidUserName},
		{name: "hyphen", input: "Alex-ey", want: core_errors.ErrInvalidUserName},
		{name: "apostrophe", input: "O'Neil", want: core_errors.ErrInvalidUserName},
		{name: "leading whitespace", input: " Alex", want: core_errors.ErrInvalidUserName},
		{name: "trailing whitespace", input: "Alex ", want: core_errors.ErrInvalidUserName},
		{name: "emoji", input: "Alex🙂", want: core_errors.ErrInvalidUserName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserName(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserSurname(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "ASCII letters", input: "Bobkov"},
		{name: "Unicode letters", input: "Бобков"},
		{name: "one rune", input: "B"},
		{name: "exactly 50 runes", input: strings.Repeat("Б", 50)},
		{name: "empty", input: "", want: core_errors.ErrInvalidUserSurname},
		{name: "whitespace only", input: " \t", want: core_errors.ErrInvalidUserSurname},
		{name: "51 runes", input: strings.Repeat("Б", 51), want: core_errors.ErrInvalidUserSurname},
		{name: "digit", input: "Bobkov2", want: core_errors.ErrInvalidUserSurname},
		{name: "hyphen", input: "Smith-Jones", want: core_errors.ErrInvalidUserSurname},
		{name: "apostrophe", input: "O'Neil", want: core_errors.ErrInvalidUserSurname},
		{name: "leading whitespace", input: " Bobkov", want: core_errors.ErrInvalidUserSurname},
		{name: "trailing whitespace", input: "Bobkov ", want: core_errors.ErrInvalidUserSurname},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserSurname(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserUsername(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "simple", input: "alexey"},
		{name: "underscore and digits", input: "alex_123"},
		{name: "Unicode username", input: "алексей"},
		{name: "exactly 100 runes", input: strings.Repeat("a", 100)},
		{name: "empty", input: "", want: core_errors.ErrInvalidUserUsername},
		{name: "whitespace only", input: " \t\n", want: core_errors.ErrInvalidUserUsername},
		{name: "101 runes", input: strings.Repeat("a", 101), want: core_errors.ErrInvalidUserUsername},
		{name: "ASCII space", input: "alex ey", want: core_errors.ErrInvalidUserUsername},
		{name: "leading ASCII space", input: " alex", want: core_errors.ErrInvalidUserUsername},
		{name: "trailing ASCII space", input: "alex ", want: core_errors.ErrInvalidUserUsername},
		{name: "tab inside username currently accepted", input: "alex\tey"},
		{name: "newline inside username currently accepted", input: "alex\ney"},
		{name: "non-ASCII whitespace currently accepted", input: "alex\u00a0ey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserUsername(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserBirthDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		date time.Time
		want error
	}{
		{name: "exactly 14 years old", date: now.AddDate(-14, 0, 0)},
		{name: "adult", date: now.AddDate(-25, 0, 0)},
		{name: "exactly 150 years old", date: now.AddDate(-150, 0, 0)},
		{name: "13 years old", date: now.AddDate(-13, 0, 0), want: core_errors.ErrInvalidUserBirthdate},
		{name: "151 years old", date: now.AddDate(-151, 0, 0), want: core_errors.ErrInvalidUserBirthdate},
		{name: "tomorrow", date: now.Add(24 * time.Hour), want: core_errors.ErrInvalidUserBirthdate},
		{name: "later today", date: now.Add(time.Minute), want: core_errors.ErrInvalidUserBirthdate},
		{name: "zero date", date: time.Time{}, want: core_errors.ErrInvalidUserBirthdate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserBirthDate(tt.date)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserDescription(t *testing.T) {
	tests := []struct {
		name string
		text string
		want error
	}{
		{name: "empty allowed", text: ""},
		{name: "one rune", text: "a"},
		{name: "exactly 1000 ASCII runes", text: strings.Repeat("a", 1000)},
		{name: "1001 ASCII runes", text: strings.Repeat("a", 1001), want: core_errors.ErrInvalidUserDescription},
		{name: "exactly 1000 Unicode runes", text: strings.Repeat("界", 1000)},
		{name: "1001 Unicode runes", text: strings.Repeat("界", 1001), want: core_errors.ErrInvalidUserDescription},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, utf8.RuneCountInString(tt.text), len([]rune(tt.text)))
			err := validateUserDescription(tt.text)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "simple email", input: "alex@example.com"},
		{name: "subdomain", input: "alex@mail.example.co.uk"},
		{name: "plus addressing", input: "alex+tripmate@example.com"},
		{name: "quoted local part", input: `"alex bob"@example.com`},
		{name: "display name currently accepted by mail.ParseAddress", input: "Alex Bob <alex@example.com>"},
		{name: "empty", input: "", want: core_errors.ErrInvalidUserEmail},
		{name: "missing at", input: "alex.example.com", want: core_errors.ErrInvalidUserEmail},
		{name: "missing domain", input: "alex@", want: core_errors.ErrInvalidUserEmail},
		{name: "missing local part", input: "@example.com", want: core_errors.ErrInvalidUserEmail},
		{name: "spaces", input: "alex @example.com", want: core_errors.ErrInvalidUserEmail},
		{name: "invalid syntax", input: "alex@@example.com", want: core_errors.ErrInvalidUserEmail},
		{name: "newline", input: "alex@example.com\nBcc:other@example.com", want: core_errors.ErrInvalidUserEmail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserEmail(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserPhoneNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "valid E164 US number", input: "+14155552671"},
		{name: "valid E164 UK number", input: "+442079460018"},
		{name: "valid E164 Russian number", input: "+74951234567"},
		{name: "empty", input: "", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "too short", input: "+123", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "missing plus", input: "14155552671", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "spaces", input: "+1 415 555 2671", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "parentheses", input: "+1 (415) 555-2671", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "national format", input: "02079460018", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "alphabetic", input: "+1ABC5552671", want: core_errors.ErrInvalidUserPhoneNumber},
		{name: "too many digits", input: "+12345678901234567890", want: core_errors.ErrInvalidUserPhoneNumber},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserPhoneNumber(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     error
	}{
		{name: "valid password", password: "Abcd123!"},
		{name: "long valid password", password: "TripMate-StrongPassword-123!"},
		{name: "unicode letters and digit", password: "Éclair123!"},
		{name: "too short", password: "Ab1!xyz", want: core_errors.ErrInvalidUserPassword},
		{name: "missing digit", password: "Abcdefg!", want: core_errors.ErrInvalidUserPassword},
		{name: "missing letter", password: "1234567!", want: core_errors.ErrInvalidUserPassword},
		{name: "missing uppercase", password: "abcd123!", want: core_errors.ErrInvalidUserPassword},
		{name: "missing lowercase", password: "ABCD123!", want: core_errors.ErrInvalidUserPassword},
		{name: "missing punctuation", password: "Abcdef12", want: core_errors.ErrInvalidUserPassword},
		{name: "symbol is not punctuation", password: "Abcd123 ", want: core_errors.ErrInvalidUserPassword},
		{name: "only symbols/digits", password: "!@#$1234", want: core_errors.ErrInvalidUserPassword},
		{name: "whitespace-trimmed length too short", password: "  Ab1!  ", want: core_errors.ErrInvalidUserPassword},
		{name: "Unicode punctuation", password: "Abcd123。"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserPassword(tt.password)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateUserDatesOrder(t *testing.T) {
	now := time.Now()
	birth := now.AddDate(-20, 0, 0)
	tests := []struct {
		name                    string
		birth, created, updated time.Time
		deleted                 *time.Time
		want                    error
	}{
		{name: "valid", birth: birth, created: now, updated: now},
		{name: "birth before creation", birth: birth, created: now, updated: now.Add(time.Hour)},
		{name: "birth equals creation", birth: now, created: now, updated: now, want: core_errors.ErrInvalidCreatedAt},
		{name: "birth after creation", birth: now.Add(time.Hour), created: now, updated: now, want: core_errors.ErrInvalidCreatedAt},
		{name: "created after updated", birth: birth, created: now, updated: now.Add(-time.Second), want: core_errors.ErrInvalidUpdatedAt},
		{name: "updated after deletion", birth: birth, created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(time.Minute)), want: core_errors.ErrInvalidDeletedAt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserDatesOrder(tt.birth, tt.created, tt.updated, tt.deleted)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
