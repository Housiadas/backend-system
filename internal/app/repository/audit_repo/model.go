package audit_repo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/internal/core/domain/entity"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
)

type auditDB struct {
	ID        uuid.UUID          `db:"id"`
	ObjID     uuid.UUID          `db:"obj_id"`
	ObjEntity string             `db:"obj_entity"`
	ObjName   string             `db:"obj_name"`
	ActorID   uuid.UUID          `db:"actor_id"`
	Action    string             `db:"action"`
	Data      types.NullJSONText `db:"data"`
	Message   string             `db:"message"`
	Timestamp time.Time          `db:"timestamp"`
}

func toDBAudit(bus audit.Audit) (auditDB, error) {
	db := auditDB{
		ID:        bus.ID,
		ObjID:     bus.ObjID,
		ObjEntity: bus.ObjEntity.String(),
		ObjName:   bus.ObjName.String(),
		ActorID:   bus.ActorID,
		Action:    bus.Action,
		Data:      types.NullJSONText{JSONText: []byte(bus.Data), Valid: true},
		Message:   bus.Message,
		Timestamp: bus.Timestamp.UTC(),
	}

	return db, nil
}

func toDomainAudit(db auditDB) (audit.Audit, error) {
	ent, err := entity.Parse(db.ObjEntity)
	if err != nil {
		return audit.Audit{}, fmt.Errorf("parse ent: %w", err)
	}

	n, err := name.Parse(db.ObjName)
	if err != nil {
		return audit.Audit{}, fmt.Errorf("parse name: %w", err)
	}

	bus := audit.Audit{
		ID:        db.ID,
		ObjID:     db.ObjID,
		ObjEntity: ent,
		ObjName:   n,
		ActorID:   db.ActorID,
		Action:    db.Action,
		Data:      json.RawMessage(db.Data.JSONText),
		Message:   db.Message,
		Timestamp: db.Timestamp.Local(),
	}

	return bus, nil
}

func toBusAudits(dbs []auditDB) ([]audit.Audit, error) {
	audits := make([]audit.Audit, len(dbs))

	for i, db := range dbs {
		a, err := toDomainAudit(db)
		if err != nil {
			return nil, err
		}

		audits[i] = a
	}

	return audits, nil
}
