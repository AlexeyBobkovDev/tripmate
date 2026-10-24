CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS app.users (
    id BIGSERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    username VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    description VARCHAR(1000) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(16) NOT NULL,
    password_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT users_version_check CHECK(version > 0),
    CONSTRAINT users_name_check CHECK(CHAR_LENGTH(TRIM(name)) > 0),
    CONSTRAINT users_surname_check CHECK(CHAR_LENGTH(TRIM(surname)) > 0),
    CONSTRAINT users_username_check CHECK(CHAR_LENGTH(TRIM(username)) > 0),
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

CREATE UNIQUE INDEX idx_unique_user_username ON app.users(username)
WHERE
    deleted_at IS NULL;

CREATE UNIQUE INDEX idx_unique_user_email ON app.users(email)
WHERE
    deleted_at IS NULL;

CREATE UNIQUE INDEX idx_unique_user_phone_number ON app.users(phone_number)
WHERE
    deleted_at IS NULL;
