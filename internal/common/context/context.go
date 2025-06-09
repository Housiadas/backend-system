package context

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/internal/core/service/authbus"
	"github.com/Housiadas/backend-system/internal/core/service/productservice"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"
)

type ctxKey string

const (
	requestID  ctxKey = "requestID"
	apiVersion ctxKey = "apiVersion"
	claimKey   ctxKey = "claimKey"
	userIDKey  ctxKey = "userIDKey"
	userKey    ctxKey = "userKey"
	productKey ctxKey = "productKey"
)

func SetRequestID(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, requestID, reqId)
}

func GetRequestID(ctx context.Context) string {
	v, ok := ctx.Value(requestID).(string)
	if !ok {
		return ""
	}
	return v
}

func SetApiVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, apiVersion, version)
}

func GetApiVersion(ctx context.Context) string {
	v, ok := ctx.Value(apiVersion).(string)
	if !ok {
		return ""
	}
	return v
}

func SetClaims(ctx context.Context, claims authbus.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) authbus.Claims {
	v, ok := ctx.Value(claimKey).(authbus.Claims)
	if !ok {
		return authbus.Claims{}
	}
	return v
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}

func SetUser(ctx context.Context, usr userservice.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (userservice.User, error) {
	v, ok := ctx.Value(userKey).(userservice.User)
	if !ok {
		return userservice.User{}, errors.New("user not found in context")
	}

	return v, nil
}

func SetProduct(ctx context.Context, prd productservice.Product) context.Context {
	return context.WithValue(ctx, productKey, prd)
}

// GetProduct returns the productapi from the context.
func GetProduct(ctx context.Context) (productservice.Product, error) {
	v, ok := ctx.Value(productKey).(productservice.Product)
	if !ok {
		return productservice.Product{}, errors.New("productapi not found in context")
	}

	return v, nil
}
