// Package usercore provides internal access to user core.
package usercore

import (
	"context"
	"fmt"

	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/sqldb"
)

// Core manages the set of APIs for user access.
type Core struct {
	log    *logger.Logger
	storer user.Storer
}

// NewCore constructs a user.User internal API for use.
func NewCore(log *logger.Logger, storer user.Storer) *Core {
	return &Core{
		log:    log,
		storer: storer,
	}
}

// NewWithTx constructs a new internal value that will use the
// specified transaction in any store-related calls.
func (c *Core) NewWithTx(tx sqldb.CommitRollbacker) (*Core, error) {
	storer, err := c.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Core{
		log:    c.log,
		storer: storer,
	}

	return &bus, nil
}

// Create adds a new User to the system.
func (c *Core) Create(ctx context.Context, nu user.NewUser) (user.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := user.User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		Department:   nu.Department,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := c.storer.Create(ctx, usr); err != nil {
		return user.User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// Update modifies information about a user.User.
func (c *Core) Update(ctx context.Context, usr user.User, uu user.UpdateUser) (user.User, error) {
	if uu.Name != nil {
		usr.Name = *uu.Name
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return user.User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		usr.PasswordHash = pw
	}

	if uu.Department != nil {
		usr.Department = *uu.Department
	}

	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}
	usr.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, usr); err != nil {
		return user.User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}

// Delete removes the specified user.
func (c *Core) Delete(ctx context.Context, usr user.User) error {
	if err := c.storer.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users.
func (c *Core) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, page page.Page) ([]user.User, error) {
	users, err := c.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// Count returns the total number of users.
func (c *Core) Count(ctx context.Context, filter user.QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

// QueryByID finds the user by the specified ID.
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	usr, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return user.User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return usr, nil
}

// QueryByEmail finds the user by a specified user email.
func (c *Core) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	usr, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		return user.User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return usr, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success, it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c *Core) Authenticate(ctx context.Context, email mail.Address, password string) (user.User, error) {
	usr, err := c.QueryByEmail(ctx, email)
	if err != nil {
		return user.User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return user.User{}, fmt.Errorf("comparehashandpassword: %w", user.ErrAuthenticationFailure)
	}

	return usr, nil
}
