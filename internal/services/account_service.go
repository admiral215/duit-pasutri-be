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

type AccountService interface {
	Create(ctx context.Context, accountDto *dto.AccountUpsert) (*dto.AccountDTO, error)
	Update(ctx context.Context, accountDto *dto.AccountUpsert) (*dto.AccountDTO, error)
	AddUser(ctx context.Context, addUser *dto.AccountUser) (*dto.AccountDTO, error)
	DeleteUser(ctx context.Context, delUser *dto.AccountUser) (*dto.AccountDTO, error)
	GetMyAccounts(ctx context.Context) ([]*dto.AccountDTO, error)
	GetMyAccountsPaginated(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedAccount, error)
	GetAccountById(ctx context.Context, id string) (*dto.AccountDTO, error)
}

type accountServiceImpl struct {
	DB                *gorm.DB
	AccountRepository repositories.AccountRepository
	UserRepository    repositories.UserRepository
}

func NewAccountService(db *gorm.DB, accountRepository repositories.AccountRepository, userRepository repositories.UserRepository) AccountService {
	return &accountServiceImpl{
		db,
		accountRepository,
		userRepository,
	}
}

func (a accountServiceImpl) Create(ctx context.Context, accountDto *dto.AccountUpsert) (*dto.AccountDTO, error) {
	errs := utils.ValidateStruct(accountDto)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors %s", errs)
	}

	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := a.UserRepository.FindByID(ctx, userCtx.UserID)
	if err != nil {
		return nil, err
	}

	newAccount := &models.Account{
		Name:        accountDto.Name,
		Description: accountDto.Description,
	}

	err = a.DB.Transaction(func(tx *gorm.DB) error {
		errTx := a.AccountRepository.Create(ctx, tx, newAccount)
		if err != nil {
			return errTx
		}

		errTx = a.AccountRepository.AddUser(ctx, tx, newAccount, user)
		if err != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	createdAccount, err := a.AccountRepository.FindByID(ctx, newAccount.ID)
	if err != nil {
		return nil, err
	}

	return createdAccount.ToDTO(), nil
}

func (a accountServiceImpl) Update(ctx context.Context, accountDto *dto.AccountUpsert) (*dto.AccountDTO, error) {
	errs := utils.ValidateStruct(accountDto)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	if accountDto.ID == "" {
		return nil, gqlerror.Errorf("ID required")
	}

	account, err := a.AccountRepository.FindByID(ctx, uuid.MustParse(accountDto.ID))
	if err != nil {
		return nil, err
	}

	account.Name = accountDto.Name
	account.Description = accountDto.Description

	err = a.DB.Transaction(func(tx *gorm.DB) error {
		if errTx := a.AccountRepository.Update(ctx, tx, account); errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedAccount, err := a.AccountRepository.FindByID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return updatedAccount.ToDTO(), nil
}

func (a accountServiceImpl) AddUser(ctx context.Context, addUser *dto.AccountUser) (*dto.AccountDTO, error) {
	errs := utils.ValidateStruct(addUser)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	accountID := uuid.MustParse(addUser.AccountID)
	userID := uuid.MustParse(addUser.UserID)

	account, err := a.AccountRepository.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	user, err := a.UserRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = a.DB.Transaction(func(tx *gorm.DB) error {
		return a.AccountRepository.AddUser(ctx, tx, account, user)
	})

	if err != nil {
		return nil, err
	}

	updatedAccount, err := a.AccountRepository.FindByIDWithPreload(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	return updatedAccount.ToDTO(), nil
}

func (a accountServiceImpl) DeleteUser(ctx context.Context, delUser *dto.AccountUser) (*dto.AccountDTO, error) {
	errs := utils.ValidateStruct(delUser)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	accountID := uuid.MustParse(delUser.AccountID)
	userID := uuid.MustParse(delUser.UserID)

	account, err := a.AccountRepository.FindByIDWithPreload(ctx, accountID)
	if err != nil {
		return nil, err
	}

	err = a.DB.Transaction(func(tx *gorm.DB) error {
		for _, user := range account.Users {
			if user.ID == userID {
				errTx := a.AccountRepository.DeleteUser(ctx, tx, account, user)
				if errTx != nil {
					return errTx
				}
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	updatedAccount, err := a.AccountRepository.FindByIDWithPreload(ctx, account.ID)
	if err != nil {
		return nil, err
	}
	return updatedAccount.ToDTO(), nil
}

func (a accountServiceImpl) GetMyAccounts(ctx context.Context) ([]*dto.AccountDTO, error) {
	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	accounts, err := a.AccountRepository.FindAllByUserID(ctx, userCtx.UserID)
	if err != nil {
		return nil, err
	}

	accountDTOs := make([]*dto.AccountDTO, len(accounts))
	for i, account := range accounts {
		accountDTOs[i] = account.ToDTO()
	}
	return accountDTOs, nil
}

func (a accountServiceImpl) GetMyAccountsPaginated(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedAccount, error) {
	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := a.AccountRepository.FindAllByUserIDPagination(ctx, userCtx.UserID, input)
	if err != nil {
		return nil, err
	}

	accounts := result.Data
	accountDTOs := make([]*dto.AccountDTO, len(accounts))
	for i, account := range accounts {
		accountDTOs[i] = account.ToDTO()
	}

	return &dto.PaginatedAccount{
		Page:       result.Page,
		Limit:      result.Limit,
		TotalRows:  result.TotalRows,
		TotalPages: result.TotalPages,
		Data:       accountDTOs,
	}, nil
}

func (a accountServiceImpl) GetAccountById(ctx context.Context, id string) (*dto.AccountDTO, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	account, err := a.AccountRepository.FindByIDWithPreload(ctx, uid)
	if err != nil {
		return nil, err
	}
	return account.ToDTO(), nil
}
