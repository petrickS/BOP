-- +goose Up
CREATE TABLE offers (
    id            BIGSERIAL PRIMARY KEY,
    supplier_id   BIGINT NOT NULL REFERENCES suppliers (id) ON DELETE CASCADE,
    product_id    BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    price         NUMERIC(12, 2) NOT NULL,
    min_order_qty INTEGER NOT NULL DEFAULT 1,
    in_stock      BOOLEAN NOT NULL DEFAULT TRUE,
    source        TEXT NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_offers_product_id ON offers (product_id);
CREATE INDEX idx_offers_supplier_id ON offers (supplier_id);

-- +goose Down
DROP TABLE offers;
