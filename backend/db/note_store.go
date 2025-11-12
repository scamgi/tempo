package db

import (
	"context"
	"fmt"
	"strings"
	"tempo-backend/types"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteStore struct {
	db *pgxpool.Pool
}

func NewNoteStore(db *pgxpool.Pool) *NoteStore {
	return &NoteStore{db: db}
}

func (s *NoteStore) CreateNote(payload types.CreateNotePayload, userID int) (*types.Note, error) {
	query := `INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3)
			   RETURNING id, user_id, title, content, created_at, updated_at`
	var note types.Note
	err := s.db.QueryRow(context.Background(), query, userID, payload.Title, payload.Content).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	return &note, err
}

func (s *NoteStore) GetNotesByUser(userID int) ([]types.Note, error) {
	query := `SELECT id, user_id, title, content, created_at, updated_at
			   FROM notes WHERE user_id = $1 ORDER BY updated_at DESC`
	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]types.Note, 0)
	for rows.Next() {
		var note types.Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (s *NoteStore) GetNoteByID(noteID, userID int) (*types.Note, error) {
	query := `SELECT id, user_id, title, content, created_at, updated_at
			   FROM notes WHERE id = $1 AND user_id = $2`
	var note types.Note
	err := s.db.QueryRow(context.Background(), query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	return &note, err
}

func (s *NoteStore) UpdateNote(noteID, userID int, payload types.UpdateNotePayload) (*types.Note, error) {
	var setParts []string
	var args []interface{}
	argID := 1

	if payload.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, *payload.Title)
		argID++
	}
	if payload.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argID))
		args = append(args, *payload.Content)
		argID++
	}
	if len(setParts) == 0 {
		return s.GetNoteByID(noteID, userID) // No update, just return the note
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now())
	argID++

	args = append(args, noteID, userID)
	query := fmt.Sprintf(`UPDATE notes SET %s WHERE id = $%d AND user_id = $%d
						   RETURNING id, user_id, title, content, created_at, updated_at`,
		strings.Join(setParts, ", "), argID, argID+1)

	var note types.Note
	err := s.db.QueryRow(context.Background(), query, args...).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	return &note, err
}

func (s *NoteStore) DeleteNote(noteID, userID int) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
	cmd, err := s.db.Exec(context.Background(), query, noteID, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("note not found or user not authorized")
	}
	return nil
}