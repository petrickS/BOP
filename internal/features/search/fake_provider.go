package search

import (
	"context"

	"github.com/petrickS/BOP/internal/features/pricing"
)

type FakeProvider struct{}

func (FakeProvider) Search(ctx context.Context, req Request) ([]pricing.Offer, error) {
	return []pricing.Offer{
		{SupplierId: 1, Price: 69, Qty: req.Qty, DeliveryCost: 12000, InStock: true},
		{SupplierId: 2, Price: 68, Qty: req.Qty, DeliveryCost: 19000, InStock: true},
		{SupplierId: 3, Price: 71, Qty: req.Qty, DeliveryCost: 0, InStock: true},
	}, nil
}
