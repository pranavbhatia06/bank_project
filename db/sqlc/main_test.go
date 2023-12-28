package sqlc

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:password@192.168.1.10:5432/simple_bank?sslmode=disable"
)

var queries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("unable to open db: %v", err)
	}
	defer conn.Close()
	testDb = conn

	queries = New(conn)

	os.Exit(m.Run())
}
