package dto

type TransactionDTO struct {
	ID          string       `json:"id"`
	Amount      float64      `json:"amount"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Account     *AccountDTO  `json:"account"`
	Category    *CategoryDTO `json:"category"`
	User        *UserDTO     `json:"user"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
}

type TransactionUpsert struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount" validate:"min=0"`
	Description string  `json:"description" validate:"required"`
	Type        string  `json:"type" validate:"required"`
	AccountID   string  `json:"accountId" validate:"required"`
	CategoryID  string  `json:"categoryId" validate:"required"`
}
