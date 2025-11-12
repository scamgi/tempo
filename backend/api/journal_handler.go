package api

import (
	"log"
	"net/http"
	"strconv"
	"tempo-backend/db"
	"tempo-backend/types"
	"time"

	"github.com/labstack/echo/v4"
)

type JournalHandler struct {
	store *db.JournalStore
}

func NewJournalHandler(store *db.JournalStore) *JournalHandler {
	return &JournalHandler{store: store}
}

func (h *JournalHandler) HandleCreateJournalEntry(c echo.Context) error {
	userID := c.Get("userID").(int)
	var payload types.CreateJournalEntryPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}
	if payload.Title == "" || payload.EntryDate.IsZero() {
		return echo.NewHTTPError(http.StatusBadRequest, "Title and entryDate are required")
	}
    // Truncate to just the date part
    payload.EntryDate = payload.EntryDate.Truncate(24 * time.Hour)


	entry, err := h.store.CreateJournalEntry(payload, userID)
	if err != nil {
		log.Printf("Error creating journal entry: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create journal entry")
	}
	return c.JSON(http.StatusCreated, entry)
}

func (h *JournalHandler) HandleGetJournalEntries(c echo.Context) error {
	userID := c.Get("userID").(int)
	entries, err := h.store.GetJournalEntriesByUser(userID)
	if err != nil {
		log.Printf("Error getting journal entries: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve journal entries")
	}
	return c.JSON(http.StatusOK, entries)
}

func (h *JournalHandler) HandleGetJournalEntry(c echo.Context) error {
	userID := c.Get("userID").(int)
	entryID, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid entry ID")
	}

	entry, err := h.store.GetJournalEntryByID(entryID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Journal entry not found")
	}

	return c.JSON(http.StatusOK, entry)
}

func (h *JournalHandler) HandleUpdateJournalEntry(c echo.Context) error {
	userID := c.Get("userID").(int)
	entryID, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid entry ID")
	}

	var payload types.UpdateJournalEntryPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	entry, err := h.store.UpdateJournalEntry(entryID, userID, payload)
	if err != nil {
		log.Printf("Error updating journal entry: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update journal entry")
	}
	return c.JSON(http.StatusOK, entry)
}

func (h *JournalHandler) HandleDeleteJournalEntry(c echo.Context) error {
	userID := c.Get("userID").(int)
	entryID, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid entry ID")
	}

	err = h.store.DeleteJournalEntry(entryID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Journal entry not found or not authorized")
	}

	return c.NoContent(http.StatusNoContent)
}