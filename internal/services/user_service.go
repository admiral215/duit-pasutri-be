package services

import (
	"context"
	"duit-pasutri-be/internal/auth"
	"duit-pasutri-be/internal/dto"
	"duit-pasutri-be/internal/middleware"
	"duit-pasutri-be/internal/models"
	"duit-pasutri-be/internal/repositories"
	"duit-pasutri-be/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetAll(ctx context.Context) ([]*dto.UserDTO, error)
	GetAllPaginated(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedUser, error)
	GetByID(ctx context.Context, id string) (*dto.UserDTO, error)
	GetByEmail(ctx context.Context, email string) (*dto.UserDTO, error)
	Me(ctx context.Context) (*dto.UserDTO, error)
	Register(ctx context.Context, register *dto.Register) (*dto.UserDTO, error)
	Login(ctx context.Context, input dto.Login) (*dto.LoginResponse, error)
	Update(ctx context.Context, user *dto.UserUpdate) (*dto.UserDTO, error)
}

type userServiceImpl struct {
	UserRepository repositories.UserRepository
	Tx             *gorm.DB
}

func NewUserService(userRepository repositories.UserRepository, DB *gorm.DB) UserService {
	return &userServiceImpl{
		userRepository,
		DB,
	}
}

func (u userServiceImpl) GetAll(ctx context.Context) ([]*dto.UserDTO, error) {
	users, err := u.UserRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	usersDTO := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		usersDTO[i] = user.ToDTO()
	}

	return usersDTO, nil
}

func (u userServiceImpl) GetAllPaginated(ctx context.Context, input *dto.PaginationInput) (*dto.PaginatedUser, error) {
	result, err := u.UserRepository.FindAllPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	users := result.Data
	usersDTO := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		usersDTO[i] = user.ToDTO()
	}

	return &dto.PaginatedUser{
		Page:       result.Page,
		Limit:      result.Limit,
		TotalRows:  result.TotalRows,
		TotalPages: result.TotalPages,
		Data:       usersDTO,
	}, nil
}

func (u userServiceImpl) GetByID(ctx context.Context, id string) (*dto.UserDTO, error) {
	uid := uuid.MustParse(id)
	user, err := u.UserRepository.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}

func (u userServiceImpl) GetByEmail(ctx context.Context, email string) (*dto.UserDTO, error) {
	user, err := u.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}

func (u userServiceImpl) Me(ctx context.Context) (*dto.UserDTO, error) {
	userCtx, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.UserRepository.FindByID(ctx, userCtx.UserID)
	if err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}

func (u userServiceImpl) Login(ctx context.Context, input dto.Login) (*dto.LoginResponse, error) {
	user, err := u.UserRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, err
	}

	token, err := auth.CreateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, err
	}

	resp := &dto.LoginResponse{
		Token: token,
		User:  user.ToDTO(),
	}

	return resp, nil
}

func (u userServiceImpl) Register(ctx context.Context, register *dto.Register) (*dto.UserDTO, error) {
	errs := utils.ValidateStruct(register)
	if errs != nil {
		return nil, gqlerror.Errorf("Validation errors: %s", errs)
	}

	if register.Password != register.ConfirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	register.Password = string(hashPassword)

	newUser := &models.User{
		Name:     register.Name,
		Password: register.Password,
		Role:     models.RoleUser,
		Email:    register.Email,
	}

	var user *models.User
	err = u.Tx.Transaction(func(tx *gorm.DB) error {
		user, err = u.UserRepository.Create(ctx, tx, newUser)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user.ToDTO(), nil
}

func (u userServiceImpl) Update(ctx context.Context, user *dto.UserUpdate) (*dto.UserDTO, error) {
	if user.ID == "" {
		return nil, fmt.Errorf("user id is empty")
	}

	id, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, err
	}

	err = u.Tx.Transaction(func(tx *gorm.DB) error {
		currentUser, errTx := u.UserRepository.FindByID(ctx, id)
		if errTx != nil {
			return errTx
		}

		currentUser.Name = user.Name
		currentUser.Email = user.Email

		errTx = u.UserRepository.Update(ctx, tx, currentUser)
		if errTx != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedUser, err := u.UserRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedUser.ToDTO(), nil
}
