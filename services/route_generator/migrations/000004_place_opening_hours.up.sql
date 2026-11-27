CREATE TABLE IF NOT EXISTS route_generator.place_opening_hours
(
    id          BIGSERIAL PRIMARY KEY,
    place_id    BIGINT   NOT NULL REFERENCES route_generator.places (id) ON DELETE CASCADE,
    day_of_week SMALLINT NOT NULL,
    opens_at    SMALLINT,
    closes_at   SMALLINT,
    is_closed   BOOLEAN  NOT NULL DEFAULT FALSE,

    CONSTRAINT place_opening_hours_id CHECK (id > 0),
    CONSTRAINT place_opening_hours_place_id CHECK (place_id > 0),
    CONSTRAINT place_opening_hours_day_of_week CHECK (day_of_week BETWEEN 1 AND 7),
    CONSTRAINT place_time_correspondence CHECK (
        (
            is_closed = true
                AND opens_at IS NULL
                AND closes_at IS NULL
            )
            OR
        (
            is_closed = false
                AND opens_at IS NOT NULL
                AND closes_at IS NOT NULL
                AND opens_at BETWEEN (day_of_week - 1) * 1440
                AND day_of_week * 1440 - 1
                AND closes_at > opens_at
                AND closes_at <= opens_at + 1440
            )
        )
);
