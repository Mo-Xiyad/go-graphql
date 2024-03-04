package resolvers

import (
	"context"
	"errors"
	"net/http"
	"server"
	auth "server/cmd/services/auth"
	company "server/cmd/services/company"
	user "server/cmd/services/user"
	gql_model "server/graph/model"
	"server/pkg/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService    auth.IAuthService
	UserService    user.IUserService
	CompanyService company.ICompanyService
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}

func buildUnauthenticatedError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusUnauthorized,
		},
	}
}

func buildForbiddenError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildNotFoundError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, server.ErrForbidden):
		return buildForbiddenError(ctx, err)
	case errors.Is(err, server.ErrUnauthenticated):
		return buildUnauthenticatedError(ctx, err)
	case errors.Is(err, server.ErrValidation):
		return buildBadRequestError(ctx, err)
	case errors.Is(err, server.ErrNotFound):
		return buildNotFoundError(ctx, err)
	default:
		return err
	}
}
func mapUser(u model.User) *model.User {
	return &model.User{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
func mapAuthResponse(a gql_model.AuthPayload) *gql_model.AuthPayload {
	return &gql_model.AuthPayload{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		User:         mapUser(*a.User),
	}
}
