package resolvers

import (
	"context"
	"errors"
	"fmt"
	"server"
	"server/graph"
	gql_model "server/graph/model"
	"server/pkg/model"
)

func mapAuthResponse(a gql_model.AuthPayload) *gql_model.AuthPayload {
	return &gql_model.AuthPayload{
		Token: a.Token,
		User:  mapUser(*a.User),
	}
}

func (r *mutationResolver) CreateUser(ctx context.Context, input gql_model.NewUser) (*model.User, error) {
	res, err := r.UserService.CreateUser(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, server.ErrValidation):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}
	return res, nil
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
	if err != nil {
		switch {
		case errors.Is(err, server.ErrValidation) ||
			errors.Is(err, server.ErrBadCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}
	return mapAuthResponse(*res), nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
