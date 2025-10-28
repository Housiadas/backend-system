package product_test

import (
	"time"

	"github.com/Housiadas/backend-system/internal/app/usecase/product_usecase"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
)

func toAppProduct(prd product.Product) product_usecase.Product {
	return product_usecase.Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProductPtr(prd product.Product) *product_usecase.Product {
	appPrd := toAppProduct(prd)
	return &appPrd
}

func toAppProducts(prds []product.Product) []product_usecase.Product {
	items := make([]product_usecase.Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
