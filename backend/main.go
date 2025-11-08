package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Represents a to-do item
type TodoItem struct {
	ID          int       `json:"id"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"is_completed"`
	DueDate     string    `json:"due_date"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

// Represents a to-do list
type TodoList struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	Items     []TodoItem `json:"items"`
}

// In-memory "database"
var todoLists = []TodoList{
	{
		ID:        1,
		Title:     "Personal",
		CreatedAt: time.Now(),
		Items: []TodoItem{
			{ID: 1, Task: "Buy groceries", IsCompleted: false, DueDate: "2025-11-15", Priority: 2, CreatedAt: time.Now()},
			{ID: 2, Task: "Go to the gym", IsCompleted: true, DueDate: "2025-11-10", Priority: 1, CreatedAt: time.Now()},
		},
	},
	{
		ID:        2,
		Title:     "Work",
		CreatedAt: time.Now(),
		Items: []TodoItem{
			{ID: 3, Task: "Finish project report", IsCompleted: false, DueDate: "2025-11-20", Priority: 3, CreatedAt: time.Now()},
		},
	},
}

var nextListID = 3
var nextItemID = 4

func main() {
	router := gin.Default()

	// To-Do Lists endpoints
	router.GET("/api/lists", getTodoLists)
	router.POST("/api/lists", createTodoList)
	router.GET("/api/lists/:listId", getTodoList)
	router.PUT("/api/lists/:listId", updateTodoList)
	router.DELETE("/api/lists/:listId", deleteTodoList)

	// To-Do Items endpoints
	router.POST("/api/lists/:listId/items", createTodoItem)
	router.PUT("/api/items/:itemId", updateTodoItem)
	router.DELETE("/api/items/:itemId", deleteTodoItem)

	router.Run(":8080")
}

// Handler to get all to-do lists
func getTodoLists(c *gin.Context) {
	c.JSON(http.StatusOK, todoLists)
}

// Handler to create a new to-do list
func createTodoList(c *gin.Context) {
	var newList TodoList
	if err := c.ShouldBindJSON(&newList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newList.ID = nextListID
	nextListID++
	newList.CreatedAt = time.Now()
	newList.Items = []TodoItem{} // Initialize with an empty slice of items

	todoLists = append(todoLists, newList)
	c.JSON(http.StatusCreated, newList)
}

// Handler to get a specific to-do list
func getTodoList(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID"})
		return
	}

	for _, list := range todoLists {
		if list.ID == listID {
			c.JSON(http.StatusOK, list)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found"})
}

// Handler to update a to-do list's title
func updateTodoList(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID"})
		return
	}

	var updatedList TodoList
	if err := c.ShouldBindJSON(&updatedList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, list := range todoLists {
		if list.ID == listID {
			todoLists[i].Title = updatedList.Title
			c.JSON(http.StatusOK, todoLists[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found"})
}

// Handler to delete a to-do list
func deleteTodoList(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID"})
		return
	}

	for i, list := range todoLists {
		if list.ID == listID {
			todoLists = append(todoLists[:i], todoLists[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "To-do list deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found"})
}

// Handler to create a new to-do item in a list
func createTodoItem(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list ID"})
		return
	}

	var newItem TodoItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newItem.ID = nextItemID
	nextItemID++
	newItem.CreatedAt = time.Now()

	for i, list := range todoLists {
		if list.ID == listID {
			todoLists[i].Items = append(todoLists[i].Items, newItem)
			c.JSON(http.StatusCreated, newItem)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "To-do list not found"})
}

// Handler to update a to-do item
func updateTodoItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var updatedItem TodoItem
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, list := range todoLists {
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

	c.JSON(http.StatusNotFound, gin.H{"error": "To-do item not found"})
}

// Handler to delete a to-do item
func deleteTodoItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	for i, list := range todoLists {
		for j, item := range list.Items {
			if item.ID == itemID {
				todoLists[i].Items = append(todoLists[i].Items[:j], todoLists[i].Items[j+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "To-do item deleted successfully"})
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "To-do item not found"})
}
