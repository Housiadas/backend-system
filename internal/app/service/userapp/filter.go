package userapp

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
)

type QueryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	ID               string
	Name             string
	Email            string
	StartCreatedDate string
	EndCreatedDate   string
}

func ParseQueryParams(r *http.Request) QueryParams {
	values := r.URL.Query()

	return QueryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("rows"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("user_id"),
		Name:             values.Get("name"),
		Email:            values.Get("email"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}
}

func parseFilter(qp QueryParams) (user.QueryFilter, error) {
	var fieldErrors validation.FieldErrors
	var filter user.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("user_id", err)
		}
	}

	if qp.Name != "" {
		n, err := name.Parse(qp.Name)
		switch err {
		case nil:
			filter.Name = &n
		default:
			fieldErrors.Add("name", err)
		}
	}

	if qp.Email != "" {
		addr, err := mail.ParseAddress(qp.Email)
		switch err {
		case nil:
			filter.Email = addr
		default:
			fieldErrors.Add("email", err)
		}
	}

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		switch err {
		case nil:
			filter.StartCreatedDate = &t
		default:
			fieldErrors.Add("start_created_date", err)
		}
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		switch err {
		case nil:
			filter.EndCreatedDate = &t
		default:
			fieldErrors.Add("end_created_date", err)
		}
	}

	if fieldErrors != nil {
		return user.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
