package db

import (
	"context"
	"fmt"
	"strings"
	"tempo-backend/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoStore struct {
	db *pgxpool.Pool
}

func NewTodoStore(db *pgxpool.Pool) *TodoStore {
	return &TodoStore{db: db}
}

// --- ToDo List Methods ---

// CreateTodoList creates a new to-do list for a specific user.
func (s *TodoStore) CreateTodoList(payload types.CreateTodoListPayload, userID int) (*types.TodoList, error) {
	query := `INSERT INTO todo_lists (title, user_id) VALUES ($1, $2)
			   RETURNING id, user_id, title, created_at`
	var list types.TodoList
	err := s.db.QueryRow(context.Background(), query, payload.Title, userID).Scan(
		&list.ID, &list.UserID, &list.Title, &list.CreatedAt,
	)
	return &list, err
}

// GetTodoListsByUser retrieves all to-do lists for a given user.
func (s *TodoStore) GetTodoListsByUser(userID int) ([]types.TodoList, error) {
	query := `SELECT id, user_id, title, created_at FROM todo_lists WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := make([]types.TodoList, 0)
	for rows.Next() {
		var list types.TodoList
		if err := rows.Scan(&list.ID, &list.UserID, &list.Title, &list.CreatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}

// GetTodoListByID retrieves a single to-do list, ensuring it belongs to the correct user.
func (s *TodoStore) GetTodoListByID(listID, userID int) (*types.TodoList, error) {
	query := `SELECT id, user_id, title, created_at FROM todo_lists WHERE id = $1 AND user_id = $2`
	var list types.TodoList
	err := s.db.QueryRow(context.Background(), query, listID, userID).Scan(
		&list.ID, &list.UserID, &list.Title, &list.CreatedAt,
	)
	return &list, err
}

// DeleteTodoList deletes a list, ensuring it belongs to the correct user.
func (s *TodoStore) DeleteTodoList(listID, userID int) error {
	query := `DELETE FROM todo_lists WHERE id = $1 AND user_id = $2`
	cmd, err := s.db.Exec(context.Background(), query, listID, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("list not found or user not authorized")
	}
	return nil
}

// --- ToDo Item Methods ---

// CreateTodoItem adds a new task to a specific to-do list.
func (s *TodoStore) CreateTodoItem(payload types.CreateTodoItemPayload, listID int) (*types.TodoItem, error) {
	query := `INSERT INTO todo_items (list_id, task, due_date) VALUES ($1, $2, $3)
			   RETURNING id, list_id, task, is_completed, due_date, priority, created_at`
	var item types.TodoItem
	err := s.db.QueryRow(context.Background(), query, listID, payload.Task, payload.DueDate).Scan(
		&item.ID, &item.ListID, &item.Task, &item.IsCompleted, &item.DueDate, &item.Priority, &item.CreatedAt,
	)
	return &item, err
}

// GetTodoItemsByListID retrieves all items for a given to-do list.
func (s *TodoStore) GetTodoItemsByListID(listID int) ([]types.TodoItem, error) {
	query := `SELECT id, list_id, task, is_completed, due_date, priority, created_at FROM todo_items
			   WHERE list_id = $1 ORDER BY created_at ASC`
	rows, err := s.db.Query(context.Background(), query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]types.TodoItem, 0)
	for rows.Next() {
		var item types.TodoItem
		if err := rows.Scan(
			&item.ID, &item.ListID, &item.Task, &item.IsCompleted, &item.DueDate, &item.Priority, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// UpdateTodoItem updates a specific todo item.
func (s *TodoStore) UpdateTodoItem(itemID int, payload types.UpdateTodoItemPayload) (*types.TodoItem, error) {
	// Dynamically build the SET part of the query
	var setParts []string
	var args []interface{}
	argID := 1

	if payload.Task != nil {
		setParts = append(setParts, fmt.Sprintf("task = $%d", argID))
		args = append(args, *payload.Task)
		argID++
	}
	if payload.IsCompleted != nil {
		setParts = append(setParts, fmt.Sprintf("is_completed = $%d", argID))
		args = append(args, *payload.IsCompleted)
		argID++
	}
	if len(setParts) == 0 {
		return nil, fmt.Errorf("no update fields provided")
	}

	args = append(args, itemID)
	query := fmt.Sprintf(`UPDATE todo_items SET %s WHERE id = $%d
						   RETURNING id, list_id, task, is_completed, due_date, priority, created_at`,
		strings.Join(setParts, ", "), argID)

	var item types.TodoItem
	err := s.db.QueryRow(context.Background(), query, args...).Scan(
		&item.ID, &item.ListID, &item.Task, &item.IsCompleted, &item.DueDate, &item.Priority, &item.CreatedAt,
	)
	return &item, err
}

// DeleteTodoItem deletes a specific todo item.
func (s *TodoStore) DeleteTodoItem(itemID int) error {
	query := `DELETE FROM todo_items WHERE id = $1`
	cmd, err := s.db.Exec(context.Background(), query, itemID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}
