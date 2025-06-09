package user_test

import (
	"context"
	"fmt"
	"github.com/Housiadas/backend-system/internal/core/service/userservice"

	"github.com/Housiadas/backend-system/internal/adapters/handlers/test"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
)

func insertSeedData(db *dbtest.Database) (test.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userservice.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := test.User{
		User: usrs[0],
	}

	tu2 := test.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	usrs, err = userservice.TestSeedUsers(ctx, 3, role.User, busDomain.User)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := test.User{
		User: usrs[0],
	}

	tu4 := test.User{
		User: usrs[1],
	}

	tu5 := test.User{
		User: usrs[2],
	}

	// -------------------------------------------------------------------------

	sd := test.SeedData{
		Users:  []test.User{tu3, tu4, tu5},
		Admins: []test.User{tu1, tu2},
	}

	return sd, nil
}
