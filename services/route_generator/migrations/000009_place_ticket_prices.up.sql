CREATE TYPE route_generator.visitor_category AS ENUM(
    'adult',
    'child',
    'student',
    'senior'
);

CREATE TYPE route_generator.price_status AS ENUM('paid', 'free', 'unknown');

CREATE TABLE IF NOT EXISTS route_generator.place_ticket_prices(
    id BIGSERIAL PRIMARY KEY,
    place_id BIGINT NOT NULL REFERENCES route_generator.places(id) ON DELETE CASCADE,
    visitor_category route_generator.visitor_category NOT NULL,
    price_status route_generator.price_status NOT NULL,
    price_minor_units BIGINT,
    currency VARCHAR(3),
    valid_from DATE NOT NULL DEFAULT CURRENT_DATE,
    valid_to DATE,
    source_url TEXT,
    verified_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT place_ticket_rices_id CHECK(id > 0),
    CONSTRAINT place_ticket_rices_place_id CHECK(place_id > 0),
    CONSTRAINT place_ticket_rices_price_minor_units CHECK(
        price_minor_units IS NULL
    OR price_minor_units >= 0
    ),


    CONSTRAINT place_ticket_prices_price CHECK (
        (
            price_status = 'paid'
            AND price_minor_units > 0
            AND currency ~ '^[A-Z]{3}$'
        )
        OR
        (
            price_status = 'free'
            AND price_minor_units = 0
            AND currency ~ '^[A-Z]{3}$'
        )
        OR
        (
            price_status = 'unknown'
            AND price_minor_units IS NULL
            AND currency IS NULL
        )
    ),

    CONSTRAINT place_ticket_prices_validity CHECK (
        valid_to IS NULL OR valid_from < valid_to
    ),

    CONSTRAINT place_ticket_prices_source_url CHECK (
        source_url IS NULL
        OR source_url ~ '^https?://.*$'
    )
);
