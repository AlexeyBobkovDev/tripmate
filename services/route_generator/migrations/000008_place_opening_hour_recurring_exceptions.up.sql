CREATE TYPE route_generator.exception_recurrence_type AS ENUM (
    'annual_date',
    'nth_weekday_of_month'
    );

CREATE TABLE IF NOT EXISTS route_generator.place_opening_hour_recurring_exceptions
(
    id              BIGSERIAL PRIMARY KEY,
    place_id        BIGINT                                    NOT NULL REFERENCES route_generator.places (id) ON DELETE CASCADE,
    recurrence_type route_generator.exception_recurrence_type NOT NULL,
    month           SMALLINT,
    day_of_month    SMALLINT,
    day_of_week     SMALLINT,
    week_of_month   SMALLINT,
    valid_from      DATE                                      NOT NULL,
    valid_until     DATE,
    is_closed       BOOLEAN                                   NOT NULL,
    reason          TEXT,
    source_url      TEXT,
    created_at      TIMESTAMPTZ                               NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ                               NOT NULL DEFAULT NOW(),

    CONSTRAINT place_recurring_exceptions_id CHECK (id > 0),
    CONSTRAINT place_recurring_exceptions_place_id CHECK (place_id > 0),
    CONSTRAINT place_recurring_exceptions_month CHECK (
        month IS NULL
            OR month BETWEEN 1 AND 12
        ),
    CONSTRAINT place_recurring_exceptions_day_of_month CHECK (
        day_of_month IS NULL
            OR
        (
            month IN (1, 3, 5, 7, 8, 10, 12) AND day_of_month BETWEEN 1 AND 31
                OR
            month = 2 AND day_of_month BETWEEN 1 AND 29
                OR
            month IN (4, 6, 9, 11) AND day_of_month BETWEEN 1 AND 30
            )
        ),
    CONSTRAINT place_recurring_exceptions_day_of_week CHECK (
        day_of_week IS NULL
            OR day_of_week BETWEEN 1 AND 7
        ),
    CONSTRAINT place_recurring_exceptions_week_of_month CHECK (
        week_of_month IS NULL
            OR week_of_month BETWEEN 1 AND 5
            OR week_of_month = -1
        ),
    CONSTRAINT place_recurring_exceptions_valid_from_until CHECK (
        valid_until IS NULL
            OR valid_from < valid_until
        ),
    CONSTRAINT place_recurring_exceptions_reason CHECK (reason IS NULL OR LENGTH(TRIM(reason)) > 0),
    CONSTRAINT place_recurring_exceptions_source_url CHECK (
        source_url IS NULL
            OR source_url ~ '^https?://.*$'
        ),

    CONSTRAINT place_recurring_recurrence_type CHECK (
        (
            recurrence_type = 'annual_date'
                AND
            month IS NOT NULL
                AND
            day_of_month IS NOT NULL
                AND
            day_of_week IS NULL
                AND
            week_of_month IS NULL
            )
            OR
        (
            recurrence_type = 'nth_weekday_of_month'
                AND
            month IS NULL
                AND
            day_of_month IS NULL
                AND
            day_of_week IS NOT NULL
                AND
            week_of_month IS NOT NULL
            )
        )
);
