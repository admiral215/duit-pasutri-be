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

type UserRepository interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	FindAllPagination(ctx context.Context, input *dto.PaginationInput) (*dto.PaginationResult[*models.User], error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, tx *gorm.DB, user *models.User) (*models.User, error)
	Update(ctx context.Context, tx *gorm.DB, user *models.User) error
}

type userRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		DB: db,
	}
}

func (u userRepositoryImpl) FindAllPagination(ctx context.Context, input *dto.PaginationInput) (*dto.PaginationResult[*models.User], error) {
	return utils.PaginateWithScope[*models.User](u.DB, &models.User{}, int(input.Page), int(input.Limit), func(tx *gorm.DB) *gorm.DB {
		allowedSortFields := map[string]bool{
			"name":       true,
			"email":      true,
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
				Or("email ILIKE ?", "%"+input.Search+"%")
		}

		return tx.Order(input.SortBy + " " + input.SortOrder)
	})
}

func (u userRepositoryImpl) FindAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := u.DB.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := u.DB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := u.DB.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, user *models.User) (*models.User, error) {
	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u userRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, user *models.User) error {
	err := tx.WithContext(ctx).Updates(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
