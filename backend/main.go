package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- Models ---

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Do not expose password hash in JSON responses
	CreatedAt    time.Time `json:"created_at"`
}

// TodoItem represents a single to-do item
type TodoItem struct {
	ID          int       `json:"id"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"is_completed"`
	DueDate     string    `json:"due_date"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

// TodoList represents a to-do list, which belongs to a user
type TodoList struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"` // Link to the User
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	Items     []TodoItem `json:"items"`
}

// --- In-Memory Database ---

var users = []User{}
var todoLists = []TodoList{
	{ID: 1, UserID: 1, Title: "Personal", CreatedAt: time.Now(), Items: []TodoItem{
		{ID: 1, Task: "Buy groceries", IsCompleted: false, DueDate: "2025-11-15", Priority: 2, CreatedAt: time.Now()},
		{ID: 2, Task: "Go to the gym", IsCompleted: true, DueDate: "2025-11-10", Priority: 1, CreatedAt: time.Now()},
	}},
	{ID: 2, UserID: 1, Title: "Work", CreatedAt: time.Now(), Items: []TodoItem{
		{ID: 3, Task: "Finish project report", IsCompleted: false, DueDate: "2025-11-20", Priority: 3, CreatedAt: time.Now()},
	}},
}

// Counters for auto-incrementing IDs
var nextUserID = 1
var nextListID = 3
var nextItemID = 4

// --- JWT Configuration ---

// IMPORTANT: In a production application, you must store this secret key securely,
// for example, in an environment variable. Do not hardcode it.
var jwtKey = []byte("my_super_secret_signing_key")

// Claims defines the structure of the JWT payload
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// --- Main Application ---

func main() {
	router := gin.Default()

	// --- Public Routes (No Authentication Required) ---
	router.POST("/api/users/register", RegisterUser)
	router.POST("/api/users/login", LoginUser)

	// --- Protected Routes (Authentication Required) ---
	protected := router.Group("/api")
	protected.Use(AuthMiddleware()) // Apply the authentication middleware to this group
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

// --- Authentication Handlers ---

// RegisterUser handles new user creation
func RegisterUser(c *gin.Context) {
	var payload struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the user's password for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create and "store" the new user
	newUser := User{
		ID:           nextUserID,
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}
	users = append(users, newUser)
	nextUserID++

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": newUser.ID})
}

// LoginUser handles user authentication and issues a JWT
func LoginUser(c *gin.Context) {
	var payload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user User
	// Find the user by email
	for _, u := range users {
		if u.Email == payload.Email {
			user = u
			break
		}
	}
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// If password is correct, generate a JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// --- Authentication Middleware ---

// AuthMiddleware validates the JWT from the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID in the context for subsequent handlers to use
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// --- API Handlers (User-Aware) ---

// getAuthenticatedUserID retrieves the user ID from the context
func getAuthenticatedUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

// getTodoLists returns only the lists belonging to the authenticated user
func getTodoLists(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	var userLists []TodoList
	for _, list := range todoLists {
		if list.UserID == userID {
			userLists = append(userLists, list)
		}
	}
	c.JSON(http.StatusOK, userLists)
}

// createTodoList creates a new list for the authenticated user
func createTodoList(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	var newList TodoList
	if err := c.ShouldBindJSON(&newList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newList.ID = nextListID
	nextListID++
	newList.UserID = userID // Assign to the current user
	newList.CreatedAt = time.Now()
	newList.Items = []TodoItem{}

	todoLists = append(todoLists, newList)
	c.JSON(http.StatusCreated, newList)
}

// getTodoList returns a specific list if it belongs to the authenticated user
func getTodoList(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	listID, _ := strconv.Atoi(c.Param("listId"))

	for _, list := range todoLists {
		if list.ID == listID && list.UserID == userID {
			c.JSON(http.StatusOK, list)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found or access denied"})
}

// updateTodoList updates a list if it belongs to the authenticated user
func updateTodoList(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	listID, _ := strconv.Atoi(c.Param("listId"))

	var updatedList TodoList
	if err := c.ShouldBindJSON(&updatedList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, list := range todoLists {
		if list.ID == listID && list.UserID == userID {
			todoLists[i].Title = updatedList.Title
			c.JSON(http.StatusOK, todoLists[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found or access denied"})
}

// deleteTodoList deletes a list if it belongs to the authenticated user
func deleteTodoList(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	listID, _ := strconv.Atoi(c.Param("listId"))

	for i, list := range todoLists {
		if list.ID == listID && list.UserID == userID {
			todoLists = append(todoLists[:i], todoLists[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "To-do list deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found or access denied"})
}

// createTodoItem adds an item to a list if the list belongs to the authenticated user
func createTodoItem(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	listID, _ := strconv.Atoi(c.Param("listId"))

	var newItem TodoItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, list := range todoLists {
		if list.ID == listID && list.UserID == userID {
			newItem.ID = nextItemID
			nextItemID++
			newItem.CreatedAt = time.Now()
			todoLists[i].Items = append(todoLists[i].Items, newItem)
			c.JSON(http.StatusCreated, newItem)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found or access denied"})
}

// updateTodoItem updates an item if it belongs to the authenticated user
func updateTodoItem(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	itemID, _ := strconv.Atoi(c.Param("itemId"))

	var updatedItem TodoItem
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, list := range todoLists {
		// Ensure the user owns the list this item is in
		if list.UserID != userID {
			continue
		}
		for j, item := range list.Items {
			if item.ID == itemID {
				todoLists[i].Items[j].Task = updatedItem.Task
				todoLists[i].Items[j].IsCompleted = updatedItem.IsCompleted
				todoLists[i].Items[j].DueDate = updatedItem.DueDate
				todoLists[i].Items[j].Priority = updatedItem.Priority
				c.JSON(http.StatusOK, todoLists[i].Items[j])
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "To-do item not found or access denied"})
}

// deleteTodoItem deletes an item if it belongs to the authenticated user
func deleteTodoItem(c *gin.Context) {
	userID, _ := getAuthenticatedUserID(c)
	itemID, _ := strconv.Atoi(c.Param("itemId"))

	for i, list := range todoLists {
		// Ensure the user owns the list this item is in
		if list.UserID != userID {
			continue
		}
		for j, item := range list.Items {
			if item.ID == itemID {
				todoLists[i].Items = append(todoLists[i].Items[:j], todoLists[i].Items[j+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "To-do item deleted successfully"})
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "To-do item not found or access denied"})
}
