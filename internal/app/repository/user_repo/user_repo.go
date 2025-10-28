// Package user_repo contains database related CRUD functionality.
package user_repo

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

// queries
var (
	//go:embed query/user_create.sql
	userCreateSql string
	//go:embed query/user_update.sql
	userUpdateSql string
	//go:embed query/user_delete.sql
	userDeleteSql string
	//go:embed query/user_query.sql
	userQuerySql string
	//go:embed query/user_query_by_id.sql
	userQueryByIdSql string
	//go:embed query/user_query_by_email.sql
	userQueryByEmailSql string
	//go:embed query/user_count.sql
	userCountSql string
)

// Store manages the set of APIs for userDB database access.
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
func (s *Store) NewWithTx(tx pgsql.CommitRollbacker) (user.Storer, error) {
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

// Create inserts a new userDB into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {
	if err := pgsql.NamedExecContext(ctx, s.log, s.db, userCreateSql, toUserDB(usr)); err != nil {
		if errors.Is(err, pgsql.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a userDB document in the database.
func (s *Store) Update(ctx context.Context, usr user.User) error {
	if err := pgsql.NamedExecContext(ctx, s.log, s.db, userUpdateSql, toUserDB(usr)); err != nil {
		if errors.Is(err, pgsql.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a userDB from the database.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	if err := pgsql.NamedExecContext(ctx, s.log, s.db, userDeleteSql, toUserDB(usr)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, page page.Page) ([]user.User, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	buf := bytes.NewBufferString(userQuerySql)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbUsrs []userDB
	if err := pgsql.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbUsrs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toUsersDomain(dbUsrs)
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter user.QueryFilter) (int, error) {
	data := map[string]any{}
	buf := bytes.NewBufferString(userCountSql)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := pgsql.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified userDB from the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	var dbUsr userDB
	if err := pgsql.NamedQueryStruct(ctx, s.log, s.db, userQueryByIdSql, data, &dbUsr); err != nil {
		if errors.Is(err, pgsql.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("db: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("db: %w", err)
	}

	return toUserDomain(dbUsr)
}

// QueryByEmail gets the specified userDB from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	var dbUsr userDB
	if err := pgsql.NamedQueryStruct(ctx, s.log, s.db, userQueryByEmailSql, data, &dbUsr); err != nil {
		if errors.Is(err, pgsql.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("db: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("db: %w", err)
	}

	return toUserDomain(dbUsr)
}
