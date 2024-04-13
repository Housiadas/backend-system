// Package page provides support for query paging.
package page

import (
	"net/http"
	"strconv"

	"github.com/Housiadas/backend-system/foundation/validate"
)

// Page represents the requested page and rows per page.
type Page struct {
	Number      int
	RowsPerPage int
}

// ParseHTTP parses the request for the page and rows query string. The
// defaults are provided as well.
func ParseHTTP(r *http.Request) (Page, error) {
	values := r.URL.Query()

	number := 1
	if page := values.Get("page"); page != "" {
		var err error
		number, err = strconv.Atoi(page)
		if err != nil {
			return Page{}, validate.NewFieldsError("page", err)
		}
	}

	rowsPerPage := 10
	if rows := values.Get("rows"); rows != "" {
		var err error
		rowsPerPage, err = strconv.Atoi(rows)
		if err != nil {
			return Page{}, validate.NewFieldsError("rows", err)
		}
	}

	p := Page{
		Number:      number,
		RowsPerPage: rowsPerPage,
	}

	return p, nil
}

// Document is the form used for API responses from query API calls.
type Document[T any] struct {
	Data     []T      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// NewDocument constructs a response value for a web paging trusted.
func NewDocument[T any](data []T, total int, page int, rowsPerPage int) Document[T] {
	metadata := CalculateMetadata(total, page, rowsPerPage)
	return Document[T]{
		Data:     data,
		Metadata: metadata,
	}
}
