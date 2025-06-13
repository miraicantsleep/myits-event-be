package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
	"gorm.io/gorm"
)

type (
	BookingRequestService interface {
		CreateBookingRequest(ctx context.Context, req dto.BookingRequestCreateRequest) (dto.BookingRequestResponse, error)
		GetBookingRequestByID(ctx context.Context, id string) (dto.BookingRequestResponse, error)
		GetAllBookingRequests(ctx context.Context) ([]dto.BookingDetailResponse, error)
		UpdateBookingRequest(ctx context.Context, id string, req dto.BookingRequestUpdateRequest, role string) (dto.BookingRequestResponse, error)
		DeleteBookingRequest(ctx context.Context, id string) error
		ApproveBookingRequest(ctx context.Context, id string) error
		RejectBookingRequest(ctx context.Context, id string) error
	}

	bookingRequestService struct {
		bookingRequestRepo repository.BookingRequestRepository
		roomRepo           repository.RoomRepository
		eventRepo          repository.EventRepository
		jwtService         JWTService
		db                 *gorm.DB
	}
)

func NewBookingRequestService(
	bookingRequestRepo repository.BookingRequestRepository,
	roomRepo repository.RoomRepository,
	eventRepo repository.EventRepository,
	jwtService JWTService,
	db *gorm.DB,
) BookingRequestService {
	return &bookingRequestService{
		bookingRequestRepo: bookingRequestRepo,
		roomRepo:           roomRepo,
		eventRepo:          eventRepo,
		jwtService:         jwtService,
		db:                 db,
	}
}

