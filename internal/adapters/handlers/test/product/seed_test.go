package product_test

import (
	"context"
	"fmt"

	testPck "github.com/Housiadas/backend-system/internal/adapters/handlers/test"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/authbus"
	"github.com/Housiadas/backend-system/internal/core/service/productbus"
	"github.com/Housiadas/backend-system/internal/core/service/userbus"
)

func insertSeedData(db *dbtest.Database, ath *authbus.Auth) (testPck.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := testPck.User{
		User:     usrs[0],
		Products: prds,

		Token: testPck.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu2 := testPck.User{
		User:     usrs[0],
		Products: prds,
		Token:    testPck.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := testPck.SeedData{
		Admins: []testPck.User{tu2},
		Users:  []testPck.User{tu1},
	}

	return sd, nil
}
