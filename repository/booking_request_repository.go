package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	BookingRequestRepository interface {
		CreateBookingRequest(ctx context.Context, tx *gorm.DB, bookingRequest *entity.BookingRequest) error
		GetBookingRequestByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entity.BookingRequest, error)
		GetAllBookingRequests(ctx context.Context, tx *gorm.DB) ([]entity.BookingRequest, error)
		UpdateBookingRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error
	}

	bookingRequestRepository struct {
		db *gorm.DB // Main DB connection, transactions will be passed or handled within methods
	}
)

func NewBookingRequestRepository(db *gorm.DB) BookingRequestRepository {
	return &bookingRequestRepository{db: db}
}

func (r *bookingRequestRepository) CreateBookingRequest(ctx context.Context, tx *gorm.DB, br *entity.BookingRequest) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Create(br).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRequestRepository) GetBookingRequestByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entity.BookingRequest, error) {
	var bookingRequest entity.BookingRequest
	db := r.db
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Preload("Rooms").Preload("Event").First(&bookingRequest, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bookingRequest, nil
}

func (r *bookingRequestRepository) GetAllBookingRequests(ctx context.Context, tx *gorm.DB) ([]entity.BookingRequest, error) {
	var bookingRequests []entity.BookingRequest
	db := r.db
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Preload("Rooms").Preload("Event").Find(&bookingRequests).Error
	if err != nil {
		return nil, err
	}
	return bookingRequests, nil
}

func (r *bookingRequestRepository) UpdateBookingRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Model(&entity.BookingRequest{}).Where("id = ?", id).Update("status", status).Error
}
