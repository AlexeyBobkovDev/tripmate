CREATE TABLE IF NOT EXISTS route_generator.place_opening_hour_exceptions
(
    id             BIGSERIAL PRIMARY KEY,
    place_id       BIGINT      NOT NULL REFERENCES route_generator.places (id) ON DELETE CASCADE,
    exception_date DATE        NOT NULL,
    is_closed      BOOLEAN     NOT NULL,
    reason         TEXT        NOT NULL,
    source_url     TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT place_opening_hour_exceptions_unique UNIQUE (place_id, exception_date),
    CONSTRAINT place_opening_hour_exceptions_id CHECK (id > 0),
    CONSTRAINT place_opening_hour_exceptions_place_id CHECK (place_id > 0),
    CONSTRAINT place_opening_hour_exceptions_reason CHECK (LENGTH(TRIM(reason)) > 0),
    CONSTRAINT place_opening_hour_exceptions_source_url CHECK (
        source_url IS NULL
            OR source_url ~ '^https?://.*$'
        )
);
