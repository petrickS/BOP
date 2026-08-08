package search

import (
	"context"

	"github.com/petrickS/BOP/internal/features/pricing"
)

type Request struct {
	Product string `json:"product"`
	Qty     int    `json:"qty"`
	Region  string `json:"region"`
}

type Provider interface {
	Search(ctx context.Context, req Request) ([]pricing.Offer, error)
}

type Result struct {
	Offers  []pricing.Offer `json:"offers"`
	Best    pricing.Offer   `json:"best"`
	Savings float64         `json:"savings"`
}

func Handle(ctx context.Context, p Provider, req Request) (Result, error) {
	offers, err := p.Search(ctx, req)
	if err != nil {
		return Result{}, err
	}
	best, savings := pricing.Best(offers)
	return Result{Offers: offers, Best: best, Savings: savings}, nil
}
