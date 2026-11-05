package routes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const getRoutesQuery = `
SELECT id,
	trip_id,
	trip_version,
	version,
	status,
	created_at,
	updated_at
FROM app.routes
WHERE trip_id=$1;
`

const getRouteDaysQuery = `
SELECT id,
	date,
	day_number
FROM app.route_days
WHERE route_id=$1;
`

const getRouteStopsQuery = `
SELECT id,
	place_id,
	position,
	name,
	description,
	latitude,
	longitude,
	arrival_time,
	departure_time,
	visit_duration
FROM app.route_stops
WHERE route_day_id=$1;
`

const getRouteLegsQuery = `
SELECT id,
	from_stop_id,
	to_stop_id,
	travel_mode,
	distance,
	duration
FROM app.route_legs
WHERE route_day_id=$1;
`

func (r *RoutesRepository) GetRoutes(
	ctx context.Context,
	tripID int,
) ([]*domain.Route, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	defer func() { _ = tx.Commit(ctx) }()

	var routeResponse []*domain.Route

	routeCounter := 0
	routeRows, err := tx.Query(
		ctx,
		getRoutesQuery,
		tripID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route rows: %w", err)
	}
	for routeRows.Next() {
		routeResponse = append(routeResponse, &domain.Route{})
		if err := routeRows.Scan(
			&routeResponse[routeCounter].ID,
			&routeResponse[routeCounter].TripID,
			&routeResponse[routeCounter].TripVersion,
			&routeResponse[routeCounter].Version,
			&routeResponse[routeCounter].Status,
			&routeResponse[routeCounter].CreatedAt,
			&routeResponse[routeCounter].UpdatedAt,
		); err != nil {
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to fetch route: %w", err)
		}
		routeCounter++

	}

	for routeRowsInd := range routeResponse {
		routeDayCounter := 0
		routeDayRows, err := tx.Query(
			ctx,
			getRouteDaysQuery,
			routeResponse[routeRowsInd].ID,
		)
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to fetch route day rows: %w", err)
		}
		for routeDayRows.Next() {
			routeResponse[routeRowsInd].Days = append(routeResponse[routeRowsInd].Days, domain.RouteDay{})
			if err := routeDayRows.Scan(
				&routeResponse[routeRowsInd].Days[routeDayCounter].ID,
				&routeResponse[routeRowsInd].Days[routeDayCounter].Date,
				&routeResponse[routeRowsInd].Days[routeDayCounter].DayNumber,
			); err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to fetch route day: %w", err)
			}
			routeDayCounter++
		}
	}

	for routeRowsInd := range routeResponse {
		for routeDaysRowsInd := range routeResponse[routeRowsInd].Days {
			routeStopCounter := 0
			routeStopRows, err := tx.Query(
				ctx,
				getRouteStopsQuery,
				routeResponse[routeRowsInd].Days[routeDaysRowsInd].ID,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to fetch route stop rows: %w", err)
			}
			for routeStopRows.Next() {
				routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops = append(routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops, domain.RouteStop{})

				var visitDurationSeconds int64
				if err := routeStopRows.Scan(
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].ID,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].PlaceID,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].Position,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].Name,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].Description,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].Coordinates.Latitude,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].Coordinates.Longitude,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].ArrivalTime,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].DepartureTime,
					&visitDurationSeconds,
				); err != nil {
					_ = tx.Rollback(ctx)
					return nil, fmt.Errorf("failed to fetch route stop: %w", err)
				}
				routeResponse[routeRowsInd].Days[routeDaysRowsInd].Stops[routeStopCounter].VisitDuration = durationFromSeconds(visitDurationSeconds)
				routeStopCounter++
			}

			routeLegCounter := 0
			routeLegRows, err := tx.Query(
				ctx,
				getRouteLegsQuery,
				routeResponse[routeRowsInd].Days[routeDaysRowsInd].ID,
			)
			if err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to fetch route leg rows: %w", err)
			}
			for routeLegRows.Next() {
				routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs = append(routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs, domain.RouteLeg{})

				var durationSeconds int64
				if err := routeLegRows.Scan(
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs[routeLegCounter].ID,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs[routeLegCounter].FromStopID,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs[routeLegCounter].ToStopID,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs[routeLegCounter].TravelMode,
					&routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs[routeLegCounter].Distance,
					&durationSeconds,
				); err != nil {
					_ = tx.Rollback(ctx)
					return nil, fmt.Errorf("failed to create route leg: %w", err)
				}
				routeResponse[routeRowsInd].Days[routeDaysRowsInd].Legs[routeLegCounter].Duration = durationFromSeconds(durationSeconds)
				routeLegCounter++
			}
		}
	}
	return routeResponse, nil
}
