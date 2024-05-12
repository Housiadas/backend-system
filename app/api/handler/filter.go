package handler

import (
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/sys/page"
)

func productParseQueryParams(r *http.Request) (productapp.QueryParams, error) {
	const (
		orderBy          = "orderBy"
		filterPage       = "page"
		filterRow        = "row"
		filterByProdID   = "product_id"
		filterByCost     = "cost"
		filterByQuantity = "quantity"
		filterByName     = "name"
	)

	values := r.URL.Query()

	var filter productapp.QueryParams

	pg, err := page.ParseHTTP(r)
	if err != nil {
		return productapp.QueryParams{}, err
	}
	filter.Page = pg.Number
	filter.Rows = pg.RowsPerPage

	if orderBy := values.Get(orderBy); orderBy != "" {
		filter.OrderBy = orderBy
	}

	if productID := values.Get(filterByProdID); productID != "" {
		filter.ID = productID
	}

	if cost := values.Get(filterByCost); cost != "" {
		filter.Cost = cost
	}

	if quantity := values.Get(filterByQuantity); quantity != "" {
		filter.Quantity = quantity
	}

	if name := values.Get(filterByName); name != "" {
		filter.Name = name
	}

	return filter, nil
}

func userParseQueryParams(r *http.Request) (userapp.QueryParams, error) {
	const (
		orderBy                  = "orderBy"
		filterPage               = "page"
		filterRow                = "row"
		filterByUserID           = "user_id"
		filterByEmail            = "email"
		filterByStartCreatedDate = "start_created_date"
		filterByEndCreatedDate   = "end_created_date"
		filterByName             = "name"
	)

	values := r.URL.Query()

	var filter userapp.QueryParams

	pg, err := page.ParseHTTP(r)
	if err != nil {
		return userapp.QueryParams{}, err
	}
	filter.Page = pg.Number
	filter.Rows = pg.RowsPerPage

	if orderBy := values.Get(orderBy); orderBy != "" {
		filter.OrderBy = orderBy
	}

	if userID := values.Get(filterByUserID); userID != "" {
		filter.ID = userID
	}

	if email := values.Get(filterByEmail); email != "" {
		filter.Email = email
	}

	if startedDate := values.Get(filterByStartCreatedDate); startedDate != "" {
		filter.StartCreatedDate = startedDate
	}

	if endDate := values.Get(filterByStartCreatedDate); endDate != "" {
		filter.EndCreatedDate = endDate
	}

	return filter, nil
}
