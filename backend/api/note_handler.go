package api

import (
	"log"
	"net/http"
	"strconv"
	"tempo-backend/db"
	"tempo-backend/types"

	"github.com/labstack/echo/v4"
)

type NoteHandler struct {
	store *db.NoteStore
}

func NewNoteHandler(store *db.NoteStore) *NoteHandler {
	return &NoteHandler{store: store}
}

func (h *NoteHandler) HandleCreateNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	var payload types.CreateNotePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}
	if payload.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title is required")
	}

	note, err := h.store.CreateNote(payload, userID)
	if err != nil {
		log.Printf("Error creating note: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create note")
	}
	return c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) HandleGetNotes(c echo.Context) error {
	userID := c.Get("userID").(int)
	notes, err := h.store.GetNotesByUser(userID)
	if err != nil {
		log.Printf("Error getting notes: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve notes")
	}
	return c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) HandleGetNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	noteID, err := strconv.Atoi(c.Param("noteId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid note ID")
	}

	note, err := h.store.GetNoteByID(noteID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Note not found")
	}

	return c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) HandleUpdateNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	noteID, err := strconv.Atoi(c.Param("noteId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid note ID")
	}

	var payload types.UpdateNotePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	note, err := h.store.UpdateNote(noteID, userID, payload)
	if err != nil {
		log.Printf("Error updating note: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update note")
	}
	return c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) HandleDeleteNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	noteID, err := strconv.Atoi(c.Param("noteId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid note ID")
	}

	err = h.store.DeleteNote(noteID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Note not found or not authorized")
	}

	return c.NoContent(http.StatusNoContent)
}