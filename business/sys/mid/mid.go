// Package mid contains the set of values the middleware is responsible
// to extract and set.
package mid

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
	productKey
	homeKey
)
