package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	RSVPStatusAccepted = "accepted"
	RSVPStatusDeclined = "declined"
	RSVPStatusPending  = "pending"
)

type Invitation struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	EventID    uuid.UUID  `gorm:"type:uuid;not null" json:"event_id"`
	InvitedAt  time.Time  `gorm:"type:timestamp;not null;default:current_timestamp" json:"invited_at"`
	RsvpStatus string     `gorm:"type:varchar(20);not null;default:'pending'" json:"rsvp_status"`
	RsvpAt     *time.Time `gorm:"type:timestamp"            json:"rsvp_at,omitempty"`
	AttendedAt *time.Time `gorm:"type:timestamp"            json:"attended_at,omitempty"`
	Users      []User     `gorm:"many2many:user_invitation" json:"users"`
}
