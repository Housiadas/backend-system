// Package test provides support for integration http tests.
package test

import (
	"log"
	"testing"

	"google.golang.org/grpc/test/bufconn"

	"github.com/Housiadas/backend-system/app/grpc/server"
	"github.com/Housiadas/backend-system/internal/sys/dbtest"
)

// Test contains functions for executing an api test.
type Test struct {
	DB     *dbtest.Database
	Server *server.Server
}

// New constructs a Test value for running api tests.
func New(db *dbtest.Database, s *server.Server) *Test {
	return &Test{
		DB:     db,
		Server: s,
	}
}

// Run performs the actual test logic based on the table data.
func (at *Test) Run(t *testing.T, table []Table, testName string) {
	for _, tt := range table {
		f := func(t *testing.T) {

			buffer := 101024 * 1024
			lis := bufconn.Listen(buffer)
			grpcServer := at.Server.Registrar()
			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					log.Printf("error serving server: %v", err)
				}
			}()

			diff := tt.CmpFunc(tt.GotResp, tt.ExpResp)
			if diff != "" {
				t.Log("DIFF")
				t.Logf("%s", diff)
				t.Log("GOT")
				t.Logf("%#v", tt.GotResp)
				t.Log("EXP")
				t.Logf("%#v", tt.ExpResp)
				t.Fatalf("Should get the expected response")
			}

			err := lis.Close()
			if err != nil {
				log.Printf("error closing listener: %v", err)
			}
			grpcServer.Stop()
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}
