package currency

// Currency Constants for all supported currencies
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	CAD Currency = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency Currency) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
