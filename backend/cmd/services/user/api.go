package services

import (
	"context"
	"errors"
	"fmt"
	"server"
	gql_model "server/graph/model"
	"server/pkg/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo IUserRepo
}

func NewUserService(ur IUserRepo) *UserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (us *UserService) CreateUser(ctx context.Context, input gql_model.CreateUserInput) (*model.User, error) {

	if _, err := us.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, server.ErrNotFound) {
		return nil, fmt.Errorf("email %s is already taken", input.Email)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user := model.User{
		ID:       uint64(uuid.New().ID()),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashPassword),
	}
	// create the user
	return us.UserRepo.Create(ctx, user)
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := us.UserRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := us.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := us.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
