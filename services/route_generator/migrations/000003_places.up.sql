CREATE TYPE route_generator.place_type AS ENUM (
    'exterior',
    'interior',
    'guided_tour',
    'event'
    );

CREATE TABLE IF NOT EXISTS route_generator.places
(
    id                     BIGSERIAL PRIMARY KEY,
    city_id                BIGINT                     NOT NULL REFERENCES route_generator.cities (id),
    name                   VARCHAR(255)               NOT NULL,
    description            TEXT                       NOT NULL,
    place_type             route_generator.place_type NOT NULL,
    categories             TEXT[]                     NOT NULL,
    interests              TEXT[]                     NOT NULL,
    location               GEOGRAPHY(POINT, 4326)     NOT NULL,
    address                TEXT                       NOT NULL,
    website_url            TEXT UNIQUE,
    phone                  VARCHAR(32),
    visit_duration_minutes SMALLINT                   NOT NULL,
    status                 VARCHAR(30)                NOT NULL,
    created_at             TIMESTAMPTZ                NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ                NOT NULL DEFAULT NOW(),

    CONSTRAINT places_id CHECK (id > 0),
    CONSTRAINT places_name CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT places_description CHECK (LENGTH(TRIM(description)) > 0),
    -- CONSTRAINT route_generator_places_location CHECK(location),
    CONSTRAINT places_address CHECK (LENGTH(TRIM(address)) > 0),
    CONSTRAINT places_categories CHECK (cardinality(categories) > 0),
    CONSTRAINT places_interests CHECK (cardinality(interests) > 0),
    CONSTRAINT places_website_url CHECK (
        website_url IS NULL
            OR website_url ~ '^https?://.*$'
        ),
    CONSTRAINT places_phone CHECK (
        phone IS NULL
            OR phone ~ '^\+[0-9]{10,15}$'
        ),
    CONSTRAINT places_visit_duration_minutes CHECK (visit_duration_minutes > 0),
    CONSTRAINT place_check_location_is_not_empty CHECK ( NOT ST_IsEmpty(location::geometry)),
    CONSTRAINT places_status CHECK (status IN ('active', 'inactive', 'temporarily_closed'))
);
