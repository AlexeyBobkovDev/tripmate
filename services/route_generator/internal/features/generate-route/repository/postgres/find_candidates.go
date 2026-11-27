package generate_route_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/repository/postgres/pool"
	generate_route_service "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/features/generate-route/service"
)

const queryGetPlaces = `
SELECT p.id,
	c.name,
	p.name,
	description,
	place_type,
	categories,
	interests,
	ST_X(location::geometry) AS longitude,
	ST_Y(location::geometry) AS latitude,
	address,
	website_url,
	phone,
	visit_duration_minutes,
	status
FROM route_generator.places p
JOIN route_generator.cities c ON p.city_id = c.id
WHERE c.name = $1 AND
	interests && $2::TEXT[];
`

const queryGetOpeningHours = `
SELECT day_of_week,
       opens_at,
       closes_at
FROM route_generator.place_opening_hours
WHERE place_id = $1
  AND is_closed = FALSE;
`

const queryGetOpeningHoursExceptions = `
SELECT e.exception_date,
       ei.opens_at,
       ei.closes_at
FROM route_generator.place_opening_hour_exceptions e
         JOIN route_generator.place_opening_hour_exception_intervals ei ON e.id = ei.exception_id
WHERE place_id = $1
  AND is_closed = FALSE
  AND e.exception_date = $2;
`

func (r *GenerateRouteRepository) FindCandidates(
	ctx context.Context,
	criteria generate_route_service.FindCandidatesCriteria,
) ([]domain.Place, error) {
	var places []domain.Place

	rows, err := r.pool.Query(
		ctx,
		queryGetPlaces,
		criteria.CityName,
		criteria.Interests,
	)
	if err != nil {
		return nil, fmt.Errorf("get places: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var routePlace domain.Place
		if err := rows.Scan(
			&routePlace.ID,
			&routePlace.CityName,
			&routePlace.Name,
			&routePlace.Description,
			&routePlace.PlaceType,
			&routePlace.Categories,
			&routePlace.Interests,
			&routePlace.Latitude,
			&routePlace.Longitude,
			&routePlace.Address,
			&routePlace.WebsiteURL,
			&routePlace.Phone,
			&routePlace.VisitDurationMinutes,
			&routePlace.Status,
		); err != nil {
			return nil, fmt.Errorf("scan places: %w", err)
		}
		places = append(places, routePlace)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	for ind, place := range places {
		openingHours, err := getOpeningHours(ctx, r.pool, place.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan openingHours")
		}
		places[ind].OpeningHours = openingHours

		openingHoursException, err := getOpeningHoursExceptions(ctx, r.pool, place.ID, criteria.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan openingHoursException")
		}
		places[ind].Exceptions = openingHoursException
	}

	return places, nil
}

func getOpeningHours(
	ctx context.Context,
	repository core_postgres_pool.Pool,
	placeID int,
) ([]domain.Schedule, error) {
	var openingHours []domain.Schedule
	rows, err := repository.Query(
		ctx,
		queryGetOpeningHours,
		placeID,
	)
	if err != nil {
		return []domain.Schedule{}, fmt.Errorf("fetch opening hours: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule domain.Schedule
		if err := rows.Scan(
			&schedule.WeekDayNumber,
			&schedule.OpensAt,
			&schedule.ClosesAt,
		); err != nil {
			return []domain.Schedule{}, fmt.Errorf("scan opening hours: %w", err)
		}
		openingHours = append(openingHours, schedule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan openingHours: %w", err)
	}

	return openingHours, nil
}

func getOpeningHoursExceptions(
	ctx context.Context,
	repository core_postgres_pool.Pool,
	placeID int,
	date time.Time,
) ([]domain.Exception, error) {
	var openingHoursException []domain.Exception
	rows, err := repository.Query(
		ctx,
		queryGetOpeningHoursExceptions,
		placeID,
		date,
	)
	if err != nil {
		return []domain.Exception{}, fmt.Errorf("fetch opening hours: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var exception domain.Exception
		if err := rows.Scan(
			&exception.ExceptionDate,
			&exception.OpensAt,
			&exception.ClosesAt,
		); err != nil {
			return []domain.Exception{}, fmt.Errorf("scan opening hours: %w", err)
		}
		openingHoursException = append(openingHoursException, exception)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan openingHoursException: %w", err)
	}

	return openingHoursException, nil
}
