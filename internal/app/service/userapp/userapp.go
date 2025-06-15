// Package userapp maintains the cli layer api for the user core.
package userapp

import (
	"context"
	"errors"
	"net/mail"

	ctxPck "github.com/Housiadas/backend-system/internal/common/context"
	"github.com/Housiadas/backend-system/internal/common/validation"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
)

// App manages the set of cli layer api functions for the user core.
type App struct {
	authCore *authcore.Auth
	userCore *usercore.Core
}

// NewApp constructs a user cli API for use.
func NewApp(userBus *usercore.Core) *App {
	return &App{
		userCore: userBus,
	}
}

// NewAppWithAuth constructs a user cli API for use.
func NewAppWithAuth(userBus *usercore.Core, authbus *authcore.Auth) *App {
	return &App{
		authCore: authbus,
		userCore: userBus,
	}
}

// Create adds a new user to the system.
func (a *App) Create(ctx context.Context, app NewUser) (User, error) {
	nc, err := toBusNewUser(app)
	if err != nil {
		return User{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userCore.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return User{}, errs.New(errs.Aborted, user.ErrUniqueEmail)
		}
		return User{}, errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr), nil
}

// Update updates an existing user.
func (a *App) Update(ctx context.Context, app UpdateUser) (User, error) {
	uu, err := toBusUpdateUser(app)
	if err != nil {
		return User{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := ctxPck.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := a.userCore.Update(ctx, usr, uu)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "update: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// UpdateRole updates an existing user's role.
func (a *App) UpdateRole(ctx context.Context, app UpdateUserRole) (User, error) {
	uu, err := toBusUpdateUserRole(app)
	if err != nil {
		return User{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := ctxPck.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := a.userCore.Update(ctx, usr, uu)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "updaterole: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// Delete removes a user from the system.
func (a *App) Delete(ctx context.Context) error {
	usr, err := ctxPck.GetUser(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "userID missing in context: %s", err)
	}

	if err := a.userCore.Delete(ctx, usr); err != nil {
		return errs.Newf(errs.Internal, "delete: userID[%s]: %s", usr.ID, err)
	}

	return nil
}

// Query returns a list of users with paging.
func (a *App) Query(ctx context.Context, qp AppQueryParams) (page.Result[User], error) {
	p, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return page.Result[User]{}, validation.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Result[User]{}, err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return page.Result[User]{}, validation.NewFieldErrors("order", err)
	}

	usrs, err := a.userCore.Query(ctx, filter, orderBy, p)
	if err != nil {
		return page.Result[User]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.userCore.Count(ctx, filter)
	if err != nil {
		return page.Result[User]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return page.NewResult(toAppUsers(usrs), total, p), nil
}

// QueryByID returns a user by its Ia.
func (a *App) QueryByID(ctx context.Context) (User, error) {
	usr, err := ctxPck.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppUser(usr), nil
}

// Authenticate provides an API to authenticate the user.
func (a *App) Authenticate(ctx context.Context, authUser AuthenticateUser) (User, error) {
	addr, err := mail.ParseAddress(authUser.Email)
	if err != nil {
		return User{}, validation.NewFieldErrors("email", err)
	}

	usr, err := a.userCore.Authenticate(ctx, *addr, authUser.Password)
	if err != nil {
		return User{}, err
	}

	return toAppUser(usr), nil
}

// Token provides an API token for the authenticated user.
func (a *App) Token(ctx context.Context) (Token, error) {
	if a.authCore == nil {
		return Token{}, errs.Newf(errs.Internal, "authapi not configured")
	}

	claims := ctxPck.GetClaims(ctx)

	tkn, err := a.authCore.GenerateToken(claims)
	if err != nil {
		return Token{}, errs.New(errs.Internal, err)
	}

	return toToken(tkn), nil
}
