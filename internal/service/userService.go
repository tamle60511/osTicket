package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"errors"
)

type UserService struct {
	Repo repository.UserRepository
	Auth helper.Auth
}

func (s *UserService) Signup(input dto.UserSignup) (string, error) {
	hPassword, err := s.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}
	user, err := s.Repo.CreateUser(domain.User{
		UserName:     input.UserName,
		Email:        input.Email,
		Password:     hPassword,
		Phone:        input.Phone,
		DepartmentID: input.DepartmentID,
		Role:         domain.Role(input.Role),
	})
	if err != nil {
		return "", err
	}
	return s.Auth.GenerateToken(user.ID, user.Email, string(user.Role))
}

func (s *UserService) FindByEmail(email string) (*domain.User, error) {
	user, err := s.Repo.FindUser(email)
	if err != nil {
		return &domain.User{}, errors.New("find by email fail")
	}
	return &user, nil
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.FindByEmail(email)
	if err != nil {
		return "", errors.New("login fail")
	}
	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, string(user.Role))
}
