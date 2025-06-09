package repository

import (
	"context"
	"errors"

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
		GetDepartmentByUserId(ctx context.Context, tx *gorm.DB, userId string) (*entity.Department, error)
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

	req.Default()
	offset := (req.Page - 1) * req.PerPage

	// Hitung total count
	var count int64
	countQuery := `
		SELECT COUNT(*) 
		FROM departments d
		LEFT JOIN users u ON d.user_id = u.id
		WHERE d.name LIKE ? AND d.deleted_at IS NULL AND u.deleted_at IS NULL
	`
	if err := tx.Raw(countQuery, "%"+req.Search+"%").Scan(&count).Error; err != nil {
		return dto.GetAllDepartmentRepositoryResponse{}, err
	}

	// Query data dengan join
	var departmentResponses []dto.DepartmentResponse
	dataQuery := `
		SELECT d.id, d.name, d.faculty, u.email
		FROM departments d
		LEFT JOIN users u ON d.user_id = u.id
		WHERE d.name LIKE ? AND d.deleted_at IS NULL AND u.deleted_at IS NULL
		LIMIT ? OFFSET ?
	`
	if err := tx.Raw(dataQuery, "%"+req.Search+"%", req.PerPage, offset).Scan(&departmentResponses).Error; err != nil {
		return dto.GetAllDepartmentRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllDepartmentRepositoryResponse{
		Departments: departmentResponses,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
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
func (r *departmentRepository) GetDepartmentByUserId(ctx context.Context, tx *gorm.DB, userId string) (*entity.Department, error) {
	if tx == nil {
		tx = r.db
	}
	var department entity.Department

	if err := tx.WithContext(ctx).
		First(&department, "admin_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrDepartmentNotFound
		}
		return nil, err
	}

	return &department, nil
}
