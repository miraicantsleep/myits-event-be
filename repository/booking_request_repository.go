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
		UpdateBookingRequest(ctx context.Context, tx *gorm.DB, bookingRequest *entity.BookingRequest) error
		UpdateBookingRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error
		DeleteBookingRequest(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
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
	err := db.WithContext(ctx).
		Joins("Event").
		Joins("left join booking_request_room on booking_request_room.booking_request_id = booking_requests.id").
		Joins("left join rooms on rooms.id = booking_request_room.room_id").
		Where("booking_requests.id = ?", id).
		First(&bookingRequest).Error
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
	err := db.WithContext(ctx).
		Joins("Event").
		Joins("left join booking_request_room on booking_request_room.booking_request_id = booking_requests.id").
		Joins("left join rooms on rooms.id = booking_request_room.room_id").
		Find(&bookingRequests).Error
	if err != nil {
		return nil, err
	}
	return bookingRequests, nil
}

func (r *bookingRequestRepository) UpdateBookingRequest(ctx context.Context, tx *gorm.DB, br *entity.BookingRequest) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	// Use Session with FullSaveAssociations to update the entity and its many-to-many relationships (Rooms)
	return db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(br).Error
}

func (r *bookingRequestRepository) UpdateBookingRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Model(&entity.BookingRequest{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookingRequestRepository) DeleteBookingRequest(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	// GORM's default delete will be a soft delete due to the gorm.DeletedAt field in the Timestamp struct.
	// We should also clear the associations in the join table.
	br := entity.BookingRequest{ID: id}
	if err := db.WithContext(ctx).Model(&br).Association("Rooms").Clear(); err != nil {
		return err
	}

	return db.WithContext(ctx).Delete(&br).Error
}
