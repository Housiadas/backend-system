package audit_repo

import (
	"fmt"

	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/pkg/order"
)

var orderByFields = map[string]string{
	audit.OrderByObjID:     "obj_id",
	audit.OrderByObjDomain: "obj_domain",
	audit.OrderByObjName:   "obj_name",
	audit.OrderByActorID:   "actor_id",
	audit.OrderByAction:    "action",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
