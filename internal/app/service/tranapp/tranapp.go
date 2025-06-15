// Package tranapp maintains the cli layer http for the tran core.
package tranapp

import (
	"context"
	"errors"

	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/sqldb"
)

// App manages the set of cli layer http functions for the tran core.
type App struct {
	userBus    *usercore.Core
	productBus *productcore.Core
}

// NewApp constructs a tran cli API for use.
func NewApp(userBus *usercore.Core, productBus *productcore.Core) *App {
	return &App{
		userBus:    userBus,
		productBus: productBus,
	}
}

// newWithTx constructs a new Handlers value with the core apis
// using a store transaction that was created via middleware.
func (a *App) newWithTx(ctx context.Context) (*App, error) {
	tx, err := sqldb.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	userBus, err := a.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	productBus, err := a.productBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := App{
		userBus:    userBus,
		productBus: productBus,
	}

	return &app, nil
}

// Create adds a new user and product at the same time under a single transaction.
func (a *App) Create(ctx context.Context, nt NewTran) (Product, error) {
	a, err := a.newWithTx(ctx)
	if err != nil {
		return Product{}, errs.New(errs.Internal, err)
	}

	np, err := toBusNewProduct(nt.Product)
	if err != nil {
		return Product{}, errs.New(errs.InvalidArgument, err)
	}

	nu, err := toBusNewUser(nt.User)
	if err != nil {
		return Product{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return Product{}, errs.New(errs.Aborted, user.ErrUniqueEmail)
		}
		return Product{}, errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	np.UserID = usr.ID

	prd, err := a.productBus.Create(ctx, np)
	if err != nil {
		return Product{}, errs.Newf(errs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}
