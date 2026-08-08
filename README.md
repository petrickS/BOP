# BOP

Поиск самого выгодного оптового предложения по стройматериалам.

## Запуск локально

Нужны Docker и Go 1.26+.

```bash
# 1. Поднять Postgres и Redis
docker compose up -d

# 2. Прогнать миграции
go tool goose -dir internal/core/db/migrations postgres \
  "postgres://bop:bop@localhost:5432/bop?sslmode=disable" up

# 3. Засеять тестовые данные (те же, что у автора)
docker compose exec -T postgres psql -U bop -d bop < internal/core/db/seed.sql

# 4. Запустить сервер
go run ./cmd/server
```

Открыть http://localhost:8080/ — форма поиска.
Пример: товар `D500`, количество `2500`, регион `Московская область`.

## Импорт цен из CSV

Формат файла: `supplier_id,product_id,price,in_stock,source` (см. пример в
`internal/features/parser/testdata/prices_example.csv`). Поставщик и товар
должны уже существовать в БД — импорт только обновляет цены (upsert по паре
supplier_id+product_id).

```bash
go run ./cmd/import path/to/prices.csv
```

## Структура

- `cmd/server` — HTTP-сервер (Chi + HTMX-шаблоны)
- `cmd/import` — CLI-импорт прайсов из CSV
- `internal/features/pricing` — расчёт лучшей цены
- `internal/features/search` — поиск (интерфейс `Provider` + Postgres-реализация)
- `internal/features/parser` — импорт цен из CSV
- `internal/core/db/migrations` — схема БД (goose)
- `internal/core/db/seed.sql` — тестовые данные для локальной разработки
