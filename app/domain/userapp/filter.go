package userapp

import (
	"net/mail"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/types/name"
	"github.com/Housiadas/backend-system/business/sys/validation"
)

func parseFilter(qp QueryParams) (userbus.QueryFilter, error) {
	var filter userbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return userbus.QueryFilter{}, validation.NewFieldsError("user_id", err)
		}
		filter.ID = &id
	}

	if qp.Name != "" {
		n, err := name.Parse(qp.Name)
		if err != nil {
			return userbus.QueryFilter{}, validation.NewFieldsError("name", err)
		}
		filter.Name = &n
	}

	if qp.Email != "" {
		addr, err := mail.ParseAddress(qp.Email)
		if err != nil {
			return userbus.QueryFilter{}, validation.NewFieldsError("email", err)
		}
		filter.Email = addr
	}

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		if err != nil {
			return userbus.QueryFilter{}, validation.NewFieldsError("start_created_date", err)
		}
		filter.StartCreatedDate = &t
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		if err != nil {
			return userbus.QueryFilter{}, validation.NewFieldsError("end_created_date", err)
		}
		filter.EndCreatedDate = &t
	}

	return filter, nil
}
