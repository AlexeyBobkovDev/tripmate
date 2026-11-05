CREATE TYPE app.travel_mode AS ENUM (
    'walking',
    'cycling',
    'driving',
    'transit'
    );

CREATE TABLE IF NOT EXISTS app.route_legs
(
    id             BIGSERIAL PRIMARY KEY,
    route_day_id   BIGINT    REFERENCES app.route_days (id) ON DELETE CASCADE,

    from_stop_id BIGINT          NOT NULL REFERENCES app.route_stops (id),
    to_stop_id   BIGINT          NOT NULL REFERENCES app.route_stops (id),

    travel_mode  app.travel_mode NOT NULL,
    distance     BIGINT          NOT NULL,
    duration     BIGINT          NOT NULL,

    CONSTRAINT route_legs_from_stop_id_check CHECK (from_stop_id > 0),
    CONSTRAINT route_legs_to_stop_id_check CHECK (to_stop_id > 0),
    CONSTRAINT route_legs_distance_check CHECK (distance > 0),
    CONSTRAINT route_legs_duration_check CHECK (duration > 0)
);
