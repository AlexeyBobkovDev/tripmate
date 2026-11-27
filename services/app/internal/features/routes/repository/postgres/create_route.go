package routes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const createRouteQuery = `
INSERT INTO app.routes (
	trip_id,
	trip_version,
	status
)
VALUES($1, $2, $3)
RETURNING
	id,
	trip_id,
	trip_version,
	version,
	status,
	created_at,
	updated_at;
`

const createRouteDayQuery = `
INSERT INTO app.route_days (
	route_id,
	date,
	day_number
)
VALUES($1, $2, $3)
RETURNING
	id,
	date,
	day_number;
`

const createRouteStopQuery = `
INSERT INTO app.route_stops (
	route_day_id,
	place_id,
	position,
	name,
	description,
	categories,
	interests,
	address,
	website_url,
	phone,
	visit_duration_minutes,
	latitude,
	longitude,
	arrival_time,
	departure_time
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING
	id,
	place_id,
	position,
	name,
	description,
	categories,
	interests,
	address,
	website_url,
	phone,
	visit_duration_minutes,
	latitude,
	longitude,
	arrival_time,
	departure_time;
`

const createRouteLegQuery = `
INSERT INTO app.route_legs (
	route_day_id,
	from_stop_id,
	to_stop_id,
	travel_mode,
	distance,
	duration
)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING
	id,
	from_stop_id,
	to_stop_id,
	travel_mode,
	distance,
	duration;
`

const createOpeningHoursQuery = `
INSERT INTO app.opening_hours (
	route_stop_id,
	opens_at,
	closes_at
)
VALUES($1, $2, $3)
RETURNING opens_at,
	closes_at;
`

func (r *RoutesRepository) CreateRoute(
	ctx context.Context,
	ownerID int,
	trip *domain.Trip,
	generatedRoute *domain.RouteCreate,
) (*domain.Route, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	defer func() { _ = tx.Commit(ctx) }()

	var routeResponse domain.Route

	if err := tx.QueryRow(
		ctx,
		createRouteQuery,
		trip.ID,
		trip.Version,
		generatedRoute.Status,
	).Scan(
		&routeResponse.ID,
		&routeResponse.TripID,
		&routeResponse.TripVersion,
		&routeResponse.Version,
		&routeResponse.Status,
		&routeResponse.CreatedAt,
		&routeResponse.UpdatedAt,
	); err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	for day_ind, day := range generatedRoute.Days {
		routeResponse.Days = append(routeResponse.Days, domain.RouteDay{})
		if err := tx.QueryRow(
			ctx,
			createRouteDayQuery,
			routeResponse.ID,
			day.Date,
			day.DayNumber,
		).Scan(
			&routeResponse.Days[day_ind].ID,
			&routeResponse.Days[day_ind].Date,
			&routeResponse.Days[day_ind].DayNumber,
		); err != nil {
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to create route day: %w", err)
		}

		mapStopIDInTripToStopIDInDB := map[int]int{}

		for stop_ind, stop := range day.Stops {
			routeResponse.Days[day_ind].Stops = append(routeResponse.Days[day_ind].Stops, domain.RouteStop{})

			if err := tx.QueryRow(
				ctx,
				createRouteStopQuery,
				routeResponse.Days[day_ind].ID,
				stop.PlaceID,
				stop.Position,
				stop.Name,
				stop.Description,
				stop.Categories,
				stop.Interests,
				stop.Address,
				stop.WebsiteURL,
				stop.Phone,
				stop.VisitDurationMinutes,
				stop.Coordinates.Latitude,
				stop.Coordinates.Longitude,
				stop.ArrivalTime,
				stop.DepartureTime,
			).Scan(
				&routeResponse.Days[day_ind].Stops[stop_ind].ID,
				&routeResponse.Days[day_ind].Stops[stop_ind].PlaceID,
				&routeResponse.Days[day_ind].Stops[stop_ind].Position,
				&routeResponse.Days[day_ind].Stops[stop_ind].Name,
				&routeResponse.Days[day_ind].Stops[stop_ind].Description,
				&routeResponse.Days[day_ind].Stops[stop_ind].Categories,
				&routeResponse.Days[day_ind].Stops[stop_ind].Interests,
				&routeResponse.Days[day_ind].Stops[stop_ind].Address,
				&routeResponse.Days[day_ind].Stops[stop_ind].WebsiteURL,
				&routeResponse.Days[day_ind].Stops[stop_ind].Phone,
				&routeResponse.Days[day_ind].Stops[stop_ind].VisitDurationMinutes,
				&routeResponse.Days[day_ind].Stops[stop_ind].Coordinates.Latitude,
				&routeResponse.Days[day_ind].Stops[stop_ind].Coordinates.Longitude,
				&routeResponse.Days[day_ind].Stops[stop_ind].ArrivalTime,
				&routeResponse.Days[day_ind].Stops[stop_ind].DepartureTime,
			); err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to create route stop: %w", err)
			}
			mapStopIDInTripToStopIDInDB[generatedRoute.Days[day_ind].Stops[stop_ind].Position] = routeResponse.Days[day_ind].Stops[stop_ind].ID

			for _, openingHour := range stop.OpeningHours {
				var scannedOpeningHour domain.OpeningHours

				if err := tx.QueryRow(
					ctx,
					createOpeningHoursQuery,
					routeResponse.Days[day_ind].Stops[stop_ind].ID,
					openingHour.OpensAt,
					openingHour.ClosesAt,
				).Scan(
					&scannedOpeningHour.OpensAt,
					&scannedOpeningHour.ClosesAt,
				); err != nil {
					_ = tx.Rollback(ctx)
					return nil, fmt.Errorf("failed to insert opening hours: %w", err)
				}
				routeResponse.Days[day_ind].Stops[stop_ind].OpeningHours = append(routeResponse.Days[day_ind].Stops[stop_ind].OpeningHours, openingHour)
			}
		}

		for leg_ind, leg := range day.Legs {
			routeResponse.Days[day_ind].Legs = append(routeResponse.Days[day_ind].Legs, domain.RouteLeg{})

			var durationSeconds int64
			if err := tx.QueryRow(
				ctx,
				createRouteLegQuery,
				routeResponse.Days[day_ind].ID,
				mapStopIDInTripToStopIDInDB[leg.FromStopID],
				mapStopIDInTripToStopIDInDB[leg.ToStopID],
				leg.TravelMode,
				leg.Distance,
				durationToSeconds(leg.Duration),
			).Scan(
				&routeResponse.Days[day_ind].Legs[leg_ind].ID,
				&routeResponse.Days[day_ind].Legs[leg_ind].FromStopID,
				&routeResponse.Days[day_ind].Legs[leg_ind].ToStopID,
				&routeResponse.Days[day_ind].Legs[leg_ind].TravelMode,
				&routeResponse.Days[day_ind].Legs[leg_ind].Distance,
				&durationSeconds,
			); err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to create route leg: %w", err)
			}
			routeResponse.Days[day_ind].Legs[leg_ind].Duration = durationFromSeconds(durationSeconds)
		}
	}
	return &routeResponse, nil
}
