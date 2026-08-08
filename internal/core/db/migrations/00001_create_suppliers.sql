-- +goose Up
CREATE TABLE suppliers (
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT NOT NULL,
    region   TEXT NOT NULL,
    contacts TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE suppliers;
