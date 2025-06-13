package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"github.com/miraicantsleep/myits-event-be/repository"
	"gorm.io/gorm"
)

type (
	RoomService interface {
		Create(ctx context.Context, req dto.RoomCreateRequest, userId string, role string) (dto.RoomResponse, error)
		GetRoomByID(ctx context.Context, id string) (dto.RoomResponse, error)
		GetRoomByName(ctx context.Context, name string) (dto.RoomResponse, error)
		// get all room without pagination
		GetAllRoom(ctx context.Context) ([]dto.RoomResponse, error)
		Update(ctx context.Context, id string, req dto.RoomUpdateRequest) (dto.RoomResponse, error)
		Delete(ctx context.Context, id string) error
	}
	roomService struct {
		roomRepository       repository.RoomRepository
		departmentRepository repository.DepartmentRepository
		jwtService           JWTService
		db                   *gorm.DB
	}
)

func NewRoomService(roomRepository repository.RoomRepository, jwtService JWTService, db *gorm.DB, departmentRepository repository.DepartmentRepository) RoomService {
	return &roomService{
		roomRepository:       roomRepository,
		departmentRepository: departmentRepository,
		jwtService:           jwtService,
		db:                   db,
	}
}

func (s *roomService) Create(ctx context.Context, req dto.RoomCreateRequest, userId string, role string) (dto.RoomResponse, error) {
	roomEntity := entity.Room{
		Name:     req.Name,
		Capacity: req.Capacity,
	}
	if role == "departemen" {
		departmentId, err := s.departmentRepository.GetDepartmentByUserId(ctx, nil, userId)
		if err != nil {
			return dto.RoomResponse{}, err
		}
		roomEntity.DepartmentID = departmentId.ID
	} else if role == "admin" {
		departmentID, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			return dto.RoomResponse{}, err
		}
		roomEntity.DepartmentID = departmentID
	}

	result, err := s.roomRepository.Create(ctx, roomEntity)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	// get the department name from the department repository
	department, err := s.departmentRepository.GetDepartmentById(ctx, nil, result.DepartmentID.String())
	if err != nil {
		return dto.RoomResponse{}, err
	}

	response := dto.RoomResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		Department:   department.Name,
		DepartmentID: result.DepartmentID.String(),
		Capacity:     result.Capacity,
	}
	return response, nil
}

func (s *roomService) GetRoomByID(ctx context.Context, id string) (dto.RoomResponse, error) {
	if id == "" {
		return dto.RoomResponse{}, errors.New("room ID is required")
	}
	log.Println("GetRoomByID called with ID:", id)
	result, err := s.roomRepository.GetRoomByID(ctx, id)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	// get the department name from the department repository
	department, err := s.departmentRepository.GetDepartmentById(ctx, nil, result.DepartmentID.String())
	if err != nil {
		return dto.RoomResponse{}, err
	}

	response := dto.RoomResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		Department:   department.Name,
		DepartmentID: result.DepartmentID.String(),
		Capacity:     result.Capacity,
	}
	return response, nil
}

func (s *roomService) GetAllRoom(ctx context.Context) ([]dto.RoomResponse, error) {
	result, err := s.roomRepository.GetAllRoom(ctx)
	if err != nil {
		return nil, err
	}
	var response []dto.RoomResponse
	for _, room := range result {
		response = append(response, dto.RoomResponse{
			ID:           room.ID,
			Name:         room.Name,
			Department:   room.Department,
			DepartmentID: room.DepartmentID,
			Capacity:     room.Capacity,
		})
	}
	return response, nil
}

func (s *roomService) Update(ctx context.Context, id string, req dto.RoomUpdateRequest) (dto.RoomResponse, error) {
	roomEntity := entity.Room{
		Name:         req.Name,
		DepartmentID: uuid.MustParse(req.DepartmentID),
		Capacity:     req.Capacity,
	}
	result, err := s.roomRepository.Update(ctx, id, roomEntity)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	// get the department name from the department repository
	department, err := s.departmentRepository.GetDepartmentById(ctx, nil, result.DepartmentID.String())
	if err != nil {
		return dto.RoomResponse{}, err
	}

	response := dto.RoomResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		Department:   department.Name,
		DepartmentID: result.DepartmentID.String(),
		Capacity:     result.Capacity,
	}
	return response, nil
}

func (s *roomService) Delete(ctx context.Context, id string) error {
	err := s.roomRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *roomService) GetRoomByName(ctx context.Context, name string) (dto.RoomResponse, error) {
	result, err := s.roomRepository.GetRoomByName(ctx, name)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	// get the department name from the department repository
	department, err := s.departmentRepository.GetDepartmentById(ctx, nil, result.DepartmentID.String())
	if err != nil {
		return dto.RoomResponse{}, err
	}

	response := dto.RoomResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		Department:   department.Name,
		DepartmentID: result.DepartmentID.String(),
		Capacity:     result.Capacity,
	}
	return response, nil
}
