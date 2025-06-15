package productcore

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/money"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/domain/quantity"
)

// TestGenerateNewProducts is a helper method for testing.
func TestGenerateNewProducts(n int, userID uuid.UUID) []product.NewProduct {
	newPrds := make([]product.NewProduct, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		np := product.NewProduct{
			Name:     name.MustParse(fmt.Sprintf("Name%d", idx)),
			Cost:     money.MustParse(float64(rand.Intn(500))),
			Quantity: quantity.MustParse(rand.Intn(50)),
			UserID:   userID,
		}

		newPrds[i] = np
	}

	return newPrds
}

// TestGenerateSeedProducts is a helper method for testing.
func TestGenerateSeedProducts(ctx context.Context, n int, api *Core, userID uuid.UUID) ([]product.Product, error) {
	newPrds := TestGenerateNewProducts(n, userID)

	prds := make([]product.Product, len(newPrds))
	for i, np := range newPrds {
		prd, err := api.Create(ctx, np)
		if err != nil {
			return nil, fmt.Errorf("seeding product: idx: %d : %w", i, err)
		}

		prds[i] = prd
	}

	return prds, nil
}
