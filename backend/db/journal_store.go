package db

import (
	"context"
	"fmt"
	"strings"
	"tempo-backend/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JournalStore struct {
	db *pgxpool.Pool
}

func NewJournalStore(db *pgxpool.Pool) *JournalStore {
	return &JournalStore{db: db}
}

func (s *JournalStore) CreateJournalEntry(payload types.CreateJournalEntryPayload, userID int) (*types.JournalEntry, error) {
	query := `INSERT INTO journal_entries (user_id, title, content, mood, entry_date) VALUES ($1, $2, $3, $4, $5)
			   RETURNING id, user_id, title, content, mood, entry_date, created_at`
	var entry types.JournalEntry
	err := s.db.QueryRow(context.Background(), query, userID, payload.Title, payload.Content, payload.Mood, payload.EntryDate).Scan(
		&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &entry.Mood, &entry.EntryDate, &entry.CreatedAt,
	)
	return &entry, err
}

func (s *JournalStore) GetJournalEntriesByUser(userID int) ([]types.JournalEntry, error) {
	query := `SELECT id, user_id, title, content, mood, entry_date, created_at
			   FROM journal_entries WHERE user_id = $1 ORDER BY entry_date DESC`
	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]types.JournalEntry, 0)
	for rows.Next() {
		var entry types.JournalEntry
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &entry.Mood, &entry.EntryDate, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *JournalStore) GetJournalEntryByID(entryID, userID int) (*types.JournalEntry, error) {
	query := `SELECT id, user_id, title, content, mood, entry_date, created_at
			   FROM journal_entries WHERE id = $1 AND user_id = $2`
	var entry types.JournalEntry
	err := s.db.QueryRow(context.Background(), query, entryID, userID).Scan(
		&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &entry.Mood, &entry.EntryDate, &entry.CreatedAt,
	)
	return &entry, err
}

func (s *JournalStore) UpdateJournalEntry(entryID, userID int, payload types.UpdateJournalEntryPayload) (*types.JournalEntry, error) {
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
	if payload.Mood != nil {
		setParts = append(setParts, fmt.Sprintf("mood = $%d", argID))
		args = append(args, *payload.Mood)
		argID++
	}
	if payload.EntryDate != nil {
		setParts = append(setParts, fmt.Sprintf("entry_date = $%d", argID))
		args = append(args, *payload.EntryDate)
		argID++
	}
	if len(setParts) == 0 {
		return s.GetJournalEntryByID(entryID, userID)
	}

	args = append(args, entryID, userID)
	query := fmt.Sprintf(`UPDATE journal_entries SET %s WHERE id = $%d AND user_id = $%d
						   RETURNING id, user_id, title, content, mood, entry_date, created_at`,
		strings.Join(setParts, ", "), argID, argID+1)

	var entry types.JournalEntry
	err := s.db.QueryRow(context.Background(), query, args...).Scan(
		&entry.ID, &entry.UserID, &entry.Title, &entry.Content, &entry.Mood, &entry.EntryDate, &entry.CreatedAt,
	)
	return &entry, err
}

func (s *JournalStore) DeleteJournalEntry(entryID, userID int) error {
	query := `DELETE FROM journal_entries WHERE id = $1 AND user_id = $2`
	cmd, err := s.db.Exec(context.Background(), query, entryID, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("entry not found or user not authorized")
	}
	return nil
}