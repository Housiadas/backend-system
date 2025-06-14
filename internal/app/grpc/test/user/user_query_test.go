package user_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	userV1 "github.com/Housiadas/backend-system/gen/go/github.com/Housiadas/backend-system/gen/user/v1"
	testPck "github.com/Housiadas/backend-system/internal/app/grpc/test"
)

func Test_GRPC_User_Query_BY_ID(t *testing.T) {
	t.Parallel()

	test, err := testPck.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	table := []testPck.Table{
		{
			Name: "basic",
			Path: "user.v1.UserService/GetUserById",
			Input: &userV1.GetUserByIdRequest{
				Id: sd.Users[0].ID.String(),
			},
			GotResp: &userV1.User{},
			ExpResp: toProtoUser(sd.Users[0].User),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp, protocmp.Transform())
			},
		},
	}

	test.Run(t, table, "grpc-user-query-by-id")
}
