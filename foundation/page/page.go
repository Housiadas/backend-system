// Package page provides support for query pagination.
package page

import (
	"fmt"
)

// Page represents the requested page and rows per page.
type Page struct {
	number int
	rows   int
}

// String implements the stringer interface.
func (p Page) String() string {
	return fmt.Sprintf("page: %d rows: %d", p.number, p.rows)
}

// Number returns the page number.
func (p Page) Number() int {
	return p.number
}

// RowsPerPage returns the rows per page.
func (p Page) RowsPerPage() int {
	return p.rows
}
