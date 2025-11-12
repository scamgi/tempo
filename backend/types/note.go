package types

import "time"

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateNotePayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNotePayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
