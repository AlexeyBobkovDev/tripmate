CREATE TABLE IF NOT EXISTS app.route_stops
(
    id             BIGSERIAL PRIMARY KEY,
    route_day_id   BIGINT           NOT NULL REFERENCES app.route_days (id) ON DELETE CASCADE,
    place_id       BIGINT           NOT NULL,

    position       SMALLINT         NOT NULL,

    name           VARCHAR(100)     NOT NULL,
    description    TEXT             NOT NULL,
    categories             TEXT[]                     NOT NULL,
    interests              TEXT[]                     NOT NULL,
    address                TEXT                       NOT NULL,
    website_url            TEXT,
    phone                  VARCHAR(32),
    visit_duration_minutes SMALLINT                   NOT NULL,

    latitude       DOUBLE PRECISION NOT NULL,
    longitude      DOUBLE PRECISION NOT NULL,

    arrival_time   TIMESTAMPTZ      NOT NULL,
    departure_time TIMESTAMPTZ      NOT NULL,

    CONSTRAINT route_stops_place_id_check CHECK (place_id > 0),
    CONSTRAINT route_stops_position_check CHECK (position > 0),
    CONSTRAINT route_stops_name_check CHECK (CHAR_LENGTH(TRIM(name)) > 0),
    CONSTRAINT route_stops_description_check CHECK (CHAR_LENGTH(TRIM(description)) > 0),
    CONSTRAINT route_stops_categories CHECK (cardinality(categories) > 0),
    CONSTRAINT route_stops_interests CHECK (cardinality(interests) > 0),
    CONSTRAINT route_stops_address CHECK (LENGTH(TRIM(address)) > 0),
    CONSTRAINT route_stops_website_url CHECK (
        website_url IS NULL
            OR website_url ~ '^https?://.*$'
        ),
    CONSTRAINT route_stops_phone CHECK (
        phone IS NULL
            OR phone ~ '^\+[0-9]{10,15}$'
        ),
    CONSTRAINT route_stops_visit_duration_minutes CHECK (visit_duration_minutes > 0),

    CONSTRAINT route_stops_check_latitude CHECK(
        latitude > -90 OR latitude < 90
    ),
    CONSTRAINT route_stops_check_longitude CHECK(
        longitude > -180 OR longitude < 180
    )

    -- CONSTRAINT route_stops_departure_time_check CHECK (departure_time >= (arrival_time + visit_duration_minutes * INTERVAL '1 minute'))
);