func (s *bookingRequestService) CreateBookingRequest(ctx context.Context, req dto.BookingRequestCreateRequest) (dto.BookingRequestResponse, error) {
	var response dto.BookingRequestResponse
	var roomsForBooking []entity.Room
	var roomResponses []dto.RoomResponse

	tx := s.db.Begin()
	if tx.Error != nil {
		return response, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, roomID := range req.RoomIDs {
		room, err := s.roomRepo.GetRoomByID(ctx, roomID.String())
		if err != nil {
			tx.Rollback()
			return response, err
		}
		roomsForBooking = append(roomsForBooking, room)
		roomResponses = append(roomResponses, dto.RoomResponse{
			ID:       room.ID.String(),
			Name:     room.Name,
			Capacity: room.Capacity,
		})
	}

	event, err := s.eventRepo.GetEventById(ctx, tx, req.EventID.String())
	if err != nil {
		tx.Rollback()
		return response, err
	}

	bookingRequest := entity.BookingRequest{
		EventID: req.EventID,
		Rooms:   roomsForBooking,
		Status:  "pending",
	}

	err = s.bookingRequestRepo.CreateBookingRequest(ctx, tx, &bookingRequest)
	if err != nil {
		tx.Rollback()
		return response, err
	}

	if err := tx.Commit().Error; err != nil {
		return response, err
	}

	response = dto.BookingRequestResponse{
		ID:          bookingRequest.ID,
		EventID:     bookingRequest.EventID,
		EventName:   event.Name,
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
		return response, err
	}

	bookingRequest, err := s.bookingRequestRepo.GetBookingRequestByID(ctx, nil, bookingRequestID)
	if err != nil {
		return response, err
	}

	var roomResponses []dto.RoomResponse
	for _, room := range bookingRequest.Rooms {
		roomResponses = append(roomResponses, dto.RoomResponse{
			ID:       room.ID.String(),
			Name:     room.Name,
			Capacity: room.Capacity,
		})
	}

	response = dto.BookingRequestResponse{
		ID:          bookingRequest.ID,
		EventID:     bookingRequest.EventID,
		EventName:   bookingRequest.Event.Name,
		RequestedAt: bookingRequest.RequestedAt.Format(time.RFC3339),
		Status:      bookingRequest.Status,
		Rooms:       roomResponses,
	}
	return response, nil
}

func (s *bookingRequestService) GetAllBookingRequests(ctx context.Context) ([]dto.BookingDetailResponse, error) {
	flatBookingData, err := s.bookingRequestRepo.GetAllBookingRequests(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 2. Use a map to group the flat data into the nested structure.
	bookingMap := make(map[uuid.UUID]*dto.BookingDetailResponse)

	for _, record := range flatBookingData {
		if _, exists := bookingMap[record.BookingID]; !exists {
			// Create the parent booking object
			bookingMap[record.BookingID] = &dto.BookingDetailResponse{
				BookingID:     record.BookingID,
				BookingStatus: record.BookingStatus,
				EventID:       record.EventID,
				EventName:     record.EventName,
				RequestedBy:   record.RequestedBy,
				Rooms:         []dto.RoomInfo{},
			}
		}

		// Add the room from the current row into the booking's "Rooms" slice
		if record.RoomID != uuid.Nil {
			room := dto.RoomInfo{
				RoomID:   record.RoomID,
				RoomName: record.RoomName,
			}
			bookingMap[record.BookingID].Rooms = append(bookingMap[record.BookingID].Rooms, room)
		}
	}

	// 3. Convert the map to a slice for the final response.
	finalResponse := make([]dto.BookingDetailResponse, 0, len(bookingMap))
	for _, booking := range bookingMap {
		finalResponse = append(finalResponse, *booking)
	}

	return finalResponse, nil
}

func (s *bookingRequestService) UpdateBookingRequest(ctx context.Context, id string, req dto.BookingRequestUpdateRequest, role string) (dto.BookingRequestResponse, error) {
	var response dto.BookingRequestResponse
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return response, err
	}

	if role == "ormawa" && req.Status != "" {
		return response, errors.New("ormawa users are not permitted to change the booking status")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	br, err := s.bookingRequestRepo.GetBookingRequestByID(ctx, tx, bookingRequestID)
	if err != nil {
		tx.Rollback()
		return response, err
	}

	if req.Status != "" {
		br.Status = req.Status
	}

	if len(req.RoomIDs) > 0 {
		var newRooms []entity.Room
		for _, roomID := range req.RoomIDs {
			room, err := s.roomRepo.GetRoomByID(ctx, roomID.String())
			if err != nil {
				tx.Rollback()
				return response, err
			}
			newRooms = append(newRooms, room)
		}
		br.Rooms = newRooms
	}

	if err := s.bookingRequestRepo.UpdateBookingRequest(ctx, tx, br); err != nil {
		tx.Rollback()
		return response, err
	}

	if err := tx.Commit().Error; err != nil {
		return response, err
	}

	// Refetch to get updated associations correctly for the response
	updatedBr, err := s.bookingRequestRepo.GetBookingRequestByID(ctx, nil, bookingRequestID)
	if err != nil {
		return response, err
	}

	var roomResponses []dto.RoomResponse
	for _, room := range updatedBr.Rooms {
		roomResponses = append(roomResponses, dto.RoomResponse{
			ID:       room.ID.String(),
			Name:     room.Name,
			Capacity: room.Capacity,
		})
	}

	response = dto.BookingRequestResponse{
		ID:          updatedBr.ID,
		EventID:     updatedBr.EventID,
		EventName:   updatedBr.Event.Name,
		RequestedAt: updatedBr.RequestedAt.Format(time.RFC3339),
		Status:      updatedBr.Status,
		Rooms:       roomResponses,
	}

	return response, nil
}

func (s *bookingRequestService) DeleteBookingRequest(ctx context.Context, id string) error {
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = s.bookingRequestRepo.GetBookingRequestByID(ctx, nil, bookingRequestID)
	if err != nil {
		return err
	}

	return s.bookingRequestRepo.DeleteBookingRequest(ctx, nil, bookingRequestID)
}

func (s *bookingRequestService) ApproveBookingRequest(ctx context.Context, id string) error {
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.bookingRequestRepo.UpdateBookingRequestStatus(ctx, nil, bookingRequestID, "approved")
}

func (s *bookingRequestService) RejectBookingRequest(ctx context.Context, id string) error {
	bookingRequestID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.bookingRequestRepo.UpdateBookingRequestStatus(ctx, nil, bookingRequestID, "rejected")
}
