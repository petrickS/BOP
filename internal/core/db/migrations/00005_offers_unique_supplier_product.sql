-- +goose Up
-- один активный оффер на пару поставщик+товар; повторный импорт обновляет цену, а не плодит дубли
ALTER TABLE offers ADD CONSTRAINT offers_supplier_product_unique UNIQUE (supplier_id, product_id);

-- +goose Down
ALTER TABLE offers DROP CONSTRAINT offers_supplier_product_unique;
