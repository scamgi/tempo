package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// --- Public Routes (No Authentication Required) ---
	router.POST("/api/users/register", RegisterUser)
	router.POST("/api/users/login", LoginUser)

	// --- Protected Routes (Authentication Required) ---
	protected := router.Group("/api")
	protected.Use(AuthMiddleware()) // Apply the authentication middleware
	{
		// To-Do Lists
		protected.GET("/lists", getTodoLists)
		protected.POST("/lists", createTodoList)
		protected.GET("/lists/:listId", getTodoList)
		protected.PUT("/lists/:listId", updateTodoList)
		protected.DELETE("/lists/:listId", deleteTodoList)

		// To-Do Items
		protected.POST("/lists/:listId/items", createTodoItem)
		protected.PUT("/items/:itemId", updateTodoItem)
		protected.DELETE("/items/:itemId", deleteTodoItem)
	}

	router.Run(":8080")
}
