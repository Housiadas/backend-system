package product_repo

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/money"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/internal/core/domain/quantity"
)

type productDB struct {
	ID          uuid.UUID `db:"product_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Cost        float64   `db:"cost"`
	Quantity    int       `db:"quantity"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBProduct(bus product.Product) productDB {
	db := productDB{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Name:        bus.Name.String(),
		Cost:        bus.Cost.Value(),
		Quantity:    bus.Quantity.Value(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return db
}

func toBusProduct(db productDB) (product.Product, error) {
	n, err := name.Parse(db.Name)
	if err != nil {
		return product.Product{}, fmt.Errorf("parse name: %w", err)
	}

	bus := product.Product{
		ID:          db.ID,
		UserID:      db.UserID,
		Name:        n,
		Cost:        money.MustParse(db.Cost),
		Quantity:    quantity.MustParse(db.Quantity),
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusProducts(dbs []productDB) ([]product.Product, error) {
	bus := make([]product.Product, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusProduct(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
