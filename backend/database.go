package main

import "time"

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
