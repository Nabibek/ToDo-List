package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	httpadapter "ToDo-List/internal/adapters/http"
	"ToDo-List/internal/adapters/repo"

	_ "github.com/lib/pq"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	todoRepo := repo.NewPostgreRepo(db)

	router := httpadapter.NewRouter(todoRepo)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	fmt.Printf("Starting server at :%s...\n", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
