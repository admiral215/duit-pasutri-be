package models

import (
	"duit-pasutri-be/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Amount      float64         `gorm:"type:numeric(10,2)"`
	Description string          `gorm:"type:text"`
	Type        TransactionType `gorm:"type:varchar(20);not null;check:type IN ('spending','income')"`
	CategoryID  uuid.UUID
	AccountID   uuid.UUID
	UserID      uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Account  *Account  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User     *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TransactionType string

const (
	Spending TransactionType = "spending"
	Income   TransactionType = "income"
)

func (t Transaction) ToDTO() *dto.TransactionDTO {
	return &dto.TransactionDTO{
		ID:          t.ID.String(),
		Amount:      t.Amount,
		Description: t.Description,
		Type:        string(t.Type),
		Account: &dto.AccountDTO{
			ID:   t.AccountID.String(),
			Name: t.Account.Name,
		},
		Category: &dto.CategoryDTO{
			ID:   t.Category.ID.String(),
			Name: t.Category.Name,
		},
		User: &dto.UserDTO{
			ID:    t.User.ID.String(),
			Name:  t.User.Name,
			Email: t.User.Email,
		},
	}
}
