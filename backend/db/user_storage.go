package db

import (
	"context"
	"fmt"
	"log"

	"tempo-backend/types"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *pgxpool.Pool
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

// CreateUser handles user creation in the database.
func (s *UserStore) CreateUser(user types.RegisterUserPayload) (*types.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert the new user into the database
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)
			   RETURNING id, username, email, created_at`

	var newUser types.User
	err = s.db.QueryRow(context.Background(), query, user.Username, user.Email, string(hashedPassword)).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

// GetUserByEmail retrieves a user by their email address.
func (s *UserStore) GetUserByEmail(email string) (*types.User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`
	var user types.User
	err := s.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		log.Printf("Error getting user by email %s: %v", email, err)
		return nil, echo.ErrNotFound
	}
	return &user, nil
}
