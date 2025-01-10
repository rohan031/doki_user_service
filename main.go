package main

import (
	"context"
	"doki.co.in/doki_user_service/permissions"
	"doki.co.in/doki_user_service/user"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading env file: %v\n", err)
	}

	// get db driver
	driver, ctx := connectDB()
	defer (*driver).Close(ctx)

	// listen for tcp connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error creating tcp listener: %v", err)
	}

	serverRegistrar := grpc.NewServer()

	// adding db driver my dependency injection
	userService := &user.UserServerImpl{Driver: driver}
	permissionService := &permissions.PermissionServerImpl{Driver: driver}

	user.RegisterUserServer(serverRegistrar, userService)
	permissions.RegisterPermissionServer(serverRegistrar, permissionService)

	fmt.Println("Doki user service")
	err = serverRegistrar.Serve(listener)
	if err != nil {
		log.Fatalf("Error starting grpc server: %v\n", err)
	}
}

func connectDB() (*neo4j.DriverWithContext, context.Context) {
	// neo4j configuration
	ctx := context.Background()
	dbUri := os.Getenv("DB_URI")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))

	if err != nil {
		log.Fatalf("Error connecting to graph: %v\n", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatalf("Error when verifying connectivity to graph: %v\n", err)
	}
	fmt.Println("Connection established.")
	return &driver, ctx
}