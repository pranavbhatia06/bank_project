package sqlc

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/razorpay/bankProject/util"
	"log"
	"os"
	"testing"
)

var queries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	cnf, err := util.LoadConfig("../cnf")
	if err != nil {
		log.Fatalf("unable to Load config: %v", err)
	}

	conn, err := sql.Open(cnf.Database.Dialect, cnf.Database.DataSourceName)
	if err != nil {
		log.Fatalf("unable to open db: %v", err)
	}
	defer conn.Close()
	testDb = conn

	queries = New(conn)

	os.Exit(m.Run())
}
