// Package user_cache contains user-related functionality with caching.
package user_cache

import (
	"context"

	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/viccon/sturdyc"

	"github.com/Housiadas/backend-system/internal/core/domain/user"
	"github.com/Housiadas/backend-system/pkg/logger"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
	"github.com/Housiadas/backend-system/pkg/pgsql"
)

// Store manages the set of APIs for user data and caching.
type Store struct {
	storer user.Storer
	log    *logger.Logger
	cache  *sturdyc.Client[user.User]
}

// NewStore constructs the api for data and caching access.
func NewStore(log *logger.Logger, storer user.Storer, ttl time.Duration) *Store {
	const capacity = 10000
	const numShards = 10
	const evictionPercentage = 10

	return &Store{
		log:    log,
		storer: storer,
		cache:  sturdyc.New[user.User](capacity, numShards, ttl, evictionPercentage),
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx pgsql.CommitRollbacker) (user.Storer, error) {
	return s.storer.NewWithTx(tx)
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {
	if err := s.storer.Create(ctx, usr); err != nil {
		return err
	}

	s.writeCache(usr)

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr user.User) error {
	if err := s.storer.Update(ctx, usr); err != nil {
		return err
	}

	s.writeCache(usr)

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	if err := s.storer.Delete(ctx, usr); err != nil {
		return err
	}

	s.deleteCache(usr)

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(
	ctx context.Context,
	filter user.QueryFilter,
	orderBy order.By,
	page page.Page,
) ([]user.User, error) {
	return s.storer.Query(ctx, filter, orderBy, page)
}

// Count returns the total number of cards in the DB.
func (s *Store) Count(ctx context.Context, filter user.QueryFilter) (int, error) {
	return s.storer.Count(ctx, filter)
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	cachedUsr, ok := s.readCache(userID.String())
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.storer.QueryByID(ctx, userID)
	if err != nil {
		return user.User{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	cachedUsr, ok := s.readCache(email.Address)
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.storer.QueryByEmail(ctx, email)
	if err != nil {
		return user.User{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

// readCache performs a safe search in the cache for the specified key.
func (s *Store) readCache(key string) (user.User, bool) {
	usr, exists := s.cache.Get(key)
	if !exists {
		return user.User{}, false
	}

	return usr, true
}

// writeCache performs a safe writing to the cache for the specified user.
func (s *Store) writeCache(bus user.User) {
	s.cache.Set(bus.ID.String(), bus)
	s.cache.Set(bus.Email.Address, bus)
}

// deleteCache performs a safe removal from the cache for the specified user.
func (s *Store) deleteCache(bus user.User) {
	s.cache.Delete(bus.ID.String())
	s.cache.Delete(bus.Email.Address)
}
