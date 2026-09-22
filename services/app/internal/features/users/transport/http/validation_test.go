package users_transport_http

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
			name:      "birthday today",
			birthDate: time.Date(2000, 9, 22, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC),
			want:      26,
		},
		{
			name:      "birthday tomorrow",
			birthDate: time.Date(2000, 9, 23, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC),
			want:      25,
		},
		{
			name:      "birthday yesterday",
			birthDate: time.Date(2000, 9, 21, 0, 0, 0, 0, time.UTC),
			now:       time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC),
			want:      26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAge(tt.birthDate, tt.now)

			if got != tt.want {
				t.Fatalf("calculateAge() = %d, want %d", got, tt.want)
			}
		})
	}
}

func FuzzCalculateAge(f *testing.F) {
	f.Add(
		time.Date(2000, 9, 22, 0, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC).Unix(),
	)

	f.Add(
		time.Date(2000, 9, 23, 0, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC).Unix(),
	)

	f.Add(
		time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC).Unix(),
		time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC).Unix(),
	)

	f.Fuzz(func(t *testing.T, birthdateUnix int64, nowUnix int64) {
		birthDate := time.Unix(birthdateUnix, 0).UTC()
		now := time.Unix(nowUnix, 0).UTC()

		if birthDate.After(now) {
			t.Skip()
		}

		age := calculateAge(birthDate, now)

		birthDateThisYear := time.Date(
			now.Year(),
			birthDate.Month(),
			birthDate.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		correctAge := now.Year() - birthDate.Year()
		if birthDateThisYear.After(now) {
			correctAge--
		}

		if age < 0 {
			t.Fatal("age is negative")
		}

		if age != correctAge {
			t.Fatalf(
				"calculateAge(%v, %v) = %d, want %d",
				birthDate,
				now,
				age,
				correctAge,
			)
		}
	})
}
