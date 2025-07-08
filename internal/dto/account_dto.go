package dto

type AccountDTO struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Balance      float64           `json:"balance"`
	Users        []*UserDTO        `json:"users,omitempty"`
	Transactions []*TransactionDTO `json:"transactions,omitempty"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
}

type AccountUpsert struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"max=255"`
}

type AccountUser struct {
	AccountID string `json:"accountId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
}
