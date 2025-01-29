package main

import (
	"database/sql"
	"log"
	"github.com/doctor12th/simple_bank_new/api"
	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/util"
	_ "github.com/lib/pq"
)

// const(
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:hello58491134@localhost:5432/simple_bank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {
	config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn,err := sql.Open(config.DBDriver, config.DBSource)
	
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store :=db.NewStore(conn)
	server,err := api.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}