CREATE TABLE IF NOT EXISTS app.route_days
(
    id         BIGSERIAL PRIMARY KEY,
    route_id   BIGINT   NOT NULL REFERENCES app.routes (id) ON DELETE CASCADE,

    date       DATE     NOT NULL,
    day_number BIGINT   NOT NULL,

    CONSTRAINT route_day_route_id_check CHECK (route_id > 0),
    CONSTRAINT route_day_day_number_check CHECK (day_number > 0)
);
