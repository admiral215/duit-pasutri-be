package utils

import (
	"duit-pasutri-be/internal/dto"
	"gorm.io/gorm"
)

func PaginateWithScope[T any](
	db *gorm.DB,
	model any,
	page int,
	limit int,
	scope func(*gorm.DB) *gorm.DB,
) (*dto.PaginationResult[T], error) {
	var totalRows int64

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := scope(db.Model(model))

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, err
	}

	var result []T
	if err := query.Limit(limit).Offset(offset).Find(&result).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalRows + int64(limit) - 1) / int64(limit))

	return &dto.PaginationResult[T]{
		Page:       int32(page),
		Limit:      int32(limit),
		TotalRows:  int32(totalRows),
		TotalPages: int32(totalPages),
		Data:       result,
	}, nil
}
