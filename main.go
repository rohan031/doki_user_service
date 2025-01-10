package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading env file: %v\n", err)
	}

	// neo4j configuration
	ctx := context.Background()
	dbUri := os.Getenv("DB_URI")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	defer driver.Close(ctx)

	if err != nil {
		log.Fatalf("Error connecting to graph: %v", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatalf("Error when verifying connectivity to graph: %v", err)
	}
	fmt.Println("Connection established.")
	fmt.Print("Doki user service")
}