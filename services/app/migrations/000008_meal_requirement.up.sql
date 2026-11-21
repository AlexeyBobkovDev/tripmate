CREATE TYPE app.meal_type AS ENUM (
    'breakfast',
    'lunch',
    'dinner'
    );

CREATE TABLE IF NOT EXISTS app.meal_requirements
(
    id                     SERIAL PRIMARY KEY,
    trip_id                BIGINT        NOT NULL REFERENCES app.trips (id) ON DELETE CASCADE,
    meal_type              app.meal_type NOT NULL,
    start_time             TIMESTAMPTZ   NOT NULL,
    end_time               TIMESTAMPTZ   NOT NULL,
    visit_duration_minutes INT           NOT NULL,
    amount_minor_units     INT           NOT NULL,
    currency               VARCHAR(3),

    CONSTRAINT meal_requirements_id CHECK (id > 0),
    CONSTRAINT meal_requirements_visit_duration_minutes CHECK (visit_duration_minutes > 0),
    CONSTRAINT meal_requirements_amount_minor_units CHECK (amount_minor_units >= 0),
    CONSTRAINT meal_requirements_currency CHECK (LENGTH(TRIM(currency)) > 0)
);
