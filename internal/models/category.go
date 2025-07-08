package models

import (
	"duit-pasutri-be/internal/dto"
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID   uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string          `gorm:"unique;not null"`
	Type TransactionType `gorm:"type:varchar(20);not null;check:type IN ('spending','income')"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c Category) ToDTO() *dto.CategoryDTO {
	return &dto.CategoryDTO{
		ID:        c.ID.String(),
		Name:      c.Name,
		Type:      string(c.Type),
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}
