package domain

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name      string
		birthDate time.Time
		now       time.Time
		want      int
	}{
		{
			name:      "birthday already occurred this year",
			birthDate: time.Date(2000, time.January, 15, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
			want:      26,
		},
		{
			name:      "birthday has not occurred yet this year",
			birthDate: time.Date(2000, time.June, 15, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
			want:      25,
		},
		{
			name:      "birthday today",
			birthDate: time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.February, 1, 12, 0, 0, 0, time.UTC),
			want:      26,
		},
		{
			name:      "day before birthday",
			birthDate: time.Date(2000, time.February, 2, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.February, 1, 23, 59, 59, 0, time.UTC),
			want:      25,
		},
		{
			name:      "day after birthday",
			birthDate: time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.February, 2, 0, 0, 0, 0, time.UTC),
			want:      26,
		},
		{
			name:      "exactly one year old",
			birthDate: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			want:      1,
		},
		{
			name:      "less than one year old",
			birthDate: time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			want:      0,
		},
		{
			name:      "same birth and current date",
			birthDate: time.Date(2026, time.March, 10, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.March, 10, 0, 0, 0, 0, time.UTC),
			want:      0,
		},
		{
			name:      "leap day birthday before birthday",
			birthDate: time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
			want:      24,
		},
		{
			name:      "leap day birthday in leap year",
			birthDate: time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			want:      24,
		},
		{
			name:      "leap day birthday after February",
			birthDate: time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
			want:      25,
		},
		{
			name:      "birthday month is later",
			birthDate: time.Date(1990, time.December, 31, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			want:      35,
		},
		{
			name:      "birthday month has passed",
			birthDate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC),
			want:      36,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAge(tt.birthDate, tt.now)

			if got != tt.want {
				t.Fatalf(
					"calculateAge(%v, %v) = %d, want %d",
					tt.birthDate,
					tt.now,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCalculateCurrentAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		birthDate time.Time
	}{
		{
			name:      "adult",
			birthDate: now.AddDate(-20, 0, 0),
		},
		{
			name:      "14 years old",
			birthDate: now.AddDate(-14, 0, 0),
		},
		{
			name:      "150 years old",
			birthDate: now.AddDate(-150, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateCurrentAge(tt.birthDate)

			// calculateCurrentAge uses time.Now() internally, so allow
			// only the current calendar-age result.
			want := time.Now().Year() - tt.birthDate.Year()

			birthdayThisYear := time.Date(
				time.Now().Year(),
				tt.birthDate.Month(),
				tt.birthDate.Day(),
				0,
				0,
				0,
				0,
				time.Now().Location(),
			)

			if time.Now().Before(birthdayThisYear) {
				want--
			}

			if got != want {
				t.Fatalf(
					"calculateCurrentAge(%v) = %d, want %d",
					tt.birthDate,
					got,
					want,
				)
			}
		})
	}
}
