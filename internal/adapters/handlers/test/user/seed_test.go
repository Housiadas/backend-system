package user_test

import (
	"context"
	"fmt"
	testPck "github.com/Housiadas/backend-system/internal/adapters/handlers/test"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/authbus"
	"github.com/Housiadas/backend-system/internal/core/service/userbus"
)

func insertSeedData(db *dbtest.Database, ath *authbus.Auth) (testPck.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := testPck.User{
		User:  usrs[0],
		Token: testPck.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tu2 := testPck.User{
		User:  usrs[1],
		Token: testPck.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 3, role.User, busDomain.User)
	if err != nil {
		return testPck.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := testPck.User{
		User:  usrs[0],
		Token: testPck.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tu4 := testPck.User{
		User:  usrs[1],
		Token: testPck.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	tu5 := testPck.User{
		User:  usrs[2],
		Token: testPck.Token(db.BusDomain.User, ath, usrs[2].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := testPck.SeedData{
		Users:  []testPck.User{tu3, tu4, tu5},
		Admins: []testPck.User{tu1, tu2},
	}

	return sd, nil
}
