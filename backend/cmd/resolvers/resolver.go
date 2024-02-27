package resolvers

import (
	"context"
	"net/http"
	auth "server/cmd/services/auth"
	user "server/cmd/services/user"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthService auth.IAuthService
	UserService user.IUserService
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

// func buildError(ctx context.Context, err error) error {
// 	switch {
// 	case errors.Is(err, twitter.ErrForbidden):
// 		return buildForbiddenError(ctx, err)
// 	case errors.Is(err, twitter.ErrUnauthenticated):
// 		return buildUnauthenticatedError(ctx, err)
// 	case errors.Is(err, twitter.ErrValidation):
// 		return buildBadRequestError(ctx, err)
// 	case errors.Is(err, twitter.ErrNotFound):
// 		return buildNotFoundError(ctx, err)
// 	default:
// 		return err
// 	}
// }
