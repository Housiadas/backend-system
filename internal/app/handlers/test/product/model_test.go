package product_test

import (
	"time"

	"github.com/Housiadas/backend-system/internal/app/service/productapp"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
)

func toAppProduct(prd product.Product) productapp.Product {
	return productapp.Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProductPtr(prd product.Product) *productapp.Product {
	appPrd := toAppProduct(prd)
	return &appPrd
}

func toAppProducts(prds []product.Product) []productapp.Product {
	items := make([]productapp.Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
