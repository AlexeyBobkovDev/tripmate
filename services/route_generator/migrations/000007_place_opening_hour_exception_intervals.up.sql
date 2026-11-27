CREATE TABLE IF NOT EXISTS route_generator.place_opening_hour_exception_intervals
(
    id           BIGSERIAL PRIMARY KEY,
    exception_id BIGINT NOT NULL REFERENCES route_generator.place_opening_hour_exceptions (id) ON DELETE CASCADE,
    opens_at     SMALLINT,
    closes_at    SMALLINT,

    CONSTRAINT place_opening_hour_exception_id CHECK (id > 0),
    CONSTRAINT place_opening_hour_exception_exception_id CHECK (exception_id > 0),
    CONSTRAINT place_opening_hour_exception_intervals CHECK (
        (
            opens_at IS NULL
                AND closes_at IS NULL
            )
            OR
        (
            opens_at IS NOT NULL
                AND closes_at IS NOT NULL
                AND opens_at BETWEEN 0 AND 1439
                AND closes_at > opens_at
                AND closes_at <= opens_at + 1440
            )
        )
);
