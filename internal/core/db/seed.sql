-- Тестовые данные для локальной разработки.
-- Применять после миграций: см. README.md.

INSERT INTO suppliers (name, region) VALUES
    ('База А', 'Московская область'),
    ('База Б', 'Московская область'),
    ('База В', 'Московская область');

INSERT INTO products (name, category, unit, normalized_key) VALUES
    ('Газоблок D500 600x300x200', 'блоки', 'шт', 'gazoblok_d500_600x300x200');

INSERT INTO offers (supplier_id, product_id, price, in_stock, source) VALUES
    (1, 1, 69, true, 'manual'),
    (2, 1, 68, true, 'manual'),
    (3, 1, 71, true, 'manual');

INSERT INTO delivery_rules (supplier_id, region, base_cost) VALUES
    (1, 'Московская область', 12000),
    (2, 'Московская область', 19000),
    (3, 'Московская область', 0);
