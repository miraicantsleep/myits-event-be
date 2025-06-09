package repository

import (
	"context"

	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	RoomRepository interface {
		Create(ctx context.Context, room entity.Room) (entity.Room, error)
		GetRoomByID(ctx context.Context, id string) (entity.Room, error)
		GetRoomByName(ctx context.Context, name string) (entity.Room, error)
		GetAllRoom(ctx context.Context) ([]entity.Room, error)
		Update(ctx context.Context, id string, room entity.Room) (entity.Room, error)
		Delete(ctx context.Context, id string) error
	}
	roomRepository struct {
		db *gorm.DB
	}
)

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) Create(ctx context.Context, room entity.Room) (entity.Room, error) {
	tx := r.db
	if tx == nil {
		return entity.Room{}, dto.ErrCreateRoom
	}

	var existingRoom entity.Room
	if err := tx.WithContext(ctx).Where("name = ?", room.Name).First(&existingRoom).Error; err == nil {
		return entity.Room{}, dto.ErrRoomAlreadyExists
	}

	if err := tx.WithContext(ctx).Create(&room).Error; err != nil {
		return entity.Room{}, err
	}
	return room, nil
}

func (r *roomRepository) GetRoomByID(ctx context.Context, id string) (entity.Room, error) {
	return entity.Room{}, nil
}

func (r *roomRepository) GetRoomByName(ctx context.Context, name string) (entity.Room, error) {
	//
	// Implementation for getting a room by name
	return entity.Room{}, nil
}

func (r *roomRepository) GetAllRoom(ctx context.Context) ([]entity.Room, error) {
	//
	// Implementation for getting all rooms
	return nil, nil
}

func (r *roomRepository) Update(ctx context.Context, id string, room entity.Room) (entity.Room, error) {
	//
	// Implementation for updating a room
	return entity.Room{}, nil
}

func (r *roomRepository) Delete(ctx context.Context, id string) error {
	//
	// Implementation for deleting a room
	return nil
}
