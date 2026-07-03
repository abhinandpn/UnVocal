package repository

import (
	"context"
	"errors"
	"time"

	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id string) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByNumber(number string) (*model.User, error)
}

// userRepository is the concrete implementation of UserRepository
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *model.User) error {

	query := `
			INSERT INTO users (id, name, email, number, password, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6)
		`

	CreatedTime := time.Now()
	_, err := r.db.Exec(
		context.Background(),
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Number,
		user.Password,
		CreatedTime,
	)

	return err
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password, created_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Number,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, number = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		user.Name,
		user.Email,
		user.Number,
		user.ID,
	)

	return err
}

func (r *userRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(
		context.Background(),
		query,
		id,
	)

	return err
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password
		FROM users
		WHERE email = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Number,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByNumber(number string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password
		FROM users
		WHERE number = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		context.Background(),
		query,
		number,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Number,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
