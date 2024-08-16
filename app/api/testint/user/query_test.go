package user_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"

	"github.com/Housiadas/backend-system/app/api/testint"
	"github.com/Housiadas/backend-system/app/domain/userapp"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/business/sys/page"
)

func query200(sd testint.SeedData) []testint.Table {
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

	table := []testint.Table{
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
					LastPage:    2,
					RowsPerPage: 10,
					Total:       len(usrs),
				},
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByID200(sd testint.SeedData) []testint.Table {
	table := []testint.Table{
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

	return table
}
