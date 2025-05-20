package repository

import (
	"errors"
	"gin_im/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser 创建用户
func (repository *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	var err error = repository.db.Create(user).Error

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repository *UserRepository) GetUserById(id uint) (*model.User, error) {
	var user model.User
	var err error = repository.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repository *UserRepository) GetUserByName(name string) (*model.User, error) {
	var user model.User
	var err error = repository.db.Where("name = ?", name).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	var err error = repository.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	err := repository.db.Where("phone = ?").First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (repository *UserRepository) UpdateUser(user *model.User) error {
	return repository.db.Save(user).Error
}
