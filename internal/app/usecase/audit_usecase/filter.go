package audit_usecase

import (
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/internal/core/domain/entity"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
)

type AppQueryParams struct {
	Page      string
	Rows      string
	OrderBy   string
	ObjID     string
	ObjEntity string
	ObjName   string
	ActorID   string
	Action    string
	Since     string
	Until     string
}

func parseFilter(qp AppQueryParams) (audit.QueryFilter, error) {
	var fieldErrors validation.FieldErrors
	var filter audit.QueryFilter

	if qp.ObjID != "" {
		id, err := uuid.Parse(qp.ObjID)
		switch err {
		case nil:
			filter.ObjID = &id
		default:
			fieldErrors.Add("obj_id", err)
		}
	}

	if qp.ObjEntity != "" {
		domain, err := entity.Parse(qp.ObjEntity)
		switch err {
		case nil:
			filter.ObjEntity = &domain
		default:
			fieldErrors.Add("obj_domain", err)
		}
	}

	if qp.ObjName != "" {
		name, err := name.Parse(qp.ObjName)
		switch err {
		case nil:
			filter.ObjName = &name
		default:
			fieldErrors.Add("obj_name", err)
		}
	}

	if qp.ActorID != "" {
		id, err := uuid.Parse(qp.ActorID)
		switch err {
		case nil:
			filter.ActorID = &id
		default:
			fieldErrors.Add("actor_id", err)
		}
	}

	if qp.Action != "" {
		filter.Action = &qp.Action
	}

	if qp.Since != "" {
		t, err := time.Parse(time.RFC3339, qp.Since)
		switch err {
		case nil:
			filter.Since = &t
		default:
			fieldErrors.Add("since", err)
		}
	}

	if qp.Until != "" {
		t, err := time.Parse(time.RFC3339, qp.Until)
		switch err {
		case nil:
			filter.Until = &t
		default:
			fieldErrors.Add("until", err)
		}
	}

	if fieldErrors != nil {
		return audit.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
