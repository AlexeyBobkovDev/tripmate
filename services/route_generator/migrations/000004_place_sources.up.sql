CREATE TABLE IF NOT EXISTS route_generator.place_sources
(
    id                BIGSERIAL PRIMARY KEY,
    place_id          BIGINT      NOT NULL REFERENCES route_generator.places (id) ON DELETE CASCADE,
    provider          VARCHAR(50) NOT NULL,
    external_type     VARCHAR(20) NOT NULL,
    external_id       BIGINT      NOT NULL,
    source_url        TEXT        NOT NULL,
    source_updated_at TIMESTAMPTZ NOT NULL,
    last_synced_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT place_sources_id CHECK (id > 0),
    CONSTRAINT place_sources_place_id CHECK (place_id > 0),
    CONSTRAINT place_sources_provider CHECK (LENGTH(TRIM(provider)) > 0),
    CONSTRAINT place_sources_external_type CHECK (LENGTH(TRIM(external_type)) > 0),
    CONSTRAINT place_sources_external_id CHECK (external_id > 0),
    CONSTRAINT place_sources_source_url CHECK (
        source_url IS NULL
            OR source_url ~ '^https?://.*$'
        ),
    CONSTRAINT route_generator_unique UNIQUE (provider, external_type, external_id)
);
