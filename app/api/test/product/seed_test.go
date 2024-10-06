package product_test

import (
	"context"
	"fmt"

	"github.com/Housiadas/backend-system/app/api/test"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/productbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/dbtest"
	"github.com/Housiadas/backend-system/business/sys/types/role"
)

func insertSeedData(db *dbtest.Database, ath *authbus.Auth) (test.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := test.User{
		User:     usrs[0],
		Products: prds,

		Token: test.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productbus.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu2 := test.User{
		User:     usrs[0],
		Products: prds,
		Token:    test.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := test.SeedData{
		Admins: []test.User{tu2},
		Users:  []test.User{tu1},
	}

	return sd, nil
}
