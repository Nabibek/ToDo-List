package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	httpadapter "ToDo-List/internal/adapters/http"
	"ToDo-List/internal/adapters/logger"
	"ToDo-List/internal/repo"

	_ "github.com/lib/pq"
)

func main() {
	// Инициализация логгера
	loggerConfig := logger.Config{
		Level: logger.ParseLevel(os.Getenv("LOG_LEVEL")),
	}
	appLogger := logger.New(loggerConfig)

	appLogger.Info("Starting ToDo application...")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		appLogger.Fatal("DATABASE_URL environment variable is not set")
	}

	db := waitForDatabase(databaseURL, appLogger)
	defer db.Close()

	appLogger.Info("Successfully connected to the database!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		appLogger.Info("Using default port: %s", port)
	}

	todoRepo := repo.NewPostgreRepo(db, appLogger)
	router := httpadapter.NewRouter(todoRepo, appLogger)

	appLogger.Info("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		appLogger.Fatal("Server failed: %v", err)
	}
}

func waitForDatabase(databaseURL string, appLogger *logger.Logger) *sql.DB {
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			appLogger.Warn("Failed to connect to DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			appLogger.Warn("Failed to ping DB (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		appLogger.Info("Successfully connected to the database!")
		return db
	}
	appLogger.Fatal("Could not connect to the database after multiple attempts: %v", err)
	return nil
}
