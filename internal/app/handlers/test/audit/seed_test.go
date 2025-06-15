package audit_test

import (
	"context"
	"fmt"

	"github.com/Housiadas/backend-system/internal/common/apitest"
	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/core/domain/entity"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/internal/core/service/authcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
)

func insertSeedData(db *dbtest.Database, ath *authcore.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.Core

	usrs, err := usercore.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	audits, err := auditcore.TestSeedAudits(ctx, 2, usrs[0].ID, entity.User, "create", busDomain.Audit)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := apitest.User{
		User:   usrs[0],
		Audits: audits,
		Token:  apitest.Token(db.Core.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{tu1},
	}

	return sd, nil
}
