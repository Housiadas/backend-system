package transaction_usecase

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/money"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/domain/quantity"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/pkg/errs"
)

// Product represents an individual product.
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

// =============================================================================

// NewTran represents an example of cross-core transaction at the
// application layer.
type NewTran struct {
	Product NewProduct `json:"product"`
	User    NewUser    `json:"user"`
}

// Validate checks the data in the model is considered clean.
func (app NewTran) Validate() error {
	if err := validation.Check(app); err != nil {
		return errs.Newf(errs.InvalidArgument, "validation: %s", err)
	}

	return nil
}

// Decode implements the decoder interface.
func (app *NewTran) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// =============================================================================

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

// Validate checks the data in the model is considered clean.
func (app NewUser) Validate() error {
	if err := validation.Check(app); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	return nil
}

func toBusNewUser(app NewUser) (user.NewUser, error) {
	roles, err := role.ParseMany(app.Roles)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(app.Name)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	department, err := name.ParseNull(app.Department)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := user.NewUser{
		Name:       nme,
		Email:      *addr,
		Roles:      roles,
		Department: department,
		Password:   app.Password,
	}

	return bus, nil
}

// =============================================================================

// NewProduct is what we require from clients when adding a Product.
type NewProduct struct {
	Name     string  `json:"name" validate:"required"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
	Quantity int     `json:"quantity" validate:"required,gte=1"`
}

// Validate checks the data in the model is considered clean.
func (app NewProduct) Validate() error {
	if err := validation.Check(app); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	return nil
}

func toBusNewProduct(app NewProduct) (product.NewProduct, error) {
	n, err := name.Parse(app.Name)
	if err != nil {
		return product.NewProduct{}, fmt.Errorf("parse: %w", err)
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
		Name:     n,
		Cost:     cost,
		Quantity: q,
	}

	return bus, nil
}
