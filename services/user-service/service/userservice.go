package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abhinandpn/UnVocal/services/user-service/auth"
	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/abhinandpn/UnVocal/services/user-service/repository"
	"github.com/abhinandpn/UnVocal/services/user-service/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(ctx context.Context, name, email, password, number string) error {

	// Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	// Check if phone number already exists
	existingUser, err = s.repo.GetUserByNumber(ctx, number)
	if err != nil {
		return fmt.Errorf("failed to check number: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with number %s already exists", number)
	}

	// TODO: Hash the password before storing it
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	newUser := &model.User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Number:   number,
		Password: string(hashedPassword), // Replace with string(hashedPassword)
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
func (s *UserService) GetUserByUserCode(ctx context.Context, userCode string) (*model.UserResponse, error) {

	isDeleted, err := s.repo.IsUserDeleted(ctx, userCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user is deleted: %w", err)
	}
	if isDeleted {
		return nil, fmt.Errorf("user with code %s is deleted", userCode)
	}

	user := model.UserResponse{}
	data, err := s.repo.GetUserByUserCode(ctx, userCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by user code: %w", err)
	}
	if data == nil {
		return nil, fmt.Errorf("user with code %s not found", userCode)
	}
	user.ID = data.ID
	user.Name = data.Name
	user.Email = data.Email
	user.Number = data.Number
	user.UserCode = data.UserCode

	return &user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userCode string) error {

	user, err := s.repo.GetUserByUserCode(ctx, userCode)
	if err != nil {
		return fmt.Errorf("failed to get user by user code: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user with code %s not found", userCode)
	}
	isDeleted, err := s.repo.IsUserDeleted(ctx, userCode)
	if err != nil {
		return fmt.Errorf("failed to check if user is deleted: %w", err)
	}
	if isDeleted {
		return fmt.Errorf("user with code %s is deleted", userCode)
	}
	return s.repo.DeleteUser(user.ID)
}

func (s *UserService) Login(ctx context.Context, identifier, password string) (*model.LoginResponse, error) {
	var (
		user *model.User
		err  error
	)

	// Find user by email or user code
	if strings.Contains(identifier, "@") {
		user, err = s.repo.GetUserByEmail(ctx, identifier)
	} else {
		user, err = s.repo.GetUserByUserCode(ctx, identifier)
	}

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email/user code or password")
	}

	// Verify password
	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, errors.New("invalid email/user code or password")
	}

	// Generate Access Token
	accessToken, err := auth.GenerateAccessToken(user.UserCode)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken, expiresAt, err := auth.GenerateRefreshToken(user.UserCode)
	if err != nil {
		return nil, err
	}

	// Save Refresh Token
	err = s.repo.CreateRefreshToken(ctx, &model.RefreshToken{
		UserCode:  user.UserCode,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	// Build user response
	userResponse := &model.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Number:   user.Number,
		UserCode: user.UserCode,
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, userCode string) error {
	// Check if user exists
	user, err := s.repo.GetUserByUserCode(ctx, userCode)
	if err != nil {
		return fmt.Errorf("failed to get user by user code: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}
	// Invalidate the token (implementation depends on your token management strategy)
	// For example, you might store invalidated tokens in a database or cache.
	// Here, we'll just return nil to indicate success.
	return nil
}
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.UpdateUser(user)
}
func (s *UserService) UserProfile(ctx context.Context, userCode string) (*model.UserResponse, error) {

	user, err := s.repo.GetUserByUserCode(ctx, userCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by user code: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with code %s not found", userCode)
	}

	isDeleted, err := s.repo.IsUserDeleted(ctx, userCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user is deleted: %w", err)
	}
	if isDeleted {
		return nil, fmt.Errorf("user with code %s is deleted", userCode)
	}

	userResponse := &model.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Number:   user.Number,
		UserCode: user.UserCode,
	}

	return userResponse, nil
}


func (s *UserService) RefreshToken(ctx context.Context,tokenString string,) (*model.LoginResponse, error) {

	// Validate JWT using JWT_REFRESH_SECRET
	claims, err := auth.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Check that token exists, is not revoked, and is not expired
	storedToken, err := s.repo.GetRefreshToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check refresh token: %w", err)
	}
	if storedToken == nil {
		return nil, errors.New("refresh token is invalid, expired, or revoked")
	}

	if storedToken.UserCode != claims.UserCode {
		return nil, errors.New("refresh token does not belong to this user")
	}

	user, err := s.repo.GetUserByUserCode(ctx, claims.UserCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	isDeleted, err := s.repo.IsUserDeleted(ctx, claims.UserCode)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, errors.New("user account is deleted")
	}

	// Revoke the old refresh token (token rotation)
	if err := s.repo.RevokeRefreshToken(ctx, tokenString); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	newAccessToken, err := auth.GenerateAccessToken(user.UserCode)
	if err != nil {
		return nil, err
	}

	newRefreshToken, expiresAt, err :=
		auth.GenerateRefreshToken(user.UserCode)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshToken(ctx, &model.RefreshToken{
		UserCode:  user.UserCode,
		Token:     newRefreshToken,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: &model.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Number:   user.Number,
			UserCode: user.UserCode,
		},
	}, nil
}