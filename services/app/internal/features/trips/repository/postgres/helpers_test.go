package trips_postgres_repository

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const tripScanDestinationsAmount = 18

// anyArgs returns `amount` gomock.Any() matchers. It is needed for variadic
// methods such as Rows.Scan, where gomock matches each variadic element.
func anyArgs(amount int) []any {
	matchers := make([]any, amount)
	for index := range matchers {
		matchers[index] = gomock.Any()
	}

	return matchers
}

// testTripFixture builds a fully populated domain.Trip without meal
// requirements. Meal requirements are loaded by dedicated queries, so they are
// not part of the trip row scan.
func testTripFixture() domain.Trip {
	startDate := time.Date(2027, time.March, 10, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.September, 9, 17, 0, 45, 0, time.UTC)
	updatedAt := time.Date(2026, time.September, 9, 17, 10, 45, 0, time.UTC)

	return domain.Trip{
		ID:               1,
		OwnerID:          1,
		Version:          1,
		Name:             "Trip to Moscow",
		PlanningArea:     "Moscow",
		StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        startDate,
		EndDate:          startDate,
		DailySchedule: domain.DailySchedule{
			StartTime: "09:00",
			EndTime:   "18:00",
		},
		Budget:    domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests: []domain.Interest{"history"},
		TimeZone:  "Europe/Moscow",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func testTripCreateFixture(ownerID int) *domain.TripCreate {
	startDate := time.Date(2027, time.March, 10, 0, 0, 0, 0, time.UTC)

	return &domain.TripCreate{
		OwnerID:          ownerID,
		Name:             "Trip to Moscow",
		PlanningArea:     "Moscow",
		StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        startDate,
		EndDate:          startDate,
		DailySchedule: domain.DailySchedule{
			StartTime: "09:00",
			EndTime:   "18:00",
		},
		Budget:    domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests: []domain.Interest{"history"},
		TimeZone:  "Europe/Moscow",
	}
}

// fillTripScan returns a gomock DoAndReturn callback that fills the 18 columns
// selected by the trip queries.
func fillTripScan(t *testing.T, trip domain.Trip) func(dest ...any) error {
	t.Helper()

	startTime, err := time.Parse("15:04", string(trip.DailySchedule.StartTime))
	if err != nil {
		t.Fatalf("parse fixture start time: %v", err)
	}
	endTime, err := time.Parse("15:04", string(trip.DailySchedule.EndTime))
	if err != nil {
		t.Fatalf("parse fixture end time: %v", err)
	}

	return func(dest ...any) error {
		if len(dest) != tripScanDestinationsAmount {
			return fmt.Errorf(
				"expected %d scan destinations, got %d",
				tripScanDestinationsAmount,
				len(dest),
			)
		}

		*dest[0].(*int) = trip.ID
		*dest[1].(*int) = trip.OwnerID
		*dest[2].(*int) = trip.Version
		*dest[3].(*string) = trip.Name
		*dest[4].(*string) = trip.PlanningArea
		*dest[5].(*float64) = trip.StartCoordinates.Latitude
		*dest[6].(*float64) = trip.StartCoordinates.Longitude
		*dest[7].(*time.Time) = trip.StartDate
		*dest[8].(*time.Time) = trip.EndDate
		*dest[9].(*time.Time) = startTime
		*dest[10].(*time.Time) = endTime
		*dest[11].(*int64) = trip.Budget.AmountMinorUnits
		*dest[12].(*string) = trip.Budget.Currency
		*dest[13].(*[]domain.Interest) = trip.Interests
		*dest[14].(*string) = trip.TimeZone
		*dest[15].(*time.Time) = trip.CreatedAt
		*dest[16].(*time.Time) = trip.UpdatedAt
		*dest[17].(**time.Time) = trip.DeletedAt

		return nil
	}
}
