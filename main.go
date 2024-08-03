package main

import (
	"database/sql"
	"github.com/Mgeorg1/simpleBank/api"
	"github.com/Mgeorg1/simpleBank/gapi"
	"github.com/Mgeorg1/simpleBank/pb"
	"github.com/Mgeorg1/simpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	db "github.com/Mgeorg1/simpleBank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
	//runGinServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	grpcServer := grpc.NewServer()
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Can not start gRPC server: %s", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
