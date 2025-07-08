package dto

import "github.com/google/uuid"

type UserDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserContext struct {
	UserID uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
}

type UserUpdate struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Register struct {
	Name            string `json:"name" validate:"required,gte=2,lte=100"`
	Password        string `json:"password" validate:"required,gte=8,lte=100"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,gte=8,lte=100"`
	Email           string `json:"email" validate:"required,email"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  *UserDTO `json:"user"`
}
