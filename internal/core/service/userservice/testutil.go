package userservice

import (
	"context"
	"fmt"
	"math/rand"
	"net/mail"

	"github.com/Housiadas/backend-system/internal/core/domain/name"
	rolePck "github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/domain/user"
)

// TestNewUsers is a helper method for testing.
func TestNewUsers(n int, role rolePck.Role) []user.NewUser {
	newUsrs := make([]user.NewUser, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nu := user.NewUser{
			Name:       name.MustParse(fmt.Sprintf("Name%d", idx)),
			Email:      mail.Address{Address: fmt.Sprintf("Email%d@gmail.com", idx)},
			Roles:      []rolePck.Role{role},
			Department: name.MustParseNull(fmt.Sprintf("Department%d", idx)),
			Password:   fmt.Sprintf("Password%d", idx),
		}

		newUsrs[i] = nu
	}

	return newUsrs
}

// TestSeedUsers is a helper method for testing.
func TestSeedUsers(ctx context.Context, n int, role rolePck.Role, api *Service) ([]user.User, error) {
	newUsrs := TestNewUsers(n, role)

	usrs := make([]user.User, len(newUsrs))
	for i, nu := range newUsrs {
		usr, err := api.Create(ctx, nu)
		if err != nil {
			return nil, fmt.Errorf("seeding user: idx: %d : %w", i, err)
		}

		usrs[i] = usr
	}

	return usrs, nil
}
