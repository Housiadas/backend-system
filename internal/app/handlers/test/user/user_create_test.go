package user_test

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Housiadas/backend-system/internal/app/usecase/user_usecase"
	"github.com/Housiadas/backend-system/internal/common/apitest"
	"github.com/Housiadas/backend-system/pkg/errs"
)

func Test_API_User_Create_200(t *testing.T) {
	t.Parallel()

	test, err := apitest.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/api/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &user_usecase.NewUser{
				Name:            "Chris Housi",
				Email:           "chris@housi.com",
				Roles:           []string{"ADMIN"},
				Department:      "IT0",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &user_usecase.User{},
			ExpResp: &user_usecase.User{
				Name:       "Chris Housi",
				Email:      "chris@housi.com",
				Roles:      []string{"ADMIN"},
				Department: "IT0",
				Enabled:    true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*user_usecase.User)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*user_usecase.User)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	test.Run(t, table, "user-create-200")
}

func Test_API_User_Create_400(t *testing.T) {
	t.Parallel()

	test, err := apitest.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/api/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &user_usecase.NewUser{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validation: [{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"email\",\"error\":\"email is a required field\"},{\"field\":\"roles\",\"error\":\"roles is a required field\"},{\"field\":\"password\",\"error\":\"password is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-role",
			URL:        "/api/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &user_usecase.NewUser{
				Name:            "Chris Housi",
				Email:           "chris@housi.com",
				Roles:           []string{"SUPER"},
				Department:      "IT0",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "parse: invalid role \"SUPER\""),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-name",
			URL:        "/api/v1/users",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &user_usecase.NewUser{
				Name:            "Bi",
				Email:           "chris@housi.com",
				Roles:           []string{"USER"},
				Department:      "IT0",
				Password:        "123",
				PasswordConfirm: "123",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "parse: invalid name \"Bi\""),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	test.Run(t, table, "user-create-400")
}

func Test_API_User_Create_401(t *testing.T) {
	t.Parallel()

	test, err := apitest.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/api/v1/users",
			Token:      "&nbsp;",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badtoken",
			URL:        "/api/v1/users",
			Token:      sd.Admins[0].Token[:10],
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        "/api/v1/users",
			Token:      sd.Admins[0].Token + "A",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        "/api/v1/users",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	test.Run(t, table, "user-create-401")
}
