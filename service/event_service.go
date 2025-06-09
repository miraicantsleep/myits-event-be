package service

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
)

type (
	EventService interface {
		Create(ctx context.Context, req dto.EventCreateRequest, userId string) (dto.EventResponse, error)
		GetAllEventWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.EventPaginationResponse, error)
		GetEventById(ctx context.Context, eventId string) (dto.EventResponse, error)
		Update(ctx context.Context, req dto.EventUpdateRequest, eventId string) (dto.EventResponse, error)
		Delete(ctx context.Context, eventId string) error
	}
	eventService struct {
		eventRepo  repository.EventRepository
		jwtService JWTService
		db         *gorm.DB
	}
)

func NewEventService(
	eventRepo repository.EventRepository,
	jwtService JWTService,
	db *gorm.DB,
) EventService {
	return &eventService{
		eventRepo:  eventRepo,
		jwtService: jwtService,
		db:         db,
	}
}

func (s *eventService) Create(ctx context.Context, req dto.EventCreateRequest, userId string) (dto.EventResponse, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return dto.EventResponse{}, err
	}

	startTime, err := time.Parse(time.RFC3339, req.Start_Time)
	if err != nil {
		return dto.EventResponse{}, err
	}

	endTime, err := time.Parse(time.RFC3339, req.End_Time)
	if err != nil {
		return dto.EventResponse{}, err
	}

	// check if event with the same name already exists
	exists, _ := s.eventRepo.CheckEventExist(ctx, nil, req.Name)
	if exists {
		return dto.EventResponse{}, errors.New("event with the same name already exists")
	}

	event := entity.Event{
		Name:        req.Name,
		Description: req.Description,
		Start_Time:  startTime,
		End_Time:    endTime,
		Event_Type:  req.Event_Type,
		Created_By:  id,
	}

	eventReg, err := s.eventRepo.Create(ctx, nil, event)
	if err != nil {
		return dto.EventResponse{}, errors.New(err.Error())
	}

	return dto.EventResponse{
		ID:          eventReg.ID.String(),
		Name:        eventReg.Name,
		Description: eventReg.Description,
		Start_Time:  eventReg.Start_Time.Format(time.RFC3339),
		End_Time:    eventReg.End_Time.Format(time.RFC3339),
		Created_By:  eventReg.Creator_Name,
		Event_Type:  eventReg.Event_Type,
	}, nil
}

func (s *eventService) GetAllEventWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.EventPaginationResponse, error) {
	EventsWithPagination, err := s.eventRepo.GetAllEventWithPagination(ctx, nil, req)
	if err != nil {
		return dto.EventPaginationResponse{}, dto.ErrGetAllEvent
	}

	var eventResponses []dto.EventResponse
	for _, event := range EventsWithPagination.Events {
		eventResponses = append(eventResponses, dto.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Start_Time:  event.Start_Time,
			End_Time:    event.End_Time,
			Created_By:  event.Created_By,
			Event_Type:  event.Event_Type,
		})
	}

	return dto.EventPaginationResponse{
		Data: eventResponses,
		PaginationResponse: dto.PaginationResponse{
			Page:    EventsWithPagination.PaginationResponse.Page,
			PerPage: EventsWithPagination.PaginationResponse.PerPage,
			MaxPage: EventsWithPagination.PaginationResponse.MaxPage,
			Count:   EventsWithPagination.PaginationResponse.Count,
		},
	}, nil
}

func (s *eventService) GetEventById(ctx context.Context, eventId string) (dto.EventResponse, error) {
	event, err := s.eventRepo.GetEventById(ctx, nil, eventId)
	if err != nil {
		return dto.EventResponse{}, dto.ErrGetEventById
	}

	return dto.EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Description: event.Description,
		Start_Time:  event.Start_Time.Format(time.RFC3339),
		End_Time:    event.End_Time.Format(time.RFC3339),
		Created_By:  event.Creator_Name,
		Event_Type:  event.Event_Type,
	}, nil
}
func (s *eventService) Update(ctx context.Context, req dto.EventUpdateRequest, eventId string) (dto.EventResponse, error) {
	id, err := uuid.Parse(eventId)
	log.Println("eventId:", eventId, "parsed ID:", id)
	if err != nil {
		return dto.EventResponse{}, err
	}

	event, err := s.eventRepo.GetEventById(ctx, nil, id.String())
	if err != nil {
		return dto.EventResponse{}, dto.ErrEventNotFound
	}

	if req.Name != "" {
		event.Name = req.Name
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Start_Time != "" {
		startTime, err := time.Parse(time.RFC3339, req.Start_Time)
		if err != nil {
			return dto.EventResponse{}, err
		}
		event.Start_Time = startTime
	}
	if req.End_Time != "" {
		endTime, err := time.Parse(time.RFC3339, req.End_Time)
		if err != nil {
			return dto.EventResponse{}, err
		}
		event.End_Time = endTime
	}
	if req.Event_Type != "" {
		event.Event_Type = req.Event_Type
	}

	// check if event with the same name already exists
	exists, _ := s.eventRepo.CheckEventExist(ctx, nil, event.Name)
	if exists {
		return dto.EventResponse{}, errors.New("event with the same name already exists")
	}

	updatedEvent, err := s.eventRepo.Update(ctx, nil, event)
	if err != nil {
		return dto.EventResponse{}, dto.ErrUpdateEvent
	}
	return dto.EventResponse{
		ID:          updatedEvent.ID.String(),
		Name:        updatedEvent.Name,
		Description: updatedEvent.Description,
		Start_Time:  updatedEvent.Start_Time.Format(time.RFC3339),
		End_Time:    updatedEvent.End_Time.Format(time.RFC3339),
		Created_By:  updatedEvent.Creator_Name,
		Event_Type:  updatedEvent.Event_Type,
	}, nil
}

func (s *eventService) Delete(ctx context.Context, eventId string) error {
	id, err := uuid.Parse(eventId)
	if err != nil {
		return err
	}

	event, err := s.eventRepo.GetEventById(ctx, nil, id.String())
	if err != nil {
		return dto.ErrGetEventById
	}

	err = s.eventRepo.Delete(ctx, nil, event.ID.String())
	if err != nil {
		return dto.ErrDeleteEvent
	}

	return nil
}
