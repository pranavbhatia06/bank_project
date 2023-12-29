package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/razorpay/bankProject/api"
	"github.com/razorpay/bankProject/db/sqlc"
	"log"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:password@192.168.1.7:5432/simple_bank?sslmode=disable"
	address        = "0.0.0.0:8080"
)

func main() {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("error occured while connecting to database %s", err)
	}
	store := sqlc.NewStore(db)
	server := api.NewServer(store)

	server.StartServer(address)

}
