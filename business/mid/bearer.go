package mid

import (
	"context"
	"github.com/google/uuid"

	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/sys/errs"
)

// Bearer processes JWT token.
func Bearer(ctx context.Context, ath *auth.Auth, authorization string, next Handler) (any, error) {
	claims, err := ath.Authenticate(ctx, authorization)
	if err != nil {
		return nil, errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return nil, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return next(ctx)
}
