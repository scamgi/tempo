package main

import "time"

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
