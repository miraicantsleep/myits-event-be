package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookingRequest struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EventID     uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	Event       Event     `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"event"`
	RequestedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp" json:"requested_at"`
	Status      string    `gorm:"type:booking_status;not null;default:'pending'" json:"status" validate:"required,oneof=pending approved rejected"`
	Rooms       []Room    `gorm:"many2many:booking_request_room" json:"rooms"`
	Timestamp
}
