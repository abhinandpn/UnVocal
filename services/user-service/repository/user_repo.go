package repository

import (
	"context"
	"errors"
	"time"

	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/abhinandpn/UnVocal/services/user-service/utils"
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
	GetUserByUserCode(userCode string) (*model.User, error)
	UserCodeExists(ctx context.Context, code string) (bool, error)
	GenerateUniqueUserCode(ctx context.Context) (string, error)
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

	ctx := context.Background()
	query := `
			INSERT INTO users (id, name, email, number, password, created_at , user_code) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
	// Generate a unique user code
	userCode, err := r.GenerateUniqueUserCode(ctx)
	if err != nil {
		return err
	}

	CreatedTime := time.Now()
	_, err = r.db.Exec(
		context.Background(),
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Number,
		user.Password,
		CreatedTime,
		userCode)

	return err
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password, created_at, user_code
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
		&user.UserCode,
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
		SELECT id, name, email, number, password, user_code
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
		&user.UserCode,
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
		SELECT id, name, email, number, password, user_code
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
		&user.UserCode,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByUserCode(userCode string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password, user_code
		FROM users
		WHERE user_code = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		context.Background(),
		query,
		userCode,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Number,
		&user.Password,
		&user.UserCode,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UserCodeExists(ctx context.Context, code string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_code = $1)`

	err := r.db.QueryRow(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *userRepository) GenerateUniqueUserCode(ctx context.Context) (string, error) {
	for {
		userCode, err := utils.GenerateUserCode()
		if err != nil {
			return "", err
		}

		exists, err := r.UserCodeExists(ctx, userCode)
		if err != nil {
			return "", err
		}

		if exists {
			continue // User code already exists, generate a new one
		}
		return userCode, nil
	}
}
