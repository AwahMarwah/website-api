package rajaongkir

type Provider interface {
	SearchDestination(search string) ([]Destination, error)
	CalculateCost(req CalculateCostRequest) ([]ShippingOption, error)
}
