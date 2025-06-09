package repository

import (
	"context"

	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	EventRepository interface {
		Create(ctx context.Context, tx *gorm.DB, event entity.Event) (entity.Event, error)
		GetAllEventWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllEventRepositoryResponse, error)
		GetEventById(ctx context.Context, tx *gorm.DB, eventId string) (entity.Event, error)
		Update(ctx context.Context, tx *gorm.DB, event entity.Event) (entity.Event, error)
		Delete(ctx context.Context, tx *gorm.DB, eventId string) error
		CheckEventExist(ctx context.Context, tx *gorm.DB, name string) (bool, error)
	}

	eventRepository struct {
		db *gorm.DB
	}
)

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, tx *gorm.DB, event entity.Event) (entity.Event, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&event).Error; err != nil {
		return entity.Event{}, err
	}

	return r.GetEventById(ctx, tx, event.ID.String())
}

func (r *eventRepository) GetAllEventWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllEventRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var events []entity.Event
	var count int64

	req.Default()

	baseQuery := tx.WithContext(ctx).Table("events").Joins("LEFT JOIN users ON events.created_by = users.id")
	if req.Search != "" {
		baseQuery = baseQuery.Where("events.name LIKE ?", "%"+req.Search+"%")
	}

	if err := baseQuery.Count(&count).Error; err != nil {
		return dto.GetAllEventRepositoryResponse{}, err
	}

	if err := baseQuery.
		Select("events.*, users.name as creator_name").
		Scopes(Paginate(req)).
		Find(&events).Error; err != nil {
		return dto.GetAllEventRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))

	// Convert []entity.Event to []dto.EventResponse
	eventResponses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = dto.EventResponse{
			ID:          event.ID.String(),
			Name:        event.Name,
			Description: event.Description,
			Start_Time:  event.Start_Time.String(),
			End_Time:    event.End_Time.String(),
			Created_By:  event.Creator_Name,
			Event_Type:  event.Event_Type,
		}
	}

	return dto.GetAllEventRepositoryResponse{
		Events: eventResponses,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *eventRepository) GetEventById(ctx context.Context, tx *gorm.DB, eventId string) (entity.Event, error) {
	if tx == nil {
		tx = r.db
	}

	var event entity.Event

	if err := tx.WithContext(ctx).
		Table("events").
		Select("events.*, users.name as creator_name").
		Joins("LEFT JOIN users ON events.created_by = users.id").
		Where("events.id = ?", eventId).
		First(&event).Error; err != nil {
		return entity.Event{}, err
	}

	return event, nil
}

func (r *eventRepository) Update(ctx context.Context, tx *gorm.DB, event entity.Event) (entity.Event, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&event).Error; err != nil {
		return entity.Event{}, err
	}

	return event, nil
}

func (r *eventRepository) Delete(ctx context.Context, tx *gorm.DB, eventId string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.Event{}, "id = ?", eventId).Error; err != nil {
		return err
	}

	return nil
}

func (r *eventRepository) CheckEventExist(ctx context.Context, tx *gorm.DB, name string) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	var count int64

	if err := tx.WithContext(ctx).
		Table("events").
		Where("events.name = ?", name).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
