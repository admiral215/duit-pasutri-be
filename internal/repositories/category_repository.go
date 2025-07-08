package repositories

import (
	"context"
	"duit-pasutri-be/internal/dto"
	"duit-pasutri-be/internal/models"
	"duit-pasutri-be/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, category *models.Category) error
	Update(ctx context.Context, tx *gorm.DB, category *models.Category) error
	Delete(ctx context.Context, tx *gorm.DB, categoryID uuid.UUID) error
	FindAll(ctx context.Context) ([]*models.Category, error)
	FindAllPagination(ctx context.Context, input *dto.PaginationInput) (*dto.PaginationResult[*models.Category], error)
	FindAllByType(ctx context.Context, categoryType string) (*[]models.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
}

type CategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewCategoryRepositoryImpl(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{
		DB: db,
	}
}

func (c CategoryRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, category *models.Category) error {
	return tx.WithContext(ctx).Create(category).Error
}

func (c CategoryRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, category *models.Category) error {
	return tx.WithContext(ctx).Updates(category).Error
}

func (c CategoryRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, categoryID uuid.UUID) error {
	result := tx.WithContext(ctx).Delete(&models.Category{}, "id = ?", categoryID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (c CategoryRepositoryImpl) FindAll(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	err := c.DB.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c CategoryRepositoryImpl) FindAllPagination(ctx context.Context, input *dto.PaginationInput) (*dto.PaginationResult[*models.Category], error) {
	return utils.PaginateWithScope[*models.Category](
		c.DB,
		&models.Category{},
		int(input.Page),
		int(input.Limit),
		func(tx *gorm.DB) *gorm.DB {
			allowedSortFields := map[string]bool{
				"name":       true,
				"created_at": true,
			}

			if input.SortBy == "" || !allowedSortFields[input.SortBy] {
				input.SortBy = "created_at"
			}

			if input.SortOrder == "" || (input.SortOrder != "asc" && input.SortOrder != "desc") {
				input.SortOrder = "desc"
			}

			if input.Search != "" {
				tx = tx.WithContext(ctx).
					Where("name ILIKE ?", "%"+input.Search+"%").
					Or("type = ?", input.Search)
			}

			return tx.Order(fmt.Sprintf("%s %s", input.SortBy, input.SortOrder))
		},
	)
}

func (c CategoryRepositoryImpl) FindAllByType(ctx context.Context, categoryType string) (*[]models.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (c CategoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := c.DB.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}
