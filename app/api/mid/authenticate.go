package mid

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/mid"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/foundation/web"
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(auth *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authorization := r.Header.Get("authorization")
			parts := strings.Split(authorization, " ")

			var err error

			switch parts[0] {
			case "Bearer":
				ctx, err = processJWT(ctx, auth, authorization)
			}

			if err != nil {
				return err
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

func processJWT(ctx context.Context, auth *auth.Auth, token string) (context.Context, error) {
	claims, err := auth.Authenticate(ctx, token)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return ctx, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	ctx = mid.SetUserID(ctx, subjectID)
	ctx = mid.SetClaims(ctx, claims)

	return ctx, nil
}
