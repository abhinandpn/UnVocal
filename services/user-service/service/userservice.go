package service

import (
	"context"
	"fmt"

	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/abhinandpn/UnVocal/services/user-service/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(ctx context.Context, name, email, password, number string) error {

	// Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	// Check if phone number already exists
	existingUser, err = s.repo.GetUserByNumber(number)
	if err != nil {
		return fmt.Errorf("failed to check number: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with number %s already exists", number)
	}
	// TODO: Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
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
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.GetUserByID(id)
}
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.UpdateUser(user)
}
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(id)
}
