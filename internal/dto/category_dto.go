package dto

type CategoryDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CategoryUpsert struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=spending income"`
}
