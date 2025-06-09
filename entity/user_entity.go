package entity

import (
	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/helpers"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleDepartemen UserRole = "departemen"
	RoleOrmawa     UserRole = "ormawa"
	RoleAdmin      UserRole = "admin"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string    `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=8"`
	Role     UserRole  `gorm:"type:user_role;not null;default:'user'" json:"role" validate:"required,oneof=user departemen ormawa admin"`

	// relationships
	Events      []Event      `gorm:"foreignKey:Created_By;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"events,omitempty"`
	Department  *Department  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"department,omitempty"`
	Invitations []Invitation `gorm:"many2many:user_invitation" json:"invitations"`
	Timestamp
}

// BeforeCreate hook to hash password and set defaults
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	// Hash password
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	// Ensure UUID is set
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Set default role if not specified
	if u.Role == "" {
		u.Role = "user"
	}

	return nil
}

// BeforeUpdate hook to handle password updates
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	// Only hash password if it has been changed
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
