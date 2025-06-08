package user_test

import (
	"context"
	"fmt"

	"github.com/Housiadas/backend-system/app/http/test"
	"github.com/Housiadas/backend-system/internal/domain/userbus"
	"github.com/Housiadas/backend-system/internal/sys/dbtest"
	"github.com/Housiadas/backend-system/internal/sys/types/role"
)

func insertSeedData(db *dbtest.Database) (test.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
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

	usrs, err = userbus.TestSeedUsers(ctx, 3, role.User, busDomain.User)
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
