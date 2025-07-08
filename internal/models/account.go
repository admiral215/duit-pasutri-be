package models

import (
	"duit-pasutri-be/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID           uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string         `gorm:"unique;type:varchar(255)"`
	Description  string         `gorm:"type:varchar(255)"`
	Balance      float64        `gorm:"type:numeric(10,2);default:0"`
	Users        []*User        `gorm:"many2many:user_accounts;"`
	Transactions []*Transaction `gorm:"foreignKey:AccountID;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a Account) ToDTO() *dto.AccountDTO {
	users := make([]*dto.UserDTO, len(a.Users))
	if len(a.Users) > 0 {
		for i, user := range a.Users {
			users[i] = user.ToDTO()
		}
	}

	transactions := make([]*dto.TransactionDTO, len(a.Transactions))
	if len(transactions) > 0 {
		for i, transaction := range a.Transactions {
			transactions[i] = transaction.ToDTO()
		}
	}

	return &dto.AccountDTO{
		ID:           a.ID.String(),
		Name:         a.Name,
		Description:  a.Description,
		Balance:      a.Balance,
		Users:        users,
		Transactions: transactions,
		CreatedAt:    a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    a.UpdatedAt.Format(time.RFC3339),
	}
}
