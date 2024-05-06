// Package userbus Package users provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	conflkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Housiadas/backend-system/business/data/transaction"
	"github.com/Housiadas/backend-system/business/sys/order"
	"github.com/Housiadas/backend-system/foundation/kafka"
	"github.com/Housiadas/backend-system/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to persists and retrieve data.
type Storer interface {
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// Business manages the set of APIs for user access.
type Business struct {
	log      *logger.Logger
	storer   Storer
	producer *kafka.ProducerClient
}

// NewBusiness constructs a user core API for use.
func NewBusiness(log *logger.Logger, storer Storer, p *kafka.ProducerClient) *Business {
	return &Business{
		log:      log,
		storer:   storer,
		producer: p,
	}
}

// ExecuteUnderTransaction constructs a new Business value that will use the
// specified transaction in any store related calls.
func (c *Business) ExecuteUnderTransaction(tx transaction.Transaction) (*Business, error) {
	trS, err := c.storer.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	core := Business{
		log:      c.log,
		storer:   trS,
		producer: c.producer,
	}

	return &core, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success, it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	usr, err := c.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}

// Create adds a new user to the system.
func (c *Business) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
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
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// Update modifies information about a user.
func (c *Business) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
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
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
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
		return User{}, fmt.Errorf("update: %w", err)
	}

	event := ActionUpdatedData(uu, usr.ID)
	err := c.producer.Produce(ctx, &conflkafka.Message{
		TopicPartition: conflkafka.TopicPartition{
			Topic:     &event.Topic,
			Partition: conflkafka.PartitionAny,
		},
		Value: event.Data,
	})
	if err != nil {
		return User{}, fmt.Errorf("failed to produce `%s` error: %w", UserUpdatedEvent, err)
	}

	return usr, nil
}

// Delete removes the specified user.
func (c *Business) Delete(ctx context.Context, usr User) error {
	if err := c.storer.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users.
func (c *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	users, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// Count returns the total number of users.
func (c *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	if err := filter.Validate(); err != nil {
		return 0, err
	}

	return c.storer.Count(ctx, filter)
}

// QueryByID finds the user by the specified ID.
func (c *Business) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}

// QueryByIDs finds the users by a specified User IDs.
func (c *Business) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error) {
	user, err := c.storer.QueryByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("query: userIDs[%s]: %w", userIDs, err)
	}

	return user, nil
}

// QueryByEmail finds the user by a specified user email.
func (c *Business) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}
