package service

import (
	"context"
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
		jwtService     JWTService
		db             *gorm.DB
	}
)

func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
	jwtService JWTService,
	db *gorm.DB,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
		jwtService:     jwtService,
		db:             db,
	}
}

func (s *departmentService) Create(ctx context.Context, req dto.DepartmentCreateRequest) (dto.DepartmentResponse, error) {
	department := entity.Department{
		Name:    req.Name,
		Faculty: req.Faculty,
	}

	departmentReg, err := s.departmentRepo.Create(ctx, nil, department)
	if err != nil {
		return dto.DepartmentResponse{}, dto.ErrCreateDepartment
	}

	return dto.DepartmentResponse{
		ID:      departmentReg.ID.String(),
		Name:    departmentReg.Name,
		Faculty: departmentReg.Faculty,
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

	var datas []dto.DepartmentResponse
	for _, department := range dataWithPaginate.Departments {
		data := dto.DepartmentResponse{
			ID:      department.ID.String(),
			Name:    department.Name,
			Faculty: department.Faculty,
		}

		datas = append(datas, data)
	}

	return dto.DepartmentPaginationResponse{
		Data: datas,
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

	return dto.DepartmentResponse{
		ID:      department.ID.String(),
		Name:    department.Name,
		Faculty: department.Faculty,
	}, nil
}

func (s *departmentService) Update(ctx context.Context, req dto.DepartmentUpdateRequest, departmentId string) (
	dto.DepartmentUpdateResponse,
	error,
) {
	department, err := s.departmentRepo.GetDepartmentById(ctx, nil, departmentId)
	if err != nil {
		return dto.DepartmentUpdateResponse{}, dto.ErrDepartmentNotFound
	}

	data := entity.Department{
		ID:      department.ID,
		Name:    req.Name,
		Faculty: department.Faculty,
	}

	departmentUpdate, err := s.departmentRepo.Update(ctx, nil, data)
	if err != nil {
		return dto.DepartmentUpdateResponse{}, dto.ErrUpdateDepartment
	}

	return dto.DepartmentUpdateResponse{
		ID:      departmentUpdate.ID.String(),
		Name:    departmentUpdate.Name,
		Faculty: department.Faculty,
	}, nil
}

func (s *departmentService) Delete(ctx context.Context, departmentId string) error {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	department, err := s.departmentRepo.GetDepartmentById(ctx, nil, departmentId)
	if err != nil {
		return dto.ErrDepartmentNotFound
	}

	err = s.departmentRepo.Delete(ctx, nil, department.ID.String())
	if err != nil {
		return dto.ErrDeleteDepartment
	}

	return nil
}
