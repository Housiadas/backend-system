package product_test

import (
	"context"
	"fmt"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"

	testPck "github.com/Housiadas/backend-system/internal/app/handlers/test"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
)

func insertSeedData(db *dbtest.Database, ath *authcore.Auth) (testPck.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := usercore.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productcore.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := testPck.User{
		User:     usrs[0],
		Products: prds,

		Token: testPck.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = usercore.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productcore.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
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
