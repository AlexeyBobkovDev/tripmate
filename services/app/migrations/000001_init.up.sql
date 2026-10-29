CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS app.users
(
    id            BIGSERIAL PRIMARY KEY,
    version       BIGINT        NOT NULL DEFAULT 1,
    name          VARCHAR(50)   NOT NULL,
    surname       VARCHAR(50)   NOT NULL,
    username      VARCHAR(100)  NOT NULL,
    birth_date    DATE          NOT NULL,
    description   VARCHAR(1000) NOT NULL,
    email         VARCHAR(100)  NOT NULL,
    phone_number  VARCHAR(16)   NOT NULL,
    password_hash TEXT          NOT NULL,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT users_version_check CHECK (version > 0),
    CONSTRAINT users_name_check CHECK (CHAR_LENGTH(TRIM(name)) > 0),
    CONSTRAINT users_surname_check CHECK (CHAR_LENGTH(TRIM(surname)) > 0),
    CONSTRAINT users_username_check CHECK (CHAR_LENGTH(TRIM(username)) > 0),
    CONSTRAINT users_birth_date_check CHECK (
        EXTRACT(
                YEAR
                FROM
                age(NOW(), birth_date)
        ) BETWEEN 14
            AND 150
        ),
    CONSTRAINT users_description_check CHECK (CHAR_LENGTH(TRIM(description)) > 0),
    CONSTRAINT users_email_check CHECK (email ~ '^[^@\s]+@[^@\s]+\.[^@\s]+$'),
    CONSTRAINT users_phone_number_check CHECK (
        CHAR_LENGTH(phone_number) BETWEEN 10
            AND 16
            AND phone_number ~ '^\+[0-9]{10,15}$'
        ),
    CONSTRAINT users_updated_at_check CHECK (updated_at >= created_at),
    CONSTRAINT users_deleted_at_check CHECK (
        deleted_at IS NULL
            OR deleted_at >= created_at
        )
);

CREATE UNIQUE INDEX idx_unique_user_username ON app.users (username)
    WHERE
        deleted_at IS NULL;

CREATE UNIQUE INDEX idx_unique_user_email ON app.users (email)
    WHERE
        deleted_at IS NULL;

CREATE UNIQUE INDEX idx_unique_user_phone_number ON app.users (phone_number)
    WHERE
        deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS app.trips
(
    id                 BIGSERIAL PRIMARY KEY,
    owner_id           BIGINT           NOT NULL REFERENCES app.users (id),
    version            BIGINT           NOT NULL DEFAULT 1,
    name               VARCHAR(100)     NOT NULL,
    planning_area      VARCHAR(100)     NOT NULL,
    location_name      VARCHAR(100)     NOT NULL,
    location_latitude  DOUBLE PRECISION NOT NULL,
    location_longitude DOUBLE PRECISION NOT NULL,
    start_date         DATE             NOT NULL,
    end_date           DATE             NOT NULL,
    amount_minor_units INT              NOT NULL,
    currency           VARCHAR(3)       NOT NULL,
    interests          TEXT[]           NOT NULL DEFAULT '{}',
    time_zone          VARCHAR(100)     NOT NULL,
    created_at         TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ,
    CONSTRAINT trips_owner_id_check CHECK (owner_id > 0),
    CONSTRAINT trips_version_check CHECK (version > 0),
    CONSTRAINT trips_name_check CHECK (CHAR_LENGTH(TRIM(name)) > 0),
    CONSTRAINT trips_planning_area_check CHECK (CHAR_LENGTH(TRIM(planning_area)) > 0),
    CONSTRAINT trips_location_name_check CHECK (CHAR_LENGTH(TRIM(location_name)) > 0),
    CONSTRAINT trips_location_latitude_check CHECK (
        location_latitude BETWEEN -90
            AND 90
        ),
    CONSTRAINT trips_location_longitude_check CHECK (
        location_longitude BETWEEN -180
            AND 180
        ),
    CONSTRAINT trips_check_trip_dates_order CHECK (end_date >= start_date),
    CONSTRAINT trips_amount_minor_units_check CHECK (amount_minor_units > 0),
    CONSTRAINT trips_currency_check CHECK (CHAR_LENGTH(TRIM(currency)) > 0),
    CONSTRAINT trips_interests_check CHECK (
        interests <@ ARRAY [
            'history',
            'nature',
            'food',
            'architecture'
            ] :: TEXT[]
        ),
    CONSTRAINT trips_time_zone_check CHECK (CHAR_LENGTH(TRIM(time_zone)) > 0),
    CONSTRAINT trips_updated_at_check CHECK (updated_at >= created_at),
    CONSTRAINT trips_deleted_at_check CHECK (
        deleted_at IS NULL
            OR deleted_at >= created_at
        )
);
