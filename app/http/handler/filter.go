package handler

import (
	"net/http"

	"github.com/Housiadas/backend-system/app/domain/productapp"
	"github.com/Housiadas/backend-system/app/domain/userapp"
)

func productParseQueryParams(r *http.Request) (productapp.QueryParams, error) {
	values := r.URL.Query()

	filter := productapp.QueryParams{
		Page:     values.Get("page"),
		Rows:     values.Get("row"),
		OrderBy:  values.Get("orderBy"),
		ID:       values.Get("product_id"),
		Name:     values.Get("name"),
		Cost:     values.Get("cost"),
		Quantity: values.Get("quantity"),
	}

	return filter, nil
}

func userParseQueryParams(r *http.Request) (userapp.QueryParams, error) {
	values := r.URL.Query()

	filter := userapp.QueryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("row"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("user_id"),
		Name:             values.Get("name"),
		Email:            values.Get("email"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}

	return filter, nil
}
