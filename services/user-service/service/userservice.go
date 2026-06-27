package service

import (
	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/abhinandpn/UnVocal/services/user-service/repository"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}
func (s *UserService) Register(name, email, password string, number string) error {

	user := model.User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Number:   number,
		Password: password,
	}

	return s.repo.CreateUser(&user)
}
func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.repo.GetUserByID(id)
}
func (s *UserService) UpdateUser(user *model.User) error {
	return s.repo.UpdateUser(user)
}
func (s *UserService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
