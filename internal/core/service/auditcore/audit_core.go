package auditcore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/otel"
	"github.com/Housiadas/backend-system/pkg/page"
)

// Core manages the set of APIs for audit access.
type Core struct {
	log    *logger.Logger
	storer audit.Storer
}

// NewCore constructs an audit business API for use.
func NewCore(log *logger.Logger, storer audit.Storer) *Core {
	return &Core{
		log:    log,
		storer: storer,
	}
}

// Create adds a new audit record to the system.
func (b *Core) Create(ctx context.Context, na audit.NewAudit) (audit.Audit, error) {
	ctx, span := otel.AddSpan(ctx, "business.auditbus.create")
	defer span.End()

	jsonData, err := json.Marshal(na.Data)
	if err != nil {
		return audit.Audit{}, fmt.Errorf("marshal object: %w", err)
	}

	aud := audit.Audit{
		ID:        uuid.New(),
		ObjID:     na.ObjID,
		ObjEntity: na.ObjEntity,
		ObjName:   na.ObjName,
		ActorID:   na.ActorID,
		Action:    na.Action,
		Data:      jsonData,
		Message:   na.Message,
		Timestamp: time.Now(),
	}

	if err := b.storer.Create(ctx, aud); err != nil {
		return audit.Audit{}, fmt.Errorf("create audit: %w", err)
	}

	return aud, nil
}

// Query retrieves a list of existing audit records.
func (b *Core) Query(ctx context.Context, filter audit.QueryFilter, orderBy order.By, page page.Page) ([]audit.Audit, error) {
	ctx, span := otel.AddSpan(ctx, "repo.audit.query")
	defer span.End()

	audits, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query audits: %w", err)
	}

	return audits, nil
}

// Count returns the total number of users.
func (b *Core) Count(ctx context.Context, filter audit.QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.auditbus.count")
	defer span.End()

	return b.storer.Count(ctx, filter)
}
