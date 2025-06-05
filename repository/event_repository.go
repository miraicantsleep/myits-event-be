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

	return event, nil
}

func (r *eventRepository) GetAllEventWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllEventRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var events []entity.Event
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.Event{})
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllEventRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&events).Error; err != nil {
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
			Created_By:  event.Created_By.String(),
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
	if err := tx.WithContext(ctx).Where("id = ?", eventId).First(&event).Error; err != nil {
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
