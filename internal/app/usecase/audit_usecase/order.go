package audit_usecase

import "github.com/Housiadas/backend-system/internal/core/domain/audit"

var orderByFields = map[string]string{
	"obj_id":     audit.OrderByObjID,
	"obj_domain": audit.OrderByObjDomain,
	"obj_name":   audit.OrderByObjName,
	"actor_id":   audit.OrderByActorID,
	"action":     audit.OrderByAction,
}
