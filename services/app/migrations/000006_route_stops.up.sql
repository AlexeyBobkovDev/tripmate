CREATE TABLE IF NOT EXISTS app.route_stops
(
    id             BIGSERIAL PRIMARY KEY,
    route_day_id   BIGINT REFERENCES app.route_days (id) ON DELETE CASCADE,
    place_id       BIGINT           NOT NULL,

    position       SMALLINT         NOT NULL,
    name           VARCHAR(100)     NOT NULL,
    description    TEXT             NOT NULL,

    latitude       DOUBLE PRECISION NOT NULL,
    longitude      DOUBLE PRECISION NOT NULL,

    arrival_time   TIMESTAMPTZ      NOT NULL,
    departure_time TIMESTAMPTZ      NOT NULL,
    visit_duration BIGINT           NOT NULL,

    CONSTRAINT route_stops_place_id_check CHECK (place_id > 0),
    CONSTRAINT route_stops_position_check CHECK (position > 0),
    CONSTRAINT route_stops_name_check CHECK (CHAR_LENGTH(TRIM(name)) > 0),
    CONSTRAINT route_stops_description_check CHECK (CHAR_LENGTH(TRIM(description)) > 0),
    CONSTRAINT route_stops_departure_time_check CHECK (departure_time >= (arrival_time + visit_duration * INTERVAL '1 second')),
    CONSTRAINT route_stops_visit_duration_check CHECK (visit_duration > 0)
);
