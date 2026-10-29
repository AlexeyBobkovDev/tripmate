package testutils

import (
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func TripTest() domain.Trip {
	startDate, _ := time.Parse("2006-01-02", "2026-11-14")
	endDate, _ := time.Parse("2006-01-02", "2026-11-16")
	now := time.Now()

	return domain.Trip{
		ID:           1,
		OwnerID:      1,
		Version:      1,
		Name:         "Walk around Moscow",
		PlanningArea: "Moscow",
		StartLocation: domain.Location{
			Name:      "Red Square",
			Latitude:  55.7539,
			Longitude: 37.6208,
		},
		StartDate: startDate,
		EndDate:   endDate,
		Budget: domain.Money{
			AmountMinorUnits: 1500000,
			Currency:         "RUB",
		},
		Interests: []domain.Interest{"history", "nature"},
		TimeZone:  "Europe/Moscow",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TripPatchTest() domain.TripPatch {
	name := "Updated Trip Name"

	tp := domain.NewTripPatch(
		1,
		1,
		domain.Nullable[string]{Value: &name, Set: true},
		domain.Nullable[string]{Set: false},
		domain.Nullable[domain.Location]{Set: false},
		domain.Nullable[time.Time]{Set: false},
		domain.Nullable[time.Time]{Set: false},
		domain.Nullable[domain.Money]{Set: false},
		domain.Nullable[[]domain.Interest]{Set: false},
		domain.Nullable[string]{Set: false},
	)
	return *tp
}

func TripCreateTest(ownerID int, name string) *domain.TripCreate {
	startDate := time.Date(time.Now().Year()+1, time.March, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(time.Now().Year()+1, time.March, 12, 0, 0, 0, 0, time.UTC)

	trip := domain.NewTripCreate(
		name,
		"Moscow",
		domain.Location{
			Name:      "Red Square",
			Latitude:  55.7539,
			Longitude: 37.6208,
		},
		startDate,
		endDate,
		domain.Money{
			AmountMinorUnits: 1500000,
			Currency:         "RUB",
		},
		[]domain.Interest{"history", "nature"},
		"Europe/Moscow",
	)
	trip.OwnerID = ownerID

	return trip
}
