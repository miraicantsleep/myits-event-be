package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
	"gorm.io/gorm"
)

type BookingRequestService interface {
	CreateBookingRequest(ctx context.Context, req dto.BookingRequestCreateRequest) (dto.BookingRequestResponse, error)
	GetBookingRequestByID(ctx context.Context, id string) (dto.BookingRequestResponse, error)
	GetAllBookingRequests(ctx context.Context) ([]dto.BookingRequestResponse, error)
	ApproveBookingRequest(ctx context.Context, id string) error
	RejectBookingRequest(ctx context.Context, id string) error
}

type bookingRequestService struct {
	bookingRequestRepo repository.BookingRequestRepository
	roomRepo           repository.RoomRepository
	eventRepo          repository.EventRepository // Assuming an EventRepository exists
	db                 *gorm.DB                   // For transactions
}

func NewBookingRequestService(
	brRepo repository.BookingRequestRepository,
	roomRepo repository.RoomRepository,
	eventRepo repository.EventRepository, // Add eventRepo to constructor
	db *gorm.DB,
) BookingRequestService {
	return &bookingRequestService{
		bookingRequestRepo: brRepo,
		roomRepo:           roomRepo,
		eventRepo:          eventRepo,
		db:                 db,
	}
}

func (s *bookingRequestService) CreateBookingRequest(ctx context.Context, req dto.BookingRequestCreateRequest) (dto.BookingRequestResponse, error) {
	var response dto.BookingRequestResponse
	var roomsForBooking []entity.Room
	var roomResponses []dto.RoomResponse

	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return response, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Fetch and validate rooms
	for _, roomID := range req.RoomIDs {
		room, err := s.roomRepo.GetRoomByID(ctx, roomID.String()) // Assuming GetRoomByID takes string
		if err != nil {
			tx.Rollback()
			return response, err // Consider specific error for room not found
		}
		roomsForBooking = append(roomsForBooking, room)
		// Assuming RoomResponse can be created from entity.Room
		// This might need adjustment based on actual RoomResponse structure and how department name is fetched
		roomResponses = append(roomResponses, dto.RoomResponse{
			ID: room.ID.String(),
			Name: room.Name,
			// Department: room.Department.Name, // This requires Department to be preloaded or fetched
			Capacity: room.Capacity,
		})
	}

	// 2. Fetch event details (optional, for response enrichment)
	event, err := s.eventRepo.GetEventByID(ctx, tx, req.EventID.String()) // Assuming GetEventByID takes string and tx
	if err != nil {
		// Decide if this is a critical error. If event name in response is nice-to-have, maybe log and continue.
		// For now, let's treat it as important for the response.
		tx.Rollback()
		return response, err
	}


	bookingRequest := entity.BookingRequest{
		EventID: req.EventID,
		Rooms:   roomsForBooking,
		Status:  "pending", // Default status
		// RequestedAt will be set by DB default or Timestamp struct
	}

	err = s.bookingRequestRepo.CreateBookingRequest(ctx, tx, &bookingRequest)
	if err != nil {
		tx.Rollback()
		return response, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return response, err
	}

	response = dto.BookingRequestResponse{
		ID:          bookingRequest.ID,
		EventID:     bookingRequest.EventID,
		EventName:   event.Name, // Populate event name
		RequestedAt: bookingRequest.RequestedAt.Format(time.RFC3339),
		Status:      bookingRequest.Status,
		Rooms:       roomResponses,
	}

	return response, nil
}

func (s *bookingRequestService) GetBookingRequestByID(ctx context.Context, id string) (dto.BookingRequestResponse, error) {
	var response dto.BookingRequestResponse
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return response, err // Invalid UUID format
	}

	bookingRequest, err := s.bookingRequestRepo.GetBookingRequestByID(ctx, nil, bookingRequestID) // Pass nil for tx if not needed
	if err != nil {
		return response, err
	}

	var roomResponses []dto.RoomResponse
	for _, room := range bookingRequest.Rooms {
		// This part might need adjustment for fetching Department.Name for each room
		// For simplicity, assuming Room entity has Department preloaded or RoomResponse doesn't strictly need it here
		roomResponses = append(roomResponses, dto.RoomResponse{
			ID: room.ID.String(),
			Name: room.Name,
			// Department: room.Department.Name, // This requires Department to be preloaded or fetched in GetBookingRequestByID
			Capacity: room.Capacity,
		})
	}

	response = dto.BookingRequestResponse{
		ID:          bookingRequest.ID,
		EventID:     bookingRequest.EventID,
		EventName:   bookingRequest.Event.Name, // Assumes Event is preloaded by GetBookingRequestByID
		RequestedAt: bookingRequest.RequestedAt.Format(time.RFC3339),
		Status:      bookingRequest.Status,
		Rooms:       roomResponses,
	}
	return response, nil
}

func (s *bookingRequestService) GetAllBookingRequests(ctx context.Context) ([]dto.BookingRequestResponse, error) {
	var responses []dto.BookingRequestResponse
	bookingRequests, err := s.bookingRequestRepo.GetAllBookingRequests(ctx, nil) // Pass nil for tx
	if err != nil {
		return nil, err
	}

	for _, br := range bookingRequests {
		var roomResponses []dto.RoomResponse
		for _, room := range br.Rooms {
			roomResponses = append(roomResponses, dto.RoomResponse{
				ID: room.ID.String(),
				Name: room.Name,
				// Department: room.Department.Name, // Requires preloading
				Capacity: room.Capacity,
			})
		}
		responses = append(responses, dto.BookingRequestResponse{
			ID:          br.ID,
			EventID:     br.EventID,
			EventName:   br.Event.Name, // Assumes Event is preloaded
			RequestedAt: br.RequestedAt.Format(time.RFC3339),
			Status:      br.Status,
			Rooms:       roomResponses,
		})
	}
	return responses, nil
}

func (s *bookingRequestService) ApproveBookingRequest(ctx context.Context, id string) error {
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return err // Invalid UUID format
	}
	// Potentially add validation here: check if request exists, current status, etc.
	return s.bookingRequestRepo.UpdateBookingRequestStatus(ctx, nil, bookingRequestID, "approved") // Pass nil for tx
}

func (s *bookingRequestService) RejectBookingRequest(ctx context.Context, id string) error {
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return err // Invalid UUID format
	}
	// Potentially add validation here
	return s.bookingRequestRepo.UpdateBookingRequestStatus(ctx, nil, bookingRequestID, "rejected") // Pass nil for tx
}
