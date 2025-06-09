package user_test

import (
	"fmt"
	"net/http"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Housiadas/backend-system/internal/adapters/domain/userapp"
	testPck "github.com/Housiadas/backend-system/internal/adapters/handlers/test"
	"github.com/Housiadas/backend-system/internal/core/service/userbus"
	"github.com/Housiadas/backend-system/pkg/errs"
	"github.com/Housiadas/backend-system/pkg/page"
)

func Test_API_User_Query_200(t *testing.T) {
	t.Parallel()

	test, err := testPck.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	usrs := make([]userbus.User, 0, len(sd.Admins)+len(sd.Users))

	for _, adm := range sd.Admins {
		usrs = append(usrs, adm.User)
	}

	for _, usr := range sd.Users {
		usrs = append(usrs, usr.User)
	}

	sort.Slice(usrs, func(i, j int) bool {
		return usrs[i].ID.String() <= usrs[j].ID.String()
	})

	table := []testPck.Table{
		{
			Name:       "basic",
			URL:        "/api/v1/users?page=1&rows=10&orderBy=user_id,ASC&name=Name",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &page.Result[userapp.User]{},
			ExpResp: &page.Result[userapp.User]{
				Data: toAppUsers(usrs),
				Metadata: page.Metadata{
					FirstPage:   1,
					CurrentPage: 1,
					LastPage:    1,
					RowsPerPage: 10,
					Total:       len(usrs),
				},
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	test.Run(t, table, "user-query-200")
}

func Test_API_User_Query_BY_ID_200(t *testing.T) {
	t.Parallel()

	test, err := testPck.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	table := []testPck.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/api/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &userapp.User{},
			ExpResp:    toAppUserPtr(sd.Users[0].User),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	test.Run(t, table, "user-query-by-id-200")
}

func Test_API_User_Query_400(t *testing.T) {
	t.Parallel()

	test, err := testPck.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	table := []testPck.Table{
		{
			Name:       "bad-query-filter",
			URL:        "/api/v1/users?page=1&rows=10&email=a.com",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "[{\"field\":\"email\",\"error\":\"mail: missing '@' or angle-addr\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-order-by-value",
			URL:        "/api/v1/users?page=1&rows=10&orderBy=ser_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "[{\"field\":\"order\",\"error\":\"unknown order: ser_id\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	test.Run(t, table, "user-query-400")
}
