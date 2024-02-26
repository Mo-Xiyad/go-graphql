package resolvers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"server/cmd/api"
	"server/graph"
	gql_model "server/graph/model"
	"server/pkg/model"
	"server/types"
)

func mapAuthResponse(a gql_model.AuthPayload) *gql_model.AuthPayload {
	return &gql_model.AuthPayload{
		Token: a.Token,
		User:  mapUser(*a.User),
	}
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input gql_model.NewUser) (*model.User, error) {
	return api.CreateUser(ctx, input)
}

// CreateCompany is the resolver for the createCompany field.
func (r *mutationResolver) CreateCompany(ctx context.Context, input gql_model.CreateCompanyInput) (*model.Company, error) {
	panic(fmt.Errorf("not implemented: CreateCompany - createCompany"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input gql_model.LoginInput) (*gql_model.AuthPayload, error) {
	res, err := r.AuthService.Login(ctx, gql_model.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	log.Println("res", res)
	log.Println("err", err)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrValidation) ||
			errors.Is(err, types.ErrBadCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}
	log.Println("res", res)
	return mapAuthResponse(*res), nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
