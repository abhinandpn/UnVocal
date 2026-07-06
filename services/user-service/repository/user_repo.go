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

	// Define the methods for user repository operations
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id string) error

	// Get user by different attributes
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByUserCode(ctx context.Context, userCode string) (*model.User, error)
	GetUserByNumber(ctx context.Context, number string) (*model.User, error)

	// user code generation and existence check
	GenerateUniqueUserCode(ctx context.Context) (string, error)
	UserCodeExists(ctx context.Context, code string) (bool, error)

	// Check if a user is deleted
	IsUserDeleted(ctx context.Context, userCode string) (bool, error)

	// Refresh Tokens
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeRefreshTokensByUserCode(ctx context.Context, userCode string) error // optional
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

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password, created_at, user_code
		FROM users
		WHERE id = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		ctx,
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

	time := time.Now()
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2`

	_, err := r.db.Exec(
		context.Background(),
		query,
		time,
		id,
	)

	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (r *userRepository) GetUserByNumber(ctx context.Context, number string) (*model.User, error) {
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

func (r *userRepository) GetUserByUserCode(ctx context.Context, userCode string) (*model.User, error) {
	query := `
		SELECT id, name, email, number, password, user_code
		FROM users
		WHERE user_code = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(
		ctx,
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

func (r *userRepository) IsUserDeleted(ctx context.Context, UserCode string) (bool, error) {

	query := `SELECT deleted_at FROM users WHERE user_code = $1`

	var deletedAt *time.Time
	err := r.db.QueryRow(ctx, query, UserCode).Scan(&deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil // User not found
		}
		return false, err
	}

	return deletedAt != nil, nil
}

func (r *userRepository) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {

	query := `
		INSERT INTO refresh_tokens (
			user_code,
			token,
			expires_at,
			created_at
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		token.UserCode,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)

	return err
}

func (r *userRepository) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	query := `
		SELECT
			id,
			user_code,
			token,
			expires_at,
			created_at,
			revoked_at
		FROM refresh_tokens
		WHERE token = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
	`

	refreshToken := &model.RefreshToken{}

	err := r.db.QueryRow(
		ctx,
		query,
		token,
	).Scan(
		&refreshToken.ID,
		&refreshToken.UserCode,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.RevokedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return refreshToken, nil
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *userRepository) RevokeRefreshTokensByUserCode(ctx context.Context, userCode string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_code = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, userCode)
	return err
}
