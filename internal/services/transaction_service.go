package services

import (
	"context"
	"duit-pasutri-be/internal/dto"
	"duit-pasutri-be/internal/middleware"
	"duit-pasutri-be/internal/models"
	"duit-pasutri-be/internal/repositories"
	"duit-pasutri-be/internal/utils"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

type TransactionService interface {
	Create(ctx context.Context, input *dto.TransactionUpsert) (*dto.TransactionDTO, error)
	Update(ctx context.Context, input *dto.TransactionUpsert) (*dto.TransactionDTO, error)
	Delete(ctx context.Context, id string) (*dto.TransactionDTO, error)
	GetMyTransactions(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedTransaction, error)
}

type transactionServiceImpl struct {
	DB                    *gorm.DB
	TransactionRepository repositories.TransactionRepository
	AccountRepository     repositories.AccountRepository
	CategoryRepository    repositories.CategoryRepository
	UserRepository        repositories.UserRepository
}

func NewTransactionService(db *gorm.DB,
	transactionRepository repositories.TransactionRepository,
	accountRepository repositories.AccountRepository,
	categoryRepository repositories.CategoryRepository,
	userRepository repositories.UserRepository,
) TransactionService {
	return &transactionServiceImpl{
		DB:                    db,
		TransactionRepository: transactionRepository,
		AccountRepository:     accountRepository,
		CategoryRepository:    categoryRepository,
		UserRepository:        userRepository,
	}
}

func (t transactionServiceImpl) Create(ctx context.Context, input *dto.TransactionUpsert) (*dto.TransactionDTO, error) {
	errs := utils.ValidateStruct(input)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	account, err := t.AccountRepository.FindByID(ctx, uuid.MustParse(input.AccountID))
	if err != nil {
		return nil, err
	}

	category, err := t.CategoryRepository.FindByID(ctx, uuid.MustParse(input.CategoryID))
	if err != nil {
		return nil, err
	}

	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := t.UserRepository.FindByID(ctx, userCtx.UserID)
	if err != nil {
		return nil, err
	}

	newTransaction := &models.Transaction{
		Amount:      input.Amount,
		Description: input.Description,
		Type:        models.TransactionType(input.Type),

		Account:  account,
		Category: category,
		User:     user,
	}

	err = t.DB.Transaction(func(tx *gorm.DB) error {
		errTx := t.TransactionRepository.Create(ctx, tx, newTransaction)
		if errTx == nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newTransaction.ToDTO(), nil
}

func (t transactionServiceImpl) Update(ctx context.Context, input *dto.TransactionUpsert) (*dto.TransactionDTO, error) {
	errs := utils.ValidateStruct(input)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	transaction, err := t.TransactionRepository.FindByIDWithPreload(ctx, id)
	if err != nil {
		return nil, err
	}

	account, err := t.AccountRepository.FindByID(ctx, uuid.MustParse(input.AccountID))
	if err != nil {
		return nil, err
	}

	category, err := t.CategoryRepository.FindByID(ctx, uuid.MustParse(input.CategoryID))
	if err != nil {
		return nil, err
	}

	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != userCtx.UserID {
		return nil, gqlerror.Errorf("You are not the user on this transaction")
	}

	transaction.Description = input.Description
	transaction.Amount = input.Amount
	transaction.Type = models.TransactionType(input.Type)
	transaction.Account = account
	transaction.Category = category
	err = t.DB.Transaction(func(tx *gorm.DB) error {
		errTx := t.TransactionRepository.Update(ctx, tx, transaction)
		if errTx == nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transaction.ToDTO(), nil
}

func (t transactionServiceImpl) Delete(ctx context.Context, id string) (*dto.TransactionDTO, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	transaction, err := t.TransactionRepository.FindByIDWithPreload(ctx, uid)
	if err != nil {
		return nil, err
	}

	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userCtx.UserID {
		return nil, gqlerror.Errorf("You are not the user on this transaction")
	}

	err = t.DB.Transaction(func(tx *gorm.DB) error {
		errTx := t.TransactionRepository.Delete(ctx, tx, transaction)
		if errTx == nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transaction.ToDTO(), nil
}

func (t transactionServiceImpl) GetMyTransactions(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedTransaction, error) {
	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := t.TransactionRepository.FindAllPaginationByUserID(ctx, input, userCtx.UserID)
	if err != nil {
		return nil, err
	}

	transactions := result.Data
	transactionsDto := make([]*dto.TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		transactionsDto[i] = transaction.ToDTO()
	}

	return &dto.PaginatedTransaction{
		Page:       result.Page,
		Limit:      result.Limit,
		TotalRows:  result.TotalRows,
		TotalPages: result.TotalPages,
		Data:       transactionsDto,
	}, nil
}
