package api

import (
	"os"
	"testing"
	"time"

	"github.com/Housiadas/simple-banking-system/business/db"
	"github.com/Housiadas/simple-banking-system/database/config"
	"github.com/Housiadas/simple-banking-system/foundation/random"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	cfg := config.Config{
		TokenSymmetricKey:   random.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(cfg, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
