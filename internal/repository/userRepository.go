package repository

import (
	"ecommerce/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(usr domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindById(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(usr domain.User) (domain.User, error) {
	err := r.db.Create(&usr).Error
	if err != nil {
		log.Printf("create user error: %v", err)
		return domain.User{}, errors.New("create user failed")
	}
	return usr, nil
}

func (r *userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found by email")
		}
		return domain.User{}, errors.New("error occurred while finding user by email")
	}
	return user, nil
}

func (r *userRepository) FindById(id uint) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found by id")
		}
		return domain.User{}, errors.New("error occurred while finding user by id")
	}
	return user, nil
}

func (r *userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}

	err = r.db.Model(&user).Clauses(clause.Returning{}).Updates(u).Error
	if err != nil {
		return domain.User{}, errors.New("update user failed")
	}

	return user, nil
}
