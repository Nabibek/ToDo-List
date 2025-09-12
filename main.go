package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	httpadapter "ToDo-List/internal/adapters/http"
	"ToDo-List/internal/adapters/repo"

	_ "github.com/lib/pq"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db := waitForDatabase(databaseURL)
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	todoRepo := repo.NewPostgreRepo(db)
	router := httpadapter.NewRouter(todoRepo)

	fmt.Printf("Starting server at :%s...\n", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func waitForDatabase(databaseURL string) *sql.DB {
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Printf("Failed to connect to DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Successfully connected to the database!")
		return db
	}
	log.Fatalf("Could not connect to the database after multiple attempts: %v", err)
	return nil
}
