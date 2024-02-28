package services

import (
	"context"
	gql_model "server/graph/model"
	"server/pkg/model"
)

type IUserService interface {
	CreateUser(ctx context.Context, user gql_model.CreateUserInput) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
}

type IUserRepo interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}
