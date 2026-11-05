CREATE TYPE app.route_status AS ENUM ('generated', 'edited');

CREATE TABLE IF NOT EXISTS app.routes
(
    id           BIGSERIAL PRIMARY KEY,
    trip_id      BIGINT           NOT NULL REFERENCES app.trips (id) ON DELETE CASCADE,
    trip_version BIGINT           NOT NULL,
    version      BIGINT           NOT NULL DEFAULT 1,
    status       app.route_status NOT NULL,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    CONSTRAINT routes_trip_version_check CHECK (trip_version > 0),
    CONSTRAINT routes_version_check CHECK (version > 0),
    CONSTRAINT routes_updated_at_check CHECK (updated_at >= created_at)
);
