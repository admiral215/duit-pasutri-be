package repositories

import (
	"context"
	"duit-pasutri-be/internal/dto"
	"duit-pasutri-be/internal/models"
	"duit-pasutri-be/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, tx *gorm.DB, account *models.Account) error
	Update(ctx context.Context, tx *gorm.DB, account *models.Account) error
	Delete(ctx context.Context, tx *gorm.DB, account *models.Account) error
	AddUser(ctx context.Context, tx *gorm.DB, account *models.Account, user *models.User) error
	DeleteUser(ctx context.Context, tx *gorm.DB, account *models.Account, user *models.User) error
	FindAllByUserID(ctx context.Context, userId uuid.UUID) ([]*models.Account, error)
	FindAllByUserIDPagination(ctx context.Context, userId uuid.UUID, input *dto.PaginationInput) (*dto.PaginationResult[*models.Account], error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	FindByIDWithPreload(ctx context.Context, id uuid.UUID) (*models.Account, error)
}

type accountRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepositoryImpl{
		db,
	}
}

func (a accountRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, account *models.Account) error {
	return tx.WithContext(ctx).Create(account).Error
}

func (a accountRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, account *models.Account) error {
	return tx.WithContext(ctx).Updates(account).Error
}

func (a accountRepositoryImpl) AddUser(ctx context.Context, tx *gorm.DB, account *models.Account, user *models.User) error {
	if err := tx.WithContext(ctx).Model(account).Association("Users").Append(user); err != nil {
		return err
	}
	return nil
}

func (a accountRepositoryImpl) DeleteUser(ctx context.Context, tx *gorm.DB, account *models.Account, user *models.User) error {
	if err := tx.WithContext(ctx).Model(account).Association("Users").Delete(user); err != nil {
		return err
	}
	return nil
}

func (a accountRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, account *models.Account) error {
	result := tx.WithContext(ctx).Delete(account)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (a accountRepositoryImpl) FindAllByUserID(ctx context.Context, userId uuid.UUID) ([]*models.Account, error) {
	var accounts []*models.Account
	err := a.DB.WithContext(ctx).
		Joins("JOIN user_accounts ua ON ua.account_id = accounts.id").
		Where("ua.user_id = ?", userId).
		Preload("Users").
		Find(&accounts).Error

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a accountRepositoryImpl) FindAllByUserIDPagination(ctx context.Context, userId uuid.UUID, input *dto.PaginationInput) (*dto.PaginationResult[*models.Account], error) {
	return utils.PaginateWithScope[*models.Account](a.DB, &models.Account{}, int(input.Page), int(input.Limit), func(tx *gorm.DB) *gorm.DB {
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

		tx = tx.WithContext(ctx).
			Joins("JOIN user_accounts ua ON ua.account_id = accounts.id").
			Where("ua.user_id = ?", userId).
			Preload("Users")

		if input.Search != "" {
			tx = tx.Where("name ILIKE ?", "%"+input.Search+"%")
		}

		return tx.Order(input.SortBy + " " + input.SortOrder)
	})
}

func (a accountRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	var account *models.Account
	err := a.DB.WithContext(ctx).First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a accountRepositoryImpl) FindByIDWithPreload(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	var account *models.Account
	err := a.DB.WithContext(ctx).
		Preload("Users").
		Preload("Transactions").
		First(&account, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return account, nil
}
