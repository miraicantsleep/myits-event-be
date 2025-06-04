package repository

import (
	"context"
	"log"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		Register(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		GetAllUserWithPagination(
			ctx context.Context,
			tx *gorm.DB,
			req dto.PaginationRequest,
		) (dto.GetAllUserRepositoryResponse, error)
		GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
		CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error)
		Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		Delete(ctx context.Context, tx *gorm.DB, userId string) error
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Register(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
    if tx == nil {
        tx = r.db
    }

    // Log input user sebelum Insert
    log.Printf("[Repo.Register] Input: %+v\n", user)

    // Misal pakai Gorm Create() bawaan:
    if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
        // Log error detail sebelum return
        log.Printf("[Repo.Register] Gagal insert ke DB: %v\n", err)
        return entity.User{}, err
    }

    log.Printf("[Repo.Register] Berhasil insert, ID: %s\n", user.ID)
    return user, nil
}

func (r *userRepository) GetAllUserWithPagination(
	ctx context.Context,
	tx *gorm.DB,
	req dto.PaginationRequest,
) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entity.User
	var count int64

	req.Default()

	dbQuery := tx.WithContext(ctx).Model(&entity.User{})
	if req.Search != "" {
		likePattern := "%" + req.Search + "%"
		dbQuery = dbQuery.Where("name LIKE ?", likePattern)
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	offset := (req.Page - 1) * req.PerPage
	if err := dbQuery.
		Limit(req.PerPage).
		Offset(offset).
		Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	err := tx.WithContext(ctx).
		First(&user, "id = ?", userId).
		Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	err := tx.WithContext(ctx).
		First(&user, "email = ?", email).
		Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	err := tx.WithContext(ctx).
		First(&user, "email = ?", email).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	query := `
		UPDATE users
		SET name = ?, email = ?, password = ?, updated_at = ?
		WHERE id = ?
	`
	if err := tx.WithContext(ctx).Exec(query,
		user.Name,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, userId string) error {
	if tx == nil {
		tx = r.db
	}

	query := "DELETE FROM users WHERE id = ?"
	if err := tx.WithContext(ctx).Exec(query, userId).Error; err != nil {
		return err
	}

	return nil
}
