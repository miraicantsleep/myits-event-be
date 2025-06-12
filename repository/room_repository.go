package repository

import (
	"context"

	"github.com/google/uuid"
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
	if err := tx.WithContext(ctx).Where("name = ? AND department_id = ?", room.Name, room.DepartmentID).First(&existingRoom).Error; err == nil {
		return entity.Room{}, dto.ErrRoomAlreadyExists
	}

	if err := tx.WithContext(ctx).Create(&room).Error; err != nil {
		return entity.Room{}, err
	}
	return room, nil
}

func (r *roomRepository) GetRoomByID(ctx context.Context, id string) (entity.Room, error) {
	tx := r.db
	if tx == nil {
		return entity.Room{}, dto.ErrGetRoomByID
	}
	var room entity.Room
	if err := tx.WithContext(ctx).Preload("Department").First(&room, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Room{}, dto.ErrRoomNotFound
		}
		return entity.Room{}, err
	}
	return room, nil
}

func (r *roomRepository) GetRoomByName(ctx context.Context, name string) (entity.Room, error) {
	tx := r.db
	if tx == nil {
		return entity.Room{}, dto.ErrGetRoomByName
	}
	var room entity.Room
	if err := tx.WithContext(ctx).Preload("Department").Where("name = ?", name).First(&room).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Room{}, dto.ErrRoomNotFound
		}
		return entity.Room{}, err
	}
	return room, nil
}

func (r *roomRepository) GetAllRoom(ctx context.Context) ([]entity.Room, error) {
	tx := r.db
	if tx == nil {
		return nil, dto.ErrGetAllRoom
	}
	var rooms []entity.Room
	if err := tx.WithContext(ctx).Preload("Department").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) Update(ctx context.Context, id string, room entity.Room) (entity.Room, error) {
	tx := r.db
	if tx == nil {
		return entity.Room{}, dto.ErrUpdateRoom
	}

	var existingRoom entity.Room
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&existingRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Room{}, dto.ErrRoomNotFound
		}
		return entity.Room{}, err
	}

	room.ID, _ = uuid.Parse(id)

	if err := tx.WithContext(ctx).Model(&existingRoom).Updates(room).Error; err != nil {
		return entity.Room{}, err
	}

	if err := tx.WithContext(ctx).Preload("Department").Where("id = ?", id).First(&room).Error; err != nil {
		return entity.Room{}, err
	}

	return room, nil
}

func (r *roomRepository) Delete(ctx context.Context, id string) error {
	tx := r.db
	if tx == nil {
		return dto.ErrDeleteRoom
	}
	var room entity.Room
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&room).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.ErrRoomNotFound
		}
		return err
	}
	if err := tx.WithContext(ctx).Delete(&room).Error; err != nil {
		return err
	}
	return nil
}
