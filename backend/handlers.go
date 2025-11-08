package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
	newList.UserID = userID
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
