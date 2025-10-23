package dto

type PricingConfig struct {
	PricePerUnitOfDistance float64
	PricingPerUnitOftime   float64
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerUnitOfDistance: 1.5,
		PricingPerUnitOftime:   0.25,
	}
}
