// Package auditrepo contains auditDB-related CRUD functionality.
package auditrepo

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/sqldb"
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
	const q = `
	INSERT INTO audit
		(id, obj_id, obj_entity, obj_name, actor_id, action, data, message, timestamp)
	VALUES
		(:id, :obj_id, :obj_entity, :obj_name, :actor_id, :action, :data, :message, :timestamp)`

	dbAudit, err := toDBAudit(a)
	if err != nil {
		return err
	}

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbAudit); err != nil {
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

	const q = `
	SELECT
		id, obj_id, obj_entity, obj_name, actor_id, action, data, message, timestamp
	FROM
		audit
	`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbAudits []auditDB
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbAudits); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusAudits(dbAudits)
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter audit.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		audit`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}
