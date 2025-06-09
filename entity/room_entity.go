package entity

import (
	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	DepartmentID uuid.UUID        `gorm:"type:uuid;not null" json:"department_id"`
	Department   Department       `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"department"`
	Name         string           `gorm:"type:varchar(255);not null" json:"name"`
	Capacity     int              `gorm:"not null" json:"capacity"`
	Bookings     []BookingRequest `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"bookings,omitempty"`
	Timestamp
}
