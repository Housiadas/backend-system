package api

import (
	c "github.com/Housiadas/simple-banking-system/business/currency"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return c.IsSupportedCurrency(currency)
	}
	return false
}
