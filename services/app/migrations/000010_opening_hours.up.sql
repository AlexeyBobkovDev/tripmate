CREATE TABLE IF NOT EXISTS app.opening_hours(
    id BIGSERIAL PRIMARY KEY,
    route_stop_id BIGINT NOT NULL REFERENCES app.route_stops(id) ON DELETE CASCADE,
    opens_at SMALLINT,
    closes_at SMALLINT,

    CONSTRAINT place_time_correspondence CHECK (
        (
                opens_at IS NULL
                AND closes_at IS NULL
            )
            OR
        (
            opens_at IS NOT NULL
                AND closes_at IS NOT NULL
                AND opens_at BETWEEN 0 AND 10080
                AND closes_at > opens_at
                AND closes_at <= opens_at + 1440
            )
        )
)
