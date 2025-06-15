package product_test

import (
	"context"
	"fmt"
	"github.com/Housiadas/backend-system/internal/common/apitest"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"

	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/productcore"
)

func insertSeedData(db *dbtest.Database, ath *authcore.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.Core

	usrs, err := usercore.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := productcore.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu1 := apitest.User{
		User:     usrs[0],
		Products: prds,

		Token: apitest.Token(db.Core.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = usercore.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = productcore.TestGenerateSeedProducts(ctx, 2, busDomain.Product, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding products : %w", err)
	}

	tu2 := apitest.User{
		User:     usrs[0],
		Products: prds,
		Token:    apitest.Token(db.Core.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{tu2},
		Users:  []apitest.User{tu1},
	}

	return sd, nil
}
