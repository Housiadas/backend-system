// Package audit_repo contains auditDB-related CRUD functionality.
package audit_repo

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

// queries
var (
	//go:embed query/audit_create.sql
	auditCreateSql string
	//go:embed query/audit_query.sql
	auditQuerySql string
	//go:embed query/audit_count.sql
	auditCountSql string
)

// Store manages the set of APIs for auditDB database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new auditDB record into the database.
func (s *Store) Create(ctx context.Context, a audit.Audit) error {
	dbAudit, err := toDBAudit(a)
	if err != nil {
		return err
	}

	if err := pgsql.NamedExecContext(ctx, s.log, s.db, auditCreateSql, dbAudit); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(
	ctx context.Context,
	filter audit.QueryFilter,
	orderBy order.By,
	page page.Page,
) ([]audit.Audit, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	buf := bytes.NewBufferString(auditQuerySql)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbAudits []auditDB
	if err := pgsql.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbAudits); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusAudits(dbAudits)
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter audit.QueryFilter) (int, error) {
	data := map[string]any{}

	buf := bytes.NewBufferString(auditCountSql)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := pgsql.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}
