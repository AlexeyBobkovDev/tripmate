CREATE TYPE route_generator.place_type AS ENUM (
    'exterior',
    'interior',
    'guided_tour',
    'event'
    );

CREATE TABLE IF NOT EXISTS route_generator.places
(
    id                     BIGSERIAL PRIMARY KEY,
    name                   VARCHAR(255)               NOT NULL,
    description            TEXT                       NOT NULL,
    short_description      VARCHAR(255)               NOT NULL,
    place_type             route_generator.place_type NOT NULL,
    categories             TEXT[]                     NOT NULL,
    interests              TEXT[]                     NOT NULL,
    location               GEOGRAPHY(POINT, 4326)     NOT NULL UNIQUE,
    address                TEXT                       NOT NULL UNIQUE,
    website_url            TEXT UNIQUE,
    phone                  VARCHAR(32) UNIQUE,
    visit_duration_minutes SMALLINT                   NOT NULL,
    price_status           VARCHAR(20)                NOT NULL,
    price_min_minor_units  BIGINT                     NOT NULL,
    price_max_minor_units  BIGINT                     NOT NULL,
    price_currency         VARCHAR(3)                 NOT NULL,
    status                 VARCHAR(30),
    created_at             TIMESTAMPTZ DEFAULT NOW(),
    updated_at             TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT places_id CHECK (id > 0),
    CONSTRAINT places_name CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT places_description CHECK (LENGTH(TRIM(description)) > 0),
    CONSTRAINT places_short_description CHECK (LENGTH(TRIM(short_description)) > 0),
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
    CONSTRAINT places_price_status CHECK (LENGTH(TRIM(price_status)) > 0),
    CONSTRAINT places_price_min_minor_units CHECK (price_min_minor_units >= 0),
    CONSTRAINT places_price_max_minor_units CHECK (price_max_minor_units >= price_min_minor_units),
    CONSTRAINT places_price_currency CHECK (LENGTH(TRIM(price_currency)) > 0),
    CONSTRAINT places_status CHECK (LENGTH(TRIM(status)) > 0)
);
