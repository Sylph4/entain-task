package service

import (
	"github.com/sylph4/entain-task/internal/record/repository"
	"github.com/sylph4/entain-task/storage/postgres"
)

type IUserService interface {
	GetAllUsers() ([]*postgres.User, error)
}

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(
	userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetAllUsers() ([]postgres.User, error) {
	users, err := s.userRepository.SelectAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
