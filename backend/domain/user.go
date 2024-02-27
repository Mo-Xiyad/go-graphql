package domain

import (
	"context"
	gql_model "server/graph/model"
	"server/pkg/model"
	"server/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo types.UserRepo
}

func NewUserService(ur types.UserRepo) *UserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (us *UserService) CreateUser(ctx context.Context, input gql_model.NewUser) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:       uint64(uuid.New().ID()),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
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
