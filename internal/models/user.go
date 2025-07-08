package models

import (
	"duit-pasutri-be/internal/dto"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string     `gorm:"type:varchar(255)"`
	Password  string     `gorm:"type:text"`
	Email     string     `gorm:"uniqueIndex"`
	Role      Role       `gorm:"type:varchar(20);not null;check:role IN ('admin','user')"`
	Accounts  []*Account `gorm:"many2many:user_accounts;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (user User) ToDTO() *dto.UserDTO {
	return &dto.UserDTO{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
