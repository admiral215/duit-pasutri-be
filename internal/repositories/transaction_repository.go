package repositories

import (
	"context"
	"duit-pasutri-be/internal/dto"
	"duit-pasutri-be/internal/models"
	"duit-pasutri-be/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type TransactionRepository interface {
	FindAllPaginationByUserID(ctx context.Context, input *dto.PaginationInput, UserID uuid.UUID) (*dto.PaginationResult[*models.Transaction], error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	FindByIDWithPreload(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	Create(ctx context.Context, tx *gorm.DB, input *models.Transaction) error
	Update(ctx context.Context, tx *gorm.DB, input *models.Transaction) error
	Delete(ctx context.Context, tx *gorm.DB, input *models.Transaction) error
}

type transactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{
		DB: db,
	}
}

func (t transactionRepositoryImpl) FindAllPaginationByUserID(ctx context.Context, input *dto.PaginationInput, UserID uuid.UUID) (*dto.PaginationResult[*models.Transaction], error) {
	return utils.PaginateWithScope[*models.Transaction](
		t.DB,
		&models.Transaction{},
		int(input.Page),
		int(input.Limit),
		func(tx *gorm.DB) *gorm.DB {
			allowedSortFields := map[string]bool{
				"created_at": true,
			}

			if input.SortBy == "" || !allowedSortFields[input.SortBy] {
				input.SortBy = "created_at"
			}

			if input.SortOrder == "" || (input.SortOrder != "asc" && input.SortOrder != "desc") {
				input.SortOrder = "desc"
			}
			tx = tx.WithContext(ctx).
				Preload("Category").
				Preload("Account").
				Preload("User")
			if input.Search != "" && (strings.ToLower(input.Search) == "income" || strings.ToLower(input.Search) == "spending") {
				tx = tx.Where("type = ?", "%"+strings.ToLower(input.Search)+"%")
			}
			return tx.Order(input.SortBy + " " + input.SortOrder)
		})
}

func (t transactionRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction *models.Transaction
	err := t.DB.WithContext(ctx).First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t transactionRepositoryImpl) FindByIDWithPreload(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction *models.Transaction
	err := t.DB.WithContext(ctx).
		Preload("Category").
		Preload("Account").
		Preload("User").
		First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t transactionRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, input *models.Transaction) error {
	return tx.WithContext(ctx).Create(input).Error
}

func (t transactionRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, input *models.Transaction) error {
	return tx.WithContext(ctx).Updates(input).Error
}

func (t transactionRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, input *models.Transaction) error {
	result := tx.WithContext(ctx).Delete(input)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
