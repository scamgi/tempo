package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"tempo-backend/api"
	"tempo-backend/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	fmt.Println("Successfully connected to the database!")

	// Initialize stores and handlers
	userStore := db.NewUserStore(dbpool)
	userHandler := api.NewUserHandler(userStore)

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- API Routes ---
	apiGroup := e.Group("/api")

	// User routes
	userGroup := apiGroup.Group("/users")
	userGroup.POST("/register", userHandler.HandleRegisterUser)
	userGroup.POST("/login", userHandler.HandleLoginUser)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Fatal(e.Start(":" + port))
}
