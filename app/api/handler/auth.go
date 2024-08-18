package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/business/web"
)

func (h *Handler) authenticate(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var requestData userapp.AuthenticateUser
	if err := web.Decode(r, &requestData); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	usr, err := h.App.User.Authenticate(ctx, requestData)
	if err != nil {
		return errs.New(errs.InvalidArgument, errors.New("invalid credentials"))
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in 8 hours.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := authbus.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID,
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	// This will generate a JWT with the claims embedded in them. The database
	// with need to be configured with the information found in the public key
	// file to validation these claims. Dgraph does not support key rotate at this time.
	token, err := h.Business.Auth.GenerateToken(claims)
	if err != nil {
		return errs.New(errs.Internal, fmt.Errorf("generating token: %w", err))
	}

	return authbus.AuthenticateResp{
		Token: token,
	}
}

func (h *Handler) authorize(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {
	var authData authbus.Authorize
	if err := web.Decode(r, &authData); err != nil {
		return errs.New(errs.FailedPrecondition, err)
	}

	if err := h.Business.Auth.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated,
			"authorize: you are not authorized for that action, claims[%v] rule[%v]: %s",
			authData.Claims.Roles,
			authData.Rule, err,
		)
	}

	return nil
}

func (h *Handler) token(ctx context.Context, _ http.ResponseWriter, r *http.Request) web.Encoder {

	token, err := h.App.User.Token(ctx)
	if err != nil {
		return errs.NewError(err)
	}

	return token
}
