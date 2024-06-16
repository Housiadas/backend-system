package page

import "math"

// Metadata a struct for holding the pagination metadata.
type Metadata struct {
	FirstPage   int `json:"firstPage,omitempty"`
	CurrentPage int `json:"currentPage,omitempty"`
	LastPage    int `json:"lastPage,omitempty"`
	RowsPerPage int `json:"rowsPerPage,omitempty"`
	Total       int `json:"total,omitempty"`
}

// calculateMetadata function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that the last page value is calculated using the math.Ceil() function, which rounds
// up a float to the nearest integer. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be math.Ceil(12/5) = 3.
func calculateMetadata(total, page, rows int) Metadata {
	if total == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage: page,
		RowsPerPage: rows,
		FirstPage:   1,
		LastPage:    int(math.Ceil(float64(total) / float64(rows))),
		Total:       total,
	}
}
