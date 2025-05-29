package repository

import (
	"context"

	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	DepartmentRepository interface {
		Create(ctx context.Context, tx *gorm.DB, department entity.Department) (entity.Department, error)
		GetAllDepartmentWithPagination(
			ctx context.Context,
			tx *gorm.DB,
			req dto.PaginationRequest,
		) (dto.GetAllDepartmentRepositoryResponse, error)
		GetDepartmentById(ctx context.Context, tx *gorm.DB, departmentId string) (entity.Department, error)
		Update(ctx context.Context, tx *gorm.DB, department entity.Department) (entity.Department, error)
		Delete(ctx context.Context, tx *gorm.DB, departmentId string) error
	}

	departmentRepository struct {
		db *gorm.DB
	}
)

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (r *departmentRepository) Create(ctx context.Context, tx *gorm.DB, department entity.Department) (entity.Department, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&department).Error; err != nil {
		return entity.Department{}, err
	}

	return department, nil
}

func (r *departmentRepository) GetAllDepartmentWithPagination(
	ctx context.Context,
	tx *gorm.DB,
	req dto.PaginationRequest,
) (dto.GetAllDepartmentRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var departments []entity.Department
	var err error
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.Department{})
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllDepartmentRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&departments).Error; err != nil {
		return dto.GetAllDepartmentRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllDepartmentRepositoryResponse{
		Departments: departments,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *departmentRepository) GetDepartmentById(ctx context.Context, tx *gorm.DB, departmentId string) (entity.Department, error) {
	if tx == nil {
		tx = r.db
	}

	var department entity.Department
	if err := tx.WithContext(ctx).Where("id = ?", departmentId).Take(&department).Error; err != nil {
		return entity.Department{}, err
	}

	return department, nil
}

func (r *departmentRepository) Update(ctx context.Context, tx *gorm.DB, department entity.Department) (entity.Department, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&department).Error; err != nil {
		return entity.Department{}, err
	}

	return department, nil
}

func (r *departmentRepository) Delete(ctx context.Context, tx *gorm.DB, departmentId string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.Department{}, "id = ?", departmentId).Error; err != nil {
		return err
	}

	return nil
}
