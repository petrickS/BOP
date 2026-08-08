-- +goose Up
CREATE TABLE products (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL,
    category       TEXT NOT NULL,
    unit           TEXT NOT NULL,
    normalized_key TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE products;
