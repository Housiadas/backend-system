package audit_test

import (
	"time"

	"github.com/Housiadas/backend-system/internal/app/service/auditapp"
	"github.com/Housiadas/backend-system/internal/core/domain/audit"
)

func toAppAudit(bus audit.Audit) auditapp.Audit {
	return auditapp.Audit{
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

func toAppAudits(audits []audit.Audit) []auditapp.Audit {
	app := make([]auditapp.Audit, len(audits))
	for i, adt := range audits {
		app[i] = toAppAudit(adt)
	}

	return app
}
