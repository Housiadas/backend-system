package productapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ctxPck "github.com/Housiadas/backend-system/internal/common/context"
	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/money"
	namePck "github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/domain/quantity"
)

// The Product represents information about an individual product.
type Product struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Product) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppProduct(prd product.Product) Product {
	return Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProducts(prds []product.Product) []Product {
	app := make([]Product, len(prds))
	for i, prd := range prds {
		app[i] = toAppProduct(prd)
	}

	return app
}

// =============================================================================

// NewProduct defines the data needed to add a new product.
type NewProduct struct {
	Name     string  `json:"name" validate:"required"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
	Quantity int     `json:"quantity" validate:"required,gte=1"`
}

// Decode implements the decoder interface.
func (app *NewProduct) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app *NewProduct) Validate() error {
	if err := validation.Check(app); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	return nil
}

func toBusNewProduct(ctx context.Context, app NewProduct) (product.NewProduct, error) {
	userID, err := ctxPck.GetUserID(ctx)
	if err != nil {
		return product.NewProduct{}, fmt.Errorf("getuserid: %w", err)
	}

	n, err := namePck.Parse(app.Name)
	if err != nil {
		return product.NewProduct{}, fmt.Errorf("parse name: %w", err)
	}

	cost, err := money.Parse(app.Cost)
	if err != nil {
		return product.NewProduct{}, fmt.Errorf("parse cost: %w", err)
	}

	q, err := quantity.Parse(app.Quantity)
	if err != nil {
		return product.NewProduct{}, fmt.Errorf("parse quantity: %w", err)
	}

	bus := product.NewProduct{
		UserID:   userID,
		Name:     n,
		Cost:     cost,
		Quantity: q,
	}

	return bus, nil
}

// =============================================================================

// UpdateProduct defines the data needed to update a product.
type UpdateProduct struct {
	Name     *string  `json:"name"`
	Cost     *float64 `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int     `json:"quantity" validate:"omitempty,gte=1"`
}

// Decode implements the decoder interface.
func (app *UpdateProduct) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app *UpdateProduct) Validate() error {
	if err := validation.Check(app); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	return nil
}

func toBusUpdateProduct(app UpdateProduct) (product.UpdateProduct, error) {
	var nme *namePck.Name
	if app.Name != nil {
		nm, err := namePck.Parse(*app.Name)
		if err != nil {
			return product.UpdateProduct{}, fmt.Errorf("parse: %w", err)
		}
		nme = &nm
	}

	var cost *money.Money
	if app.Cost != nil {
		cst, err := money.Parse(*app.Cost)
		if err != nil {
			return product.UpdateProduct{}, fmt.Errorf("parse: %w", err)
		}
		cost = &cst
	}

	var qnt *quantity.Quantity
	if app.Cost != nil {
		qn, err := quantity.Parse(*app.Quantity)
		if err != nil {
			return product.UpdateProduct{}, fmt.Errorf("parse: %w", err)
		}
		qnt = &qn
	}

	bus := product.UpdateProduct{
		Name:     nme,
		Cost:     cost,
		Quantity: qnt,
	}

	return bus, nil
}
