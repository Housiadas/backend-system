// Package authapi maintains the web based api for auth access.
package authapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/auth"
	"github.com/Housiadas/backend-system/business/sys/errs"
	"github.com/Housiadas/backend-system/foundation/validate"
	"github.com/Housiadas/backend-system/foundation/web"
)

type api struct {
	auth    *auth.Auth
	userApp *userapp.App
}

func newAPI(userApp *userapp.App, auth *auth.Auth) *api {
	return &api{
		userApp: userApp,
		auth:    auth,
	}
}

func (api *api) authenticate(ctx context.Context, _ http.ResponseWriter, r *http.Request) (any, error) {
	kid := web.Header(r, "kid")
	if kid == "" {
		return nil, errs.New(errs.FailedPrecondition, validate.NewFieldsError("kid", errors.New("missing kid")))
	}

	var requestData userapp.AuthenticateUser
	if err := web.Decode(r, &requestData); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	usr, err := api.userApp.Authenticate(ctx, requestData)
	if err != nil {
		return nil, errs.New(errs.InvalidArgument, invalidCredentials)
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
	claims := auth.Claims{
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
	// file to validate these claims. Dgraph does not support key rotate at this time.
	token, err := api.auth.GenerateToken(kid, claims)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	return data, nil
}

func (api *api) authorize(ctx context.Context, _ http.ResponseWriter, r *http.Request) (any, error) {
	var authData auth.Authorize
	if err := web.Decode(r, &authData); err != nil {
		return nil, errs.New(errs.FailedPrecondition, err)
	}

	if err := api.auth.Authorize(ctx, authData.Claims, authData.UserID, authData.Rule); err != nil {
		return nil, errs.Newf(errs.Unauthenticated,
			"authorize: you are not authorized for that action, claims[%v] rule[%v]: %s",
			authData.Claims.Roles,
			authData.Rule, err,
		)
	}

	return nil, nil
}

func (api *api) token(ctx context.Context, _ http.ResponseWriter, r *http.Request) (any, error) {
	kid := web.Param(r, "kid")
	if kid == "" {
		return nil, errs.New(errs.FailedPrecondition, validate.NewFieldsError("kid", errors.New("missing kid")))
	}

	token, err := api.userApp.Token(ctx, kid)
	if err != nil {
		return nil, err
	}

	return token, nil
}
