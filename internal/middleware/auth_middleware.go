package middleware

import (
	"context"
	"duit-pasutri-be/internal/auth"
	"duit-pasutri-be/internal/dto"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := auth.ValidateToken(tokenString)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user.Claims.(jwt.MapClaims))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*dto.UserContext, error) {
	claims, ok := ctx.Value(UserContextKey).(jwt.MapClaims)
	if !ok {
		return nil, errors.New("not logged in")
	}

	userId, err := uuid.Parse(claims["userId"].(string))
	if err != nil {
		return nil, errors.New("could not parse user id")
	}

	user := dto.UserContext{
		UserID: userId,
		Role:   claims["role"].(string),
	}
	return &user, nil
}
