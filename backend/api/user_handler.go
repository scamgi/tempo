package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"tempo-backend/db"
	"tempo-backend/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store *db.UserStore
}

func NewUserHandler(store *db.UserStore) *UserHandler {
	return &UserHandler{store: store}
}

// HandleRegisterUser handles the user registration request.
func (h *UserHandler) HandleRegisterUser(c echo.Context) error {
	var payload types.RegisterUserPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Basic validation
	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username, email, and password are required")
	}

	// Create user in the database
	user, err := h.store.CreateUser(payload)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		// This could be a unique constraint violation (e.g., email already exists)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create user")
	}

	return c.JSON(http.StatusCreated, user)
}

// HandleLoginUser handles the user login request.
func (h *UserHandler) HandleLoginUser(c echo.Context) error {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password))
	if err != nil {
		// Passwords don't match
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// --- Generate JWT Token ---
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expires in 3 days

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
