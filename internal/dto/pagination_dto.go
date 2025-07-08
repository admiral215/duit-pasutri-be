package dto

type PaginationInput struct {
	Page      int32  `json:"page"`
	Limit     int32  `json:"limit"`
	Search    string `json:"search"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

type PaginationResult[T any] struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalRows  int32 `json:"totalRows"`
	TotalPages int32 `json:"totalPages"`
	Data       []T   `json:"data"`
}

type PaginatedCategory struct {
	Page       int32          `json:"page"`
	Limit      int32          `json:"limit"`
	TotalRows  int32          `json:"totalRows"`
	TotalPages int32          `json:"totalPages"`
	Data       []*CategoryDTO `json:"data"`
}

type PaginatedUser struct {
	Page       int32      `json:"page"`
	Limit      int32      `json:"limit"`
	TotalRows  int32      `json:"totalRows"`
	TotalPages int32      `json:"totalPages"`
	Data       []*UserDTO `json:"data"`
}

type PaginatedAccount struct {
	Page       int32         `json:"page"`
	Limit      int32         `json:"limit"`
	TotalRows  int32         `json:"totalRows"`
	TotalPages int32         `json:"totalPages"`
	Data       []*AccountDTO `json:"data"`
}

type PaginatedTransaction struct {
	Page       int32             `json:"page"`
	Limit      int32             `json:"limit"`
	TotalRows  int32             `json:"totalRows"`
	TotalPages int32             `json:"totalPages"`
	Data       []*TransactionDTO `json:"data"`
}
