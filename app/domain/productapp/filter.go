package productapp

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/domain/productbus"
	"github.com/Housiadas/backend-system/internal/sys/types/name"
	"github.com/Housiadas/backend-system/internal/sys/validation"
)

type QueryParams struct {
	Page     string
	Rows     string
	OrderBy  string
	ID       string
	Name     string
	Cost     string
	Quantity string
}

func ParseQueryParams(r *http.Request) QueryParams {
	values := r.URL.Query()

	return QueryParams{
		Page:     values.Get("page"),
		Rows:     values.Get("rows"),
		OrderBy:  values.Get("orderBy"),
		ID:       values.Get("product_id"),
		Name:     values.Get("name"),
		Cost:     values.Get("cost"),
		Quantity: values.Get("quantity"),
	}

}

func parseFilter(qp QueryParams) (productbus.QueryFilter, error) {
	var fieldErrors validation.FieldErrors
	var filter productbus.QueryFilter

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
		return productbus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
