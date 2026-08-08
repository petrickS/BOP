-- +goose Up
CREATE TABLE delivery_rules (
    id             BIGSERIAL PRIMARY KEY,
    supplier_id    BIGINT NOT NULL REFERENCES suppliers (id) ON DELETE CASCADE,
    region         TEXT NOT NULL,
    base_cost      NUMERIC(12, 2) NOT NULL DEFAULT 0,
    free_from_sum  NUMERIC(12, 2)
);

CREATE INDEX idx_delivery_rules_supplier_region ON delivery_rules (supplier_id, region);

-- +goose Down
DROP TABLE delivery_rules;
