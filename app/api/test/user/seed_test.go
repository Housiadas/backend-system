package user_test

import (
	"context"
	"fmt"

	"github.com/Housiadas/backend-system/app/api/test"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/dbtest"
	"github.com/Housiadas/backend-system/business/sys/types/role"
)

func insertSeedData(db *dbtest.Database, ath *authbus.Auth) (test.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := test.User{
		User:  usrs[0],
		Token: test.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tu2 := test.User{
		User:  usrs[1],
		Token: test.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 3, role.User, busDomain.User)
	if err != nil {
		return test.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := test.User{
		User:  usrs[0],
		Token: test.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tu4 := test.User{
		User:  usrs[1],
		Token: test.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	tu5 := test.User{
		User:  usrs[2],
		Token: test.Token(db.BusDomain.User, ath, usrs[2].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := test.SeedData{
		Users:  []test.User{tu3, tu4, tu5},
		Admins: []test.User{tu1, tu2},
	}

	return sd, nil
}
