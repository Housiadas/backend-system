package audit_repo_test

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/Housiadas/backend-system/internal/common/dbtest"
	"github.com/Housiadas/backend-system/internal/common/unitest"
	"github.com/Housiadas/backend-system/internal/core/domain/audit"
	"github.com/Housiadas/backend-system/internal/core/domain/entity"
	"github.com/Housiadas/backend-system/internal/core/domain/role"
	"github.com/Housiadas/backend-system/internal/core/service/auditcore"
	"github.com/Housiadas/backend-system/internal/core/service/usercore"
	"github.com/Housiadas/backend-system/pkg/order"
	"github.com/Housiadas/backend-system/pkg/page"
)

func Test_Audit(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Audit")

	sd, err := insertSeedData(db.Core)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	unitest.Run(t, query(db.Core, sd), "query")
}

// =============================================================================

func insertSeedData(core dbtest.Core) (unitest.SeedData, error) {
	ctx := context.Background()

	usrs, err := usercore.TestSeedUsers(ctx, 1, role.Admin, core.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	audits, err := auditcore.TestSeedAudits(ctx, 2, usrs[0].ID, entity.User, "create", core.Audit)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := unitest.User{
		User:   usrs[0],
		Audits: audits,
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Admins: []unitest.User{tu1},
	}

	return sd, nil
}

// =============================================================================

func query(core dbtest.Core, sd unitest.SeedData) []unitest.Table {
	sort.Slice(sd.Admins[0].Audits, func(i, j int) bool {
		return sd.Admins[0].Audits[i].ObjName.String() <= sd.Admins[0].Audits[j].ObjName.String()
	})

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: sd.Admins[0].Audits,
			ExcFunc: func(ctx context.Context) any {
				filter := audit.QueryFilter{
					Action: dbtest.StringPointer("create"),
				}

				orderBy := order.NewBy(audit.OrderByObjName, order.ASC)

				resp, err := core.Audit.Query(ctx, filter, orderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]audit.Audit)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]audit.Audit)

				for i := range gotResp {
					if gotResp[i].Timestamp.Format(time.RFC3339) == expResp[i].Timestamp.Format(time.RFC3339) {
						expResp[i].Timestamp = gotResp[i].Timestamp
					}

					gotResp[i].Data = bytes.ReplaceAll(gotResp[i].Data, []byte{' '}, []byte{})
					expResp[i].Data = bytes.ReplaceAll(expResp[i].Data, []byte{' '}, []byte{})
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
