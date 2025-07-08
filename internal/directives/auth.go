package directives

import (
	"context"
	"duit-pasutri-be/internal/middleware"
	"duit-pasutri-be/internal/models"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
)

func Auth(ctx context.Context, obj any, next graphql.Resolver, roles []models.Role) (any, error) {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if models.Role(user.Role) == role {
			return next(ctx)
		}
	}

	return nil, fmt.Errorf("unauthorized: role %s is not allowed", user.Role)
}
