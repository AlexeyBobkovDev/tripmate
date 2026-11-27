CREATE TABLE IF NOT EXISTS route_generator.cities
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,

    CONSTRAINT cities_id CHECK (id > 0),
    CONSTRAINT cities_name CHECK (LENGTH(TRIM(name)) > 0)
);
