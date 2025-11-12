package types

import "time"

type JournalEntry struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Mood      *string   `json:"mood"`
	EntryDate time.Time `json:"entryDate"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateJournalEntryPayload struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Mood      *string   `json:"mood"`
	EntryDate time.Time `json:"entryDate"`
}

type UpdateJournalEntryPayload struct {
	Title   *string    `json:"title"`
	Content *string    `json:"content"`
	Mood    *string    `json:"mood"`
	EntryDate *time.Time `json:"entryDate"`
}