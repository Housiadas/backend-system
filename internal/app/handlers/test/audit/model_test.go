package audit_test

import (
	"time"

	"github.com/Housiadas/backend-system/internal/app/usecase/audit_usecase"
	"github.com/Housiadas/backend-system/internal/core/domain/audit"
)

func toAppAudit(bus audit.Audit) audit_usecase.Audit {
	return audit_usecase.Audit{
		ID:        bus.ID.String(),
		ObjID:     bus.ObjID.String(),
		ObjEntity: bus.ObjEntity.String(),
		ObjName:   bus.ObjName.String(),
		ActorID:   bus.ActorID.String(),
		Action:    bus.Action,
		Data:      string(bus.Data),
		Message:   bus.Message,
		Timestamp: bus.Timestamp.Format(time.RFC3339),
	}
}

func toAppAudits(audits []audit.Audit) []audit_usecase.Audit {
	app := make([]audit_usecase.Audit, len(audits))
	for i, adt := range audits {
		app[i] = toAppAudit(adt)
	}

	return app
}
