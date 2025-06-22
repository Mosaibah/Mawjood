package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mawjoodv1 "github.com/mosaibah/Mawjood/gen/go/packages/proto/v1"

	"github.com/mosaibah/Mawjood/packages/discovery/store"
	v1 "github.com/mosaibah/Mawjood/packages/discovery/v1"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "26257")
	dbName := getEnv("DB_NAME", "mawjood")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")
	servicePort := getEnv("SERVICE_PORT", "9002")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	store := store.New(db)
	service := v1.New(store)

	lis, err := net.Listen("tcp", ":"+servicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mawjoodv1.RegisterDiscoveryServiceServer(grpcServer, service)

	reflection.Register(grpcServer)

	log.Printf("Discovery server starting on :%s", servicePort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
