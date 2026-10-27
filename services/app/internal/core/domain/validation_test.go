package domain

import (
	"math"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want error
	}{
		{name: "minimum invalid", id: 0, want: core_errors.ErrInvalidID},
		{name: "negative", id: -1, want: core_errors.ErrInvalidID},
		{name: "minimum integer", id: math.MinInt, want: core_errors.ErrInvalidID},
		{name: "one", id: 1},
		{name: "positive", id: 42},
		{name: "maximum integer", id: math.MaxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version int
		want    error
	}{
		{name: "zero", version: 0, want: core_errors.ErrInvalidVersion},
		{name: "negative", version: -1, want: core_errors.ErrInvalidVersion},
		{name: "minimum integer", version: math.MinInt, want: core_errors.ErrInvalidVersion},
		{name: "one", version: 1},
		{name: "positive", version: 42},
		{name: "maximum integer", version: math.MaxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersion(tt.version)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
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
			err := validateName(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
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
			err := validateSurname(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
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
			err := validateUserName(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateBirthDate(t *testing.T) {
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
			err := validateBirthDate(tt.date)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
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
			err := validateDescription(tt.text)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
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
			err := validateEmail(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
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
			err := validatePhoneNumber(tt.input)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
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
			err := validatePassword(tt.password)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCreatedAt(t *testing.T) {
	now := time.Now()
	require.NoError(t, validateCreatedAt(now.Add(-time.Second)))
	require.NoError(t, validateCreatedAt(now))
	require.ErrorIs(t, validateCreatedAt(now.Add(time.Second)), core_errors.ErrInvalidCreatedAt)
	require.NoError(t, validateCreatedAt(time.Time{}))
}

func TestValidateUpdatedAt(t *testing.T) {
	now := time.Now()
	require.NoError(t, validateUpdatedAt(now.Add(-time.Second)))
	require.NoError(t, validateUpdatedAt(now))
	require.ErrorIs(t, validateUpdatedAt(now.Add(time.Second)), core_errors.ErrInvalidUpdatedAt)
	require.NoError(t, validateUpdatedAt(time.Time{}))
}

func TestValidateDeletedAt(t *testing.T) {
	now := time.Now()
	require.NoError(t, validateDeletedAt(nil))
	require.NoError(t, validateDeletedAt(timePtr(now.Add(-time.Second))))
	require.NoError(t, validateDeletedAt(timePtr(now)))
	require.ErrorIs(t, validateDeletedAt(timePtr(now.Add(time.Second))), core_errors.ErrInvalidDeletedAt)
	require.NoError(t, validateDeletedAt(timePtr(time.Time{})))
}

func TestValidateSystemDatesOrder(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		created time.Time
		updated time.Time
		deleted *time.Time
		want    error
	}{
		{name: "strict chronological order", created: now, updated: now.Add(time.Hour)},
		{name: "created equals updated", created: now, updated: now},
		{name: "deleted later than updated", created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(2 * time.Hour))},
		{name: "deleted equals updated", created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(time.Hour))},
		{name: "created after updated", created: now, updated: now.Add(-time.Second), want: core_errors.ErrInvalidUpdatedAt},
		{name: "updated after deleted", created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(time.Minute)), want: core_errors.ErrInvalidDeletedAt},
		{name: "no deleted timestamp", created: now, updated: now, deleted: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSystemDatesOrder(tt.created, tt.updated, tt.deleted)
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

func TestValidateStartDate(t *testing.T) {
	require.ErrorIs(t, validateStartDate(time.Time{}), core_errors.ErrInvalidTripStartDate)
	require.NoError(t, validateStartDate(time.Now()))
	require.NoError(t, validateStartDate(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestValidateEndDate(t *testing.T) {
	require.ErrorIs(t, validateEndDate(time.Time{}), core_errors.ErrInvalidTripEndDate)
	require.NoError(t, validateEndDate(time.Now()))
	require.NoError(t, validateEndDate(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestValidateTripDatesOrder(t *testing.T) {
	created := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name                             string
		start, end, createdAt, updatedAt time.Time
		deleted                          *time.Time
		want                             error
	}{
		{name: "valid future trip", start: created.Add(time.Hour), end: created.Add(2 * time.Hour), createdAt: created, updatedAt: created},
		{name: "start before end", start: created.Add(time.Hour), end: created.Add(time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripEndDate},
		{name: "start after end", start: created.Add(2 * time.Hour), end: created.Add(time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripEndDate},
		{name: "start before creation", start: created.Add(-time.Hour), end: created.Add(time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripStartDate},
		{name: "start equals creation", start: created, end: created.Add(time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripStartDate},
		{name: "system dates invalid", start: created.Add(time.Hour), end: created.Add(2 * time.Hour), createdAt: created, updatedAt: created.Add(-time.Second), want: core_errors.ErrInvalidUpdatedAt},
		{name: "deleted before updated", start: created.Add(time.Hour), end: created.Add(2 * time.Hour), createdAt: created, updatedAt: created.Add(time.Hour), deleted: timePtr(created.Add(30 * time.Minute)), want: core_errors.ErrInvalidDeletedAt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTripDatesOrder(tt.start, tt.end, tt.createdAt, tt.updatedAt, tt.deleted)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateBudget(t *testing.T) {
	require.NoError(t, validateBudget(nil))
	require.NoError(t, validateBudget(&Money{AmountMinorUnits: 1, Currency: "RUB"}))
	require.ErrorIs(t, validateBudget(&Money{AmountMinorUnits: 0, Currency: "RUB"}), core_errors.ErrInvalidBudget)
	require.ErrorIs(t, validateBudget(&Money{AmountMinorUnits: 100, Currency: "USD"}), core_errors.ErrUnsupportedCurrency)
}

func TestValidateTripName(t *testing.T) {
	for _, value := range []string{"Trip", "Поездка", "Trip 2026", "   Trip   "} {
		require.NoError(t, validateTripName(value), "value: %q", value)
	}
	for _, value := range []string{"", " ", "\t\n"} {
		require.ErrorIs(t, validateTripName(value), core_errors.ErrInvalidTripName, "value: %q", value)
	}
}

func TestValidatePlanningArea(t *testing.T) {
	for _, value := range []string{"Nizhny Novgorod", "Paris, France", "Москва"} {
		require.NoError(t, validatePlanningArea(value), "value: %q", value)
	}
	for _, value := range []string{"", " ", "\t\n"} {
		require.ErrorIs(t, validatePlanningArea(value), core_errors.ErrInvalidPlanningArea, "value: %q", value)
	}
}

func TestValidateStartLocation(t *testing.T) {
	require.NoError(t, validateStartLocation(nil))
	require.NoError(t, validateStartLocation(&Location{Name: "Moscow", Latitude: 55.75, Longitude: 37.62}))
	require.ErrorIs(t, validateStartLocation(&Location{Name: "", Latitude: 55.75, Longitude: 37.62}), core_errors.ErrInvalidLocationName)
	require.ErrorIs(t, validateStartLocation(&Location{Name: "Moscow", Latitude: 91, Longitude: 37.62}), core_errors.ErrInvalidLatitude)
	require.ErrorIs(t, validateStartLocation(&Location{Name: "Moscow", Latitude: 55.75, Longitude: 181}), core_errors.ErrInvalidLongitude)
}

func TestValidateInterests(t *testing.T) {
	tests := []struct {
		name   string
		values []Interest
		want   error
	}{
		{name: "nil list"},
		{name: "empty list", values: []Interest{}},
		{name: "single supported", values: []Interest{InterestHistory}},
		{name: "multiple supported", values: []Interest{InterestHistory, InterestNature, InterestFood, InterestArchitecture}},
		{name: "blank interest", values: []Interest{" "}, want: core_errors.ErrInvalidInterest},
		{name: "unsupported interest", values: []Interest{"shopping"}, want: core_errors.ErrUnsupportedInterest},
		{name: "supported then unsupported", values: []Interest{InterestFood, "shopping"}, want: core_errors.ErrUnsupportedInterest},
		{name: "blank precedes unsupported", values: []Interest{" ", "shopping"}, want: core_errors.ErrInvalidInterest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInterests(tt.values)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateTimeZone(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  error
	}{
		{name: "supported Moscow", value: "Europe/Moscow"},
		{name: "empty", value: "", want: core_errors.ErrInvalidTimeZone},
		{name: "whitespace", value: " \t ", want: core_errors.ErrInvalidTimeZone},
		{name: "invalid IANA zone", value: "Mars/Olympus", want: core_errors.ErrUnsupportedTimezone},
		{name: "valid UTC but unsupported by product policy", value: "UTC", want: core_errors.ErrUnsupportedTimezone},
		{name: "valid New York but unsupported by product policy", value: "America/New_York", want: core_errors.ErrUnsupportedTimezone},
		{name: "Moscow with whitespace currently rejected", value: " Europe/Moscow ", want: core_errors.ErrUnsupportedTimezone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTimeZone(tt.value)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidationInputsDoNotNeedASCII(t *testing.T) {
	// Sanity check that Unicode cases above are actually multi-byte input;
	// the validators count runes for their length limits.
	assert.Greater(t, len("Алексей"), utf8.RuneCountInString("Алексей"))
	assert.NoError(t, validateName("Алексей"))
	assert.NoError(t, validateSurname("Бобков"))
	assert.NoError(t, validateDescription(strings.Repeat("界", 1000)))
}

func timePtr(value time.Time) *time.Time { return &value }
