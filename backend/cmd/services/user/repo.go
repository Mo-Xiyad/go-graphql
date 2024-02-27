package services

import (
	"context"
	"server/pkg/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}
func (u *UserRepo) Create(ctx context.Context, user model.User) (*model.User, error) {
	if err := u.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (u *UserRepo) GetAll() ([]*model.User, error) {
// 	var users []*model.User
// 	if err := u.DB.Find(&users).Error; err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }

// func (u *UserRepo) Update(ctx context.Context, user model.User) (*model.User, error) {
// 	if err := u.DB.Save(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (u *UserRepo) Delete(ctx context.Context, id string) error {
// 	if err := u.DB.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (u *UserRepo) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
