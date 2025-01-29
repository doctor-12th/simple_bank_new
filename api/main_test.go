package api

import (
	// "database/sql"
	// "log"
	"os"
	"testing"
	"time"

	// "github.com/doctor12th/simple_bank_new/util"
	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)
func newTestServer(t *testing.T,store db.Store) *Server{
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server,err := NewServer(config,store)
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestMain(m *testing.M) {
    gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}