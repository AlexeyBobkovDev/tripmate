CREATE TYPE app.visitor_category AS ENUM(
    'adult',
    'child',
    'student',
    'senior'
);

CREATE TABLE IF NOT EXISTS app.trip_visitor_categories(
    id BIGINT NOT NULL REFERENCES app.trips(id) ON DELETE CASCADE,
    category app.visitor_category NOT NULL,
    quantity SMALLINT NOT NULL,

    PRIMARY KEY(id, category),
    CHECK(id > 0),
    CHECK(quantity >= 0)
);
