package entity

import (
	"github.com/google/uuid"
)

type RoomBookingRequest struct {
	RoomID           uuid.UUID      `gorm:"type:uuid;not null" json:"room_id"`
	Room             Room           `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	BookingRequestID uuid.UUID      `gorm:"type:uuid;not null" json:"booking_request_id"`
	BookingRequest   BookingRequest `gorm:"foreignKey:BookingRequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Timestamp
}
