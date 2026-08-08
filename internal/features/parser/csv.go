package parser

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ожидаемые колонки CSV: supplier_id,product_id,price,in_stock,source
// (supplier и product должны уже существовать — импорт только обновляет цены)

const upsertOfferQuery = `
INSERT INTO offers (supplier_id, product_id, price, in_stock, source, updated_at)
VALUES ($1, $2, $3, $4, $5, now())
ON CONFLICT (supplier_id, product_id)
DO UPDATE SET price = excluded.price,
              in_stock = excluded.in_stock,
              source = excluded.source,
              updated_at = now()
`

func ImportCSV(ctx context.Context, pool *pgxpool.Pool, r io.Reader) (int, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 5

	header, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("read header: %w", err)
	}
	if err := validateHeader(header); err != nil {
		return 0, err
	}

	count := 0
	for lineNum := 2; ; lineNum++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("line %d: %w", lineNum, err)
		}

		supplierID, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return count, fmt.Errorf("line %d: invalid supplier_id %q: %w", lineNum, record[0], err)
		}
		productID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return count, fmt.Errorf("line %d: invalid product_id %q: %w", lineNum, record[1], err)
		}
		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return count, fmt.Errorf("line %d: invalid price %q: %w", lineNum, record[2], err)
		}
		inStock, err := strconv.ParseBool(record[3])
		if err != nil {
			return count, fmt.Errorf("line %d: invalid in_stock %q: %w", lineNum, record[3], err)
		}
		source := record[4]

		if _, err := pool.Exec(ctx, upsertOfferQuery, supplierID, productID, price, inStock, source); err != nil {
			return count, fmt.Errorf("line %d: upsert offer: %w", lineNum, err)
		}
		count++
	}

	return count, nil
}

func validateHeader(header []string) error {
	want := []string{"supplier_id", "product_id", "price", "in_stock", "source"}
	if len(header) != len(want) {
		return fmt.Errorf("expected %d columns, got %d", len(want), len(header))
	}
	for i, col := range want {
		if header[i] != col {
			return fmt.Errorf("expected column %q at position %d, got %q", col, i, header[i])
		}
	}
	return nil
}
