package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RSVPStatusAccepted = "accepted"
	RSVPStatusDeclined = "declined"
	RSVPStatusPending  = "pending"
)

type Invitation struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	EventID uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	// Relationships
	Event Event  `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"event,omitempty"`
	Users []User `gorm:"many2many:user_invitation;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"users,omitempty"`
}

type UserInvitation struct {
	UserID       uuid.UUID  `gorm:"primaryKey"`
	InvitationID uuid.UUID  `gorm:"primaryKey"`
	QRCode       string     `gorm:"type:varchar(255);uniqueIndex" json:"qr_code,omitempty"`
	InvitedAt    time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"invited_at"`
	RSVPStatus   string     `gorm:"type:rsvp_status;not null;default:'pending'" json:"rsvp_status" validate:"required,oneof=accepted declined pending"`
	RsvpAt       *time.Time `gorm:"type:timestamp;default:null" json:"rsvp_at,omitempty"`
	AttendedAt   *time.Time `gorm:"type:timestamp;default:null" json:"attended_at,omitempty"`
}

func (ui *UserInvitation) BeforeCreate(tx *gorm.DB) (err error) {
	if ui.QRCode == "" { // Generate QR code only if not already set (e.g., during seeding or testing)
		ui.QRCode = uuid.New().String()
	}
	return
}

func (UserInvitation) TableName() string { return "user_invitation" }
