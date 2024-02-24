package api

import (
	"context"
	"fmt"
	"math/rand"
	model1 "server/graph/model"
	"server/pkg/model"

	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, input model1.NewUser) (*model.User, error) {
	db, ok := ctx.Value("DB").(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to get database instance from context")
	}

	user := model.User{
		ID:    rand.Uint64(),
		Name:  input.Name,
		Email: input.Email,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &model.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func GetAllUsers(ctx context.Context) ([]*model.User, error) {
	db, ok := ctx.Value("DB").(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to get database instance from context")
	}

	var usersFromDb []*model.User
	if err := db.Find(&usersFromDb).Error; err != nil {
		return nil, err
	}

	var users []*model.User
	for _, userModel := range usersFromDb {
		user := &model.User{
			ID:    userModel.ID,
			Name:  userModel.Name,
			Email: userModel.Email,
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUser(ctx context.Context, id string) (*model.User, error) {
	db, ok := ctx.Value("DB").(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to get database instance from context")
	}

	var userFromDb model.User
	if err := db.Where("id = ?", id).First(&userFromDb).Error; err != nil {
		return nil, err
	}

	user := &model.User{
		ID:    userFromDb.ID,
		Name:  userFromDb.Name,
		Email: userFromDb.Email,
	}
	return user, nil
}
