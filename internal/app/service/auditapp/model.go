package auditapp

import (
	"encoding/json"
	"time"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
)

// Audit represents information about an individual audit record.
type Audit struct {
	ID        string
	ObjID     string
	ObjEntity string
	ObjName   string
	ActorID   string
	Action    string
	Data      string
	Message   string
	Timestamp string
}

// Encode implements the encoder interface.
func (app Audit) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppAudit(aud audit.Audit) Audit {
	return Audit{
		ID:        aud.ID.String(),
		ObjID:     aud.ObjID.String(),
		ObjEntity: aud.ObjEntity.String(),
		ObjName:   aud.ObjName.String(),
		ActorID:   aud.ActorID.String(),
		Action:    aud.Action,
		Data:      string(aud.Data),
		Message:   aud.Message,
		Timestamp: aud.Timestamp.Format(time.RFC3339),
	}
}

func toAppAudits(audits []audit.Audit) []Audit {
	app := make([]Audit, len(audits))
	for i, adt := range audits {
		app[i] = toAppAudit(adt)
	}

	return app
}
