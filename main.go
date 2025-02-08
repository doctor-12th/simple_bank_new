package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/doctor12th/simple_bank_new/api"
	db "github.com/doctor12th/simple_bank_new/db/sqlc"
	"github.com/doctor12th/simple_bank_new/gapi"
	"github.com/doctor12th/simple_bank_new/pb"
	"github.com/doctor12th/simple_bank_new/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGrpcServer(config,store)

	
}

func runGinServer(config util.Config,store db.Store) {
	server,err := api.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runGrpcServer(config util.Config,store db.Store) {
	server,err:=gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankNewServer(grpcServer,server)
	reflection.Register(grpcServer)
	listener,err :=  net.Listen("tcp",config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Printf("start grpc at %s",listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}