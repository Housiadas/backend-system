package user_test

import (
	"testing"

	testPck "github.com/Housiadas/backend-system/app/api/testint"
)

func Test_User(t *testing.T) {
	t.Parallel()

	test, err := testPck.StartTest(t, "Test_API_User")
	if err != nil {
		t.Fatalf("Start error: %s", err)
	}

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	//test.Run(t, query200(sd), "query-200")
	//test.Run(t, queryByID200(sd), "query-by-id-200")
	//
	test.Run(t, create200(sd), "create-200")
	//test.Run(t, create401(sd), "create-401")
	//test.Run(t, create400(sd), "create-400")
	//
	////test.Run(t, update200(sd), "update-200")
	//test.Run(t, update401(sd), "update-401")
	//test.Run(t, update400(sd), "update-400")
	//
	//test.Run(t, delete200(sd), "delete-200")
	//test.Run(t, delete401(sd), "delete-401")
}
