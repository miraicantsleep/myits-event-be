package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	InvitationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, invitation []entity.Invitation) ([]entity.Invitation, error)
		GetInvitationByID(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) (entity.Invitation, error)
		GetInvitationByEventID(ctx context.Context, tx *gorm.DB, eventID uuid.UUID) ([]entity.Invitation, error)
		GetAllInvitations(ctx context.Context, tx *gorm.DB) ([]entity.Invitation, error)
		Update(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error)
		Delete(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) error
		CheckInvitationExist(ctx context.Context, tx *gorm.DB, eventID uuid.UUID, userID uuid.UUID) (bool, error)
		CheckInvitationRSVPStatus(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID, rsvpStatus string) (bool, error)
	}

	invitationRepository struct {
		db *gorm.DB
	}
)

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{
		db: db,
	}
}

func (r *invitationRepository) GetInvitationByID(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) (entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	var invitation entity.Invitation
	if err := tx.WithContext(ctx).Where("id = ?", invitationID).First(&invitation).Error; err != nil {
		return entity.Invitation{}, err
	}

	return invitation, nil
}

func (r *invitationRepository) GetInvitationByEventID(ctx context.Context, tx *gorm.DB, eventID uuid.UUID) ([]entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	var invitations []entity.Invitation
	if err := tx.WithContext(ctx).Where("event_id = ?", eventID).Find(&invitations).Error; err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *invitationRepository) GetAllInvitations(ctx context.Context, tx *gorm.DB) ([]entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	var invitations []entity.Invitation
	if err := tx.WithContext(ctx).Find(&invitations).Error; err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *invitationRepository) Create(ctx context.Context, tx *gorm.DB, invitation []entity.Invitation) ([]entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&invitation).Error; err != nil {
		return []entity.Invitation{}, err
	}

	return invitation, nil
}

func (r *invitationRepository) Update(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&invitation).Error; err != nil {
		return entity.Invitation{}, err
	}

	return invitation, nil
}

func (r *invitationRepository) Delete(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.Invitation{}, "id = ?", invitationID).Error; err != nil {
		return err
	}

	return nil
}

func (r *invitationRepository) CheckInvitationExist(ctx context.Context, tx *gorm.DB, eventID uuid.UUID, userID uuid.UUID) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	var count int64
	if err := tx.WithContext(ctx).
		Table("invitations").
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *invitationRepository) CheckInvitationRSVPStatus(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID, rsvpStatus string) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	var count int64
	if err := tx.WithContext(ctx).
		Table("invitations").
		Where("id = ? AND rsvp_status = ?", invitationID, rsvpStatus).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
