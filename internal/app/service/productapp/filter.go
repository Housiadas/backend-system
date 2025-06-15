package productapp

import (
	"strconv"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/name"
	"github.com/Housiadas/backend-system/internal/core/domain/product"
)

type AppQueryParams struct {
	Page     string
	Rows     string
	OrderBy  string
	ID       string
	Name     string
	Cost     string
	Quantity string
}

func parseFilter(qp AppQueryParams) (product.QueryFilter, error) {
	var fieldErrors validation.FieldErrors
	var filter product.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("product_id", err)
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

	if qp.Cost != "" {
		cst, err := strconv.ParseFloat(qp.Cost, 64)
		switch err {
		case nil:
			filter.Cost = &cst
		default:
			fieldErrors.Add("cost", err)
		}
	}

	if qp.Quantity != "" {
		qua, err := strconv.ParseInt(qp.Quantity, 10, 64)
		switch err {
		case nil:
			i := int(qua)
			filter.Quantity = &i
		default:
			fieldErrors.Add("quantity", err)
		}
	}

	if fieldErrors != nil {
		return product.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
