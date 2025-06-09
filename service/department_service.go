package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
)

type (
	DepartmentService interface {
		Create(ctx context.Context, req dto.DepartmentCreateRequest) (dto.DepartmentResponse, error)
		GetAllDepartmentWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.DepartmentPaginationResponse, error)
		GetDepartmentById(ctx context.Context, departmentId string) (dto.DepartmentResponse, error)
		Update(ctx context.Context, req dto.DepartmentUpdateRequest, departmentId string) (dto.DepartmentUpdateResponse, error)
		Delete(ctx context.Context, departmentId string) error
	}

	departmentService struct {
		departmentRepo repository.DepartmentRepository
		userRepo       repository.UserRepository
		jwtService     JWTService
		db             *gorm.DB
	}
)

func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
	userRepo repository.UserRepository,
	jwtService JWTService,
	db *gorm.DB,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
		userRepo:       userRepo,
		jwtService:     jwtService,
		db:             db,
	}
}

func (s *departmentService) Create(ctx context.Context, req dto.DepartmentCreateRequest) (dto.DepartmentResponse, error) {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	// Convert DTO to entity for user creation
	userEntity := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "departemen",
	}

	// Create user via repository (assuming you have userRepo)
	createdUser, err := s.userRepo.Register(ctx, tx, userEntity)
	if err != nil {
		return dto.DepartmentResponse{}, dto.ErrCreateUser
	}

	// Create department with user_id
	department := entity.Department{
		Name:    req.Name,
		Faculty: req.Faculty,
		UserID:  createdUser.ID,
	}

	departmentReg, err := s.departmentRepo.Create(ctx, tx, department)
	if err != nil {
		return dto.DepartmentResponse{}, errors.New(err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return dto.DepartmentResponse{}, err
	}

	return dto.DepartmentResponse{
		ID:      departmentReg.ID.String(),
		Name:    departmentReg.Name,
		Faculty: departmentReg.Faculty,
		Email:   createdUser.Email,
	}, nil
}

func (s *departmentService) GetAllDepartmentWithPagination(
	ctx context.Context,
	req dto.PaginationRequest,
) (dto.DepartmentPaginationResponse, error) {
	dataWithPaginate, err := s.departmentRepo.GetAllDepartmentWithPagination(ctx, nil, req)
	if err != nil {
		return dto.DepartmentPaginationResponse{}, err
	}

	return dto.DepartmentPaginationResponse{
		Data: dataWithPaginate.Departments,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *departmentService) GetDepartmentById(ctx context.Context, departmentId string) (dto.DepartmentResponse, error) {
	department, err := s.departmentRepo.GetDepartmentById(ctx, nil, departmentId)
	if err != nil {
		return dto.DepartmentResponse{}, dto.ErrGetDepartmentById
	}

	user, err := s.userRepo.GetUserById(ctx, nil, department.UserID.String())
	if err != nil {
		return dto.DepartmentResponse{}, dto.ErrGetUserById
	}

	return dto.DepartmentResponse{
		ID:      department.ID.String(),
		Name:    department.Name,
		Faculty: department.Faculty,
		Email:   user.Email,
	}, nil
}

func (s *departmentService) Update(ctx context.Context, req dto.DepartmentUpdateRequest, departmentId string) (
	dto.DepartmentUpdateResponse,
	error,
) {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	department, err := s.departmentRepo.GetDepartmentById(ctx, tx, departmentId)
	if err != nil {
		tx.Rollback()
		return dto.DepartmentUpdateResponse{}, dto.ErrDepartmentNotFound
	}

	if req.Email != "" {
		userUpdate := entity.User{
			ID:    department.UserID,
			Name:  req.Name,
			Email: req.Email,
		}

		if err := tx.WithContext(ctx).Model(&userUpdate).Select("name", "email").Updates(userUpdate).Error; err != nil {
			tx.Rollback()
			return dto.DepartmentUpdateResponse{}, err
		}
	}

	departmentData := entity.Department{
		ID:      department.ID,
		Name:    req.Name,
		Faculty: req.Faculty,
		UserID:  department.UserID,
	}

	departmentUpdate, err := s.departmentRepo.Update(ctx, tx, departmentData)
	if err != nil {
		tx.Rollback()
		return dto.DepartmentUpdateResponse{}, dto.ErrUpdateDepartment
	}

	if err := tx.Commit().Error; err != nil {
		return dto.DepartmentUpdateResponse{}, err
	}

	return dto.DepartmentUpdateResponse{
		ID:      departmentUpdate.ID.String(),
		Name:    departmentUpdate.Name,
		Faculty: departmentUpdate.Faculty,
		Email:   req.Email,
	}, nil
}

func (s *departmentService) Delete(ctx context.Context, departmentId string) error {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	department, err := s.departmentRepo.GetDepartmentById(ctx, tx, departmentId)
	if err != nil {
		return dto.ErrDepartmentNotFound
	}

	err = s.departmentRepo.Delete(ctx, tx, department.ID.String())
	if err != nil {
		return dto.ErrDeleteDepartment
	}

	err = s.userRepo.Delete(ctx, tx, department.UserID.String())
	if err != nil {
		return dto.ErrDeleteUser
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
