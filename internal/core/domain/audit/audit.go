package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/domain/entity"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
)

// Audit represents information about an individual audit record.
type Audit struct {
	ID        uuid.UUID
	ObjID     uuid.UUID
	ObjEntity entity.Entity
	ObjName   name.Name
	ActorID   uuid.UUID
	Action    string
	Data      json.RawMessage
	Message   string
	Timestamp time.Time
}

// NewAudit represents the information needed to create a new audit record.
type NewAudit struct {
	ObjID     uuid.UUID
	ObjEntity entity.Entity
	ObjName   name.Name
	ActorID   uuid.UUID
	Action    string
	Data      any
	Message   string
}

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ObjID     *uuid.UUID
	ObjEntity *entity.Entity
	ObjName   *name.Name
	ActorID   *uuid.UUID
	Action    *string
	Since     *time.Time
	Until     *time.Time
}
