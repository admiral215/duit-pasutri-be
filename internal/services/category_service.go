package services

import (
	"context"
	"duit-pasutri-be/internal/dto"
	"duit-pasutri-be/internal/models"
	"duit-pasutri-be/internal/repositories"
	"duit-pasutri-be/internal/utils"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

type CategoryService interface {
	GetAllPaginated(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedCategory, error)
	Create(ctx context.Context, input *dto.CategoryUpsert) (*dto.CategoryDTO, error)
	Update(ctx context.Context, input *dto.CategoryUpsert) (*dto.CategoryDTO, error)
	Delete(ctx context.Context, id string) (*dto.CategoryDTO, error)
}

type categoryServiceImpl struct {
	DB                 *gorm.DB
	CategoryRepository repositories.CategoryRepository
}

func NewCategoryService(DB *gorm.DB, categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryServiceImpl{
		DB:                 DB,
		CategoryRepository: categoryRepo,
	}
}

func (c categoryServiceImpl) Create(ctx context.Context, input *dto.CategoryUpsert) (*dto.CategoryDTO, error) {
	errs := utils.ValidateStruct(input)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	newCategory := &models.Category{
		Name: input.Name,
		Type: models.TransactionType(input.Type),
	}

	err := c.DB.Transaction(func(tx *gorm.DB) error {
		err := c.CategoryRepository.Create(ctx, tx, newCategory)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newCategory.ToDTO(), nil
}

func (c categoryServiceImpl) Update(ctx context.Context, input *dto.CategoryUpsert) (*dto.CategoryDTO, error) {
	errs := utils.ValidateStruct(input)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors : %s", errs)
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	category, err := c.CategoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = input.Name
	category.Type = models.TransactionType(input.Type)
	err = c.DB.Transaction(func(tx *gorm.DB) error {
		errTx := c.CategoryRepository.Update(ctx, tx, category)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return category.ToDTO(), nil
}

func (c categoryServiceImpl) GetAllPaginated(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedCategory, error) {
	result, err := c.CategoryRepository.FindAllPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	categories := result.Data

	dataDto := make([]*dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dataDto = append(dataDto, category.ToDTO())
	}

	return &dto.PaginatedCategory{
		Page:       result.Page,
		Limit:      result.Limit,
		TotalRows:  result.TotalRows,
		TotalPages: result.TotalPages,
		Data:       dataDto,
	}, nil
}

func (c categoryServiceImpl) Delete(ctx context.Context, id string) (*dto.CategoryDTO, error) {
	uid := uuid.MustParse(id)

	deletedCategory, err := c.CategoryRepository.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	err = c.DB.Transaction(func(tx *gorm.DB) error {
		errTx := c.CategoryRepository.Delete(ctx, tx, uid)
		if errTx != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return deletedCategory.ToDTO(), nil
}
