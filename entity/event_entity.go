package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventTypeOnline  = "online"
	EventTypeOffline = "offline"
)

type Event struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name              string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Description       string    `gorm:"type:text;not null" json:"description" validate:"required,min=10,max=500"`
	Start_Time        time.Time `gorm:"type:timestamp;not null" json:"start_time" validate:"required"`
	End_Time          time.Time `gorm:"type:timestamp;not null" json:"end_time" validate:"required"`
	Created_By        uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	Event_Type        string    `gorm:"type:event_type;not null;default:'offline'" json:"event_type" validate:"required,oneof=online offline"`
	DurationInMinutes int       `gorm:"type:integer;" json:"duration_in_minutes"`

	// Relationships
	Invitations []Invitation `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"invitations,omitempty"`
	// temp
	Creator_Name string `gorm:"->;column:creator_name" json:"creator_name,omitempty"`
	Timestamp
}
