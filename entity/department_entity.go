package entity

import (
	"github.com/google/uuid"
)

type Department struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name    string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Faculty string    `gorm:"type:varchar(100);not null" json:"faculty" validate:"required,min=2,max=100"`

	Timestamp
}
