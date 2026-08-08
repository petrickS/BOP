package search

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/petrickS/BOP/internal/features/pricing"
)

type PostgresProvider struct {
	pool *pgxpool.Pool
}

func NewPostgresProvider(pool *pgxpool.Pool) PostgresProvider {
	return PostgresProvider{pool: pool}
}

// подбирает товар по частичному совпадению названия; полноценную
// нормализацию (normalized_key) можно добавить позже.
const searchQuery = `
SELECT o.supplier_id, o.price, o.in_stock, dr.base_cost, dr.free_from_sum
FROM offers o
JOIN products p ON p.id = o.product_id
LEFT JOIN delivery_rules dr ON dr.supplier_id = o.supplier_id AND dr.region = $2
WHERE p.name ILIKE '%' || $1 || '%'
`

func (p PostgresProvider) Search(ctx context.Context, req Request) ([]pricing.Offer, error) {
	rows, err := p.pool.Query(ctx, searchQuery, req.Product, req.Region)
	if err != nil {
		return nil, fmt.Errorf("query offers: %w", err)
	}
	defer rows.Close()

	var offers []pricing.Offer
	for rows.Next() {
		var (
			supplierID  int64
			price       float64
			inStock     bool
			baseCost    *float64
			freeFromSum *float64
		)
		if err := rows.Scan(&supplierID, &price, &inStock, &baseCost, &freeFromSum); err != nil {
			return nil, fmt.Errorf("scan offer: %w", err)
		}

		orderSum := price * float64(req.Qty)
		deliveryCost := 0.0
		if baseCost != nil {
			deliveryCost = *baseCost
		}
		if freeFromSum != nil && orderSum >= *freeFromSum {
			deliveryCost = 0
		}

		offers = append(offers, pricing.Offer{
			SupplierId:   supplierID,
			Price:        price,
			Qty:          req.Qty,
			DeliveryCost: deliveryCost,
			InStock:      inStock,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate offers: %w", err)
	}

	return offers, nil
}
