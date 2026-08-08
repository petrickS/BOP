package pricing

type Offer struct {
	SupplierId   int64   `json:"supplier_id"`
	Price        float64 `json:"price"`
	Qty          int     `json:"qty"`
	DeliveryCost float64 `json:"delivery_cost"`
	InStock      bool    `json:"in_stock"`
}

func (o Offer) Total() float64 {
	sum := o.Price*float64(o.Qty) + o.DeliveryCost
	return sum
}

// находит лучшую цену из слайса
func Best(offers []Offer) (best Offer, savings float64) {
	if len(offers) == 0 {
		return Offer{}, 0
	}

	var sum float64
	best = offers[0]
	bestTotal := best.Total()

	for _, offer := range offers {
		total := offer.Total()
		sum += total
		if total < bestTotal {
			best = offer
			bestTotal = total
		}
	}

	avg := sum / float64(len(offers))
	savings = avg - bestTotal

	return best, savings
}
