package service

import (
	"context"
	"testing"

	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/abhinandpn/UnVocal/services/user-service/repository"
)

type mockUserRepository struct {
	// Embedding satisfies methods that a particular test does not use.
	repository.UserRepository

	getUserByEmailFn  func(context.Context, string) (*model.User, error)
	getUserByNumberFn func(context.Context, string) (*model.User, error)
	createUserFn      func(*model.User) error
}

func (m *mockUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	return m.getUserByEmailFn(ctx, email)
}

func (m *mockUserRepository) GetUserByNumber(
	ctx context.Context,
	number string,
) (*model.User, error) {
	return m.getUserByNumberFn(ctx, number)
}

func (m *mockUserRepository) CreateUser(user *model.User) error {
	return m.createUserFn(user)
}
func TestRegisterSuccess(t *testing.T) {
	var createdUser *model.User

	repo := &mockUserRepository{
		getUserByEmailFn: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return nil, nil
		},

		getUserByNumberFn: func(
			ctx context.Context,
			number string,
		) (*model.User, error) {
			return nil, nil
		},

		createUserFn: func(user *model.User) error {
			createdUser = user
			return nil
		},
	}

	service := NewUserService(repo)

	err := service.Register(
		context.Background(),
		"Abhinand",
		"abhinand@example.com",
		"Password123",
		"8086324067",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdUser == nil {
		t.Fatal("expected user to be created")
	}

	if createdUser.Email != "abhinand@example.com" {
		t.Errorf(
			"expected email abhinand@example.com, got %s",
			createdUser.Email,
		)
	}

	if createdUser.Password == "Password123" {
		t.Error("expected password to be hashed")
	}
}
func TestRegisterDuplicateEmail(t *testing.T) {
	repo := &mockUserRepository{
		getUserByEmailFn: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return &model.User{Email: email}, nil
		},
	}

	service := NewUserService(repo)

	err := service.Register(
		context.Background(),
		"Abhinand",
		"abhinand@example.com",
		"Password123",
		"8086324067",
	)

	if err == nil {
		t.Fatal("expected duplicate-email error")
	}
}
