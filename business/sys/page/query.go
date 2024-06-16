package page

import (
	"encoding/json"
)

// Result is the data model used when returning a query result.
type Result[T any] struct {
	Data     []T      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// NewResult constructs a result value to return query results.
func NewResult[T any](data []T, total int, page Page) Result[T] {
	metadata := calculateMetadata(total, page.Number(), page.RowsPerPage())
	return Result[T]{
		Data:     data,
		Metadata: metadata,
	}
}

// Encode implments the encoder interface.
func (r Result[T]) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}
