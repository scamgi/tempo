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

	todoStore := db.NewTodoStore(dbpool)
	todoHandler := api.NewTodoHandler(todoStore)

	noteStore := db.NewNoteStore(dbpool)
	noteHandler := api.NewNoteHandler(noteStore)

	journalStore := db.NewJournalStore(dbpool)
	journalHandler := api.NewJournalHandler(journalStore)

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- API Routes ---
	apiGroup := e.Group("/api")

	// User routes (public)
	userGroup := apiGroup.Group("/users")
	userGroup.POST("/register", userHandler.HandleRegisterUser)
	userGroup.POST("/login", userHandler.HandleLoginUser)

	// To-Do List routes (protected)
	listGroup := apiGroup.Group("/lists")
	listGroup.Use(api.JWTAuthMiddleware) // Apply the middleware to all routes in this group
	listGroup.POST("", todoHandler.HandleCreateTodoList)
	listGroup.GET("", todoHandler.HandleGetTodoLists)
	listGroup.GET("/:listId", todoHandler.HandleGetTodoListAndItems)
	listGroup.DELETE("/:listId", todoHandler.HandleDeleteTodoList)
	listGroup.POST("/:listId/items", todoHandler.HandleCreateTodoItem)

	// To-Do Item routes (protected)
	itemGroup := apiGroup.Group("/items")
	itemGroup.Use(api.JWTAuthMiddleware)
	itemGroup.PUT("/:itemId", todoHandler.HandleUpdateTodoItem)
	itemGroup.DELETE("/:itemId", todoHandler.HandleDeleteTodoItem)

	// Notes routes (protected)
	noteGroup := apiGroup.Group("/notes")
	noteGroup.Use(api.JWTAuthMiddleware)
	noteGroup.POST("", noteHandler.HandleCreateNote)
	noteGroup.GET("", noteHandler.HandleGetNotes)
	noteGroup.GET("/:noteId", noteHandler.HandleGetNote)
	noteGroup.PUT("/:noteId", noteHandler.HandleUpdateNote)
	noteGroup.DELETE("/:noteId", noteHandler.HandleDeleteNote)

	// Journal routes (protected)
	journalGroup := apiGroup.Group("/journal")
	journalGroup.Use(api.JWTAuthMiddleware)
	journalGroup.POST("", journalHandler.HandleCreateJournalEntry)
	journalGroup.GET("", journalHandler.HandleGetJournalEntries)
	journalGroup.GET("/:entryId", journalHandler.HandleGetJournalEntry)
	journalGroup.PUT("/:entryId", journalHandler.HandleUpdateJournalEntry)
	journalGroup.DELETE("/:entryId", journalHandler.HandleDeleteJournalEntry)


	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}