package main

import (
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mawjoodv1 "mawjood/gen/go/packages/proto/v1"

	"github.com/mosaibah/Mawjood/packages/discovery/store"
	v1 "github.com/mosaibah/Mawjood/packages/discovery/v1"
)

func main() {
	const connStr = "postgres://root:@localhost:26257/mawjood?sslmode=disable&parseTime=true"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	store := store.New(db)
	service := v1.New(store)

	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mawjoodv1.RegisterDiscoveryServiceServer(grpcServer, service)

	// Enable reflection for grpcui
	reflection.Register(grpcServer)

	log.Printf("Discovery server starting on :9002")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
