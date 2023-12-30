package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/razorpay/bankProject/api"
	"github.com/razorpay/bankProject/db/sqlc"
	"github.com/razorpay/bankProject/util"
	"log"
)

const (
	driverName     = "postgres"
	dataSourceName = "postgresql://root:password@192.168.1.7:5432/simple_bank?sslmode=disable"
	address        = "0.0.0.0:8080"
)

func main() {
	cnf, err := util.LoadConfig("./cnf")
	if err != nil {
		log.Fatalf("error loading config file %s", err)
	}
	fmt.Println(cnf, err)
	db, err := sql.Open(cnf.Database.Dialect, cnf.Database.DataSourceName)
	if err != nil {
		log.Fatalf("error occured while connecting to database %s", err)
	}
	store := sqlc.NewStore(db)
	server := api.NewServer(store)

	server.StartServer(cnf.Application.LISTEN_IP)

}
