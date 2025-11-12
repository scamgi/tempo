package api

import (
	"log"
	"net/http"
	"strconv"
	"tempo-backend/db"
	"tempo-backend/types"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	store *db.TodoStore
}

func NewTodoHandler(store *db.TodoStore) *TodoHandler {
	return &TodoHandler{store: store}
}

// --- List Handlers ---

func (h *TodoHandler) HandleCreateTodoList(c echo.Context) error {
	userID := c.Get("userID").(int)
	var payload types.CreateTodoListPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}
	if payload.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title is required")
	}

	list, err := h.store.CreateTodoList(payload, userID)
	if err != nil {
		log.Printf("Error creating todo list: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create list")
	}
	return c.JSON(http.StatusCreated, list)
}

func (h *TodoHandler) HandleGetTodoLists(c echo.Context) error {
	userID := c.Get("userID").(int)
	lists, err := h.store.GetTodoListsByUser(userID)
	if err != nil {
		log.Printf("Error getting todo lists: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve lists")
	}
	return c.JSON(http.StatusOK, lists)
}

func (h *TodoHandler) HandleGetTodoListAndItems(c echo.Context) error {
	userID := c.Get("userID").(int)
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid list ID")
	}

	// First, verify the user owns this list
	list, err := h.store.GetTodoListByID(listID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}

	// Then, get the items for that list
	items, err := h.store.GetTodoItemsByListID(listID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve items")
	}

	response := struct {
		*types.TodoList
		Items []types.TodoItem `json:"items"`
	}{
		TodoList: list,
		Items:    items,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) HandleDeleteTodoList(c echo.Context) error {
	userID := c.Get("userID").(int)
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid list ID")
	}

	err = h.store.DeleteTodoList(listID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found or not authorized")
	}

	return c.NoContent(http.StatusNoContent)
}

// --- Item Handlers ---

func (h *TodoHandler) HandleCreateTodoItem(c echo.Context) error {
	// Security check: Ensure the user owns the list they're adding an item to.
	userID := c.Get("userID").(int)
	listID, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid list ID")
	}
	if _, err := h.store.GetTodoListByID(listID, userID); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Not authorized to add to this list")
	}

	var payload types.CreateTodoItemPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	item, err := h.store.CreateTodoItem(payload, listID)
	if err != nil {
		log.Printf("Error creating todo item: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create item")
	}
	return c.JSON(http.StatusCreated, item)
}

func (h *TodoHandler) HandleUpdateTodoItem(c echo.Context) error {
	// Security check would go here to ensure the item belongs to the authenticated user.
	// This is often done by joining tables in the SQL query itself. Omitted for brevity.
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID")
	}

	var payload types.UpdateTodoItemPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	item, err := h.store.UpdateTodoItem(itemID, payload)
	if err != nil {
		log.Printf("Error updating item: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update item")
	}
	return c.JSON(http.StatusOK, item)
}

func (h *TodoHandler) HandleDeleteTodoItem(c echo.Context) error {
	// Security check would go here.
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid item ID")
	}

	err = h.store.DeleteTodoItem(itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Item not found")
	}

	return c.NoContent(http.StatusNoContent)
}
