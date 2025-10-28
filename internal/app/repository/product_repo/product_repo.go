// Package product_repo contains productDB related CRUD functionality.
package product_repo

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/internal/core/domain/product"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

// Store manages the set of APIs for productDB database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx pgsql.CommitRollbacker) (product.Storer, error) {
	ec, err := pgsql.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Create adds a Product to the pgsql. It returns the created Product with
// fields like ID and DateCreated populated.
func (s *Store) Create(ctx context.Context, prd product.Product) error {
	const q = `
	INSERT INTO products
		(product_id, user_id, name, cost, quantity, date_created, date_updated)
	VALUES
		(:product_id, :user_id, :name, :cost, :quantity, :date_created, :date_updated)`

	if err := pgsql.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update modifies data about a product. It will error if the specified ID is
// invalid or does not reference an existing product.
func (s *Store) Update(ctx context.Context, prd product.Product) error {
	const q = `
	UPDATE
		products
	SET
		"name" = :name,
		"cost" = :cost,
		"quantity" = :quantity,
		"date_updated" = :date_updated
	WHERE
		product_id = :product_id`

	if err := pgsql.NamedExecContext(ctx, s.log, s.db, q, toDBProduct(prd)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes the productDB identified by a given ID.
func (s *Store) Delete(ctx context.Context, prd product.Product) error {
	data := struct {
		ID string `db:"product_id"`
	}{
		ID: prd.ID.String(),
	}

	const q = `
	DELETE FROM
		products
	WHERE
		product_id = :product_id`

	if err := pgsql.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query gets all Products from the database.
func (s *Store) Query(ctx context.Context, filter product.QueryFilter, orderBy order.By, page page.Page) ([]product.Product, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
	    product_id, user_id, name, cost, quantity, date_created, date_updated
	FROM
		products`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbPrds []productDB
	if err := pgsql.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusProducts(dbPrds)
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter product.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		products`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count   int `db:"count"`
		Sold    int `db:"sold"`
		Revenue int `db:"revenue"`
	}
	if err := pgsql.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID finds the productDB identified by a given ID.
func (s *Store) QueryByID(ctx context.Context, productID uuid.UUID) (product.Product, error) {
	data := struct {
		ID string `db:"product_id"`
	}{
		ID: productID.String(),
	}

	const q = `
	SELECT
	    product_id, user_id, name, cost, quantity, date_created, date_updated
	FROM
		products
	WHERE
		product_id = :product_id`

	var dbPrd productDB
	if err := pgsql.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPrd); err != nil {
		if errors.Is(err, pgsql.ErrDBNotFound) {
			return product.Product{}, fmt.Errorf("db: %w", product.ErrNotFound)
		}
		return product.Product{}, fmt.Errorf("db: %w", err)
	}

	return toBusProduct(dbPrd)
}

// QueryByUserID finds the productDB identified by a given User ID.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]product.Product, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    product_id, user_id, name, cost, quantity, date_created, date_updated
	FROM
		products
	WHERE
		user_id = :user_id`

	var dbPrds []productDB
	if err := pgsql.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbPrds); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toBusProducts(dbPrds)
}
