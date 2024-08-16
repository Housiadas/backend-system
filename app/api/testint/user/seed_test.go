package user_test

import (
	"context"
	"fmt"

	"github.com/Housiadas/backend-system/app/api/testint"
	"github.com/Housiadas/backend-system/business/domain/authbus"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/dbtest"
)

func insertSeedData(db *dbtest.Database, ath *authbus.Auth) (testint.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, userbus.Roles.Admin, busDomain.User)
	if err != nil {
		return testint.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := testint.User{
		User:  usrs[0],
		Token: testint.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tu2 := testint.User{
		User:  usrs[1],
		Token: testint.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 3, userbus.Roles.User, busDomain.User)
	if err != nil {
		return testint.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := testint.User{
		User:  usrs[0],
		Token: testint.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tu4 := testint.User{
		User:  usrs[1],
		Token: testint.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	tu5 := testint.User{
		User:  usrs[2],
		Token: testint.Token(db.BusDomain.User, ath, usrs[2].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := testint.SeedData{
		Users:  []testint.User{tu3, tu4, tu5},
		Admins: []testint.User{tu1, tu2},
	}

	return sd, nil
}
