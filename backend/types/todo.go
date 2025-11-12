package types

import "time"

type TodoList struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

type TodoItem struct {
	ID          int        `json:"id"`
	ListID      int        `json:"listId"`
	Task        string     `json:"task"`
	IsCompleted bool       `json:"isCompleted"`
	DueDate     *time.Time `json:"dueDate,omitempty"` // Use a pointer for optional fields
	Priority    int        `json:"priority"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// Payloads for creating data
type CreateTodoListPayload struct {
	Title string `json:"title"`
}

type CreateTodoItemPayload struct {
	Task    string     `json:"task"`
	DueDate *time.Time `json:"dueDate"`
}

// Payload for updating a todo item
type UpdateTodoItemPayload struct {
	Task        *string `json:"task"`
	IsCompleted *bool   `json:"isCompleted"`
}
