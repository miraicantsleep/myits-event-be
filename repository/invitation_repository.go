package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity"
	"gorm.io/gorm"
)

type (
	InvitationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error)
		GetInvitationByID(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) (entity.Invitation, error)
		GetInvitationByEvent(ctx context.Context, tx *gorm.DB, eventID uuid.UUID) ([]dto.InvitationResponse, error)
		GetAllUserInvitations(ctx context.Context, tx *gorm.DB) ([]entity.Invitation, error)
		Update(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error)
		Delete(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) error
		CheckInvitationExist(ctx context.Context, tx *gorm.DB, eventID uuid.UUID, userID ...uuid.UUID) (bool, error)
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
	if err := tx.WithContext(ctx).Preload("Event").Preload("Users").Where("id = ?", invitationID).First(&invitation).Error; err != nil {
		return entity.Invitation{}, err
	}

	return invitation, nil
}

func (r *invitationRepository) GetInvitationByEvent(ctx context.Context, tx *gorm.DB, eventID uuid.UUID) ([]dto.InvitationResponse, error) {
	if tx == nil {
		tx = r.db
	}
	var resp []dto.InvitationResponse
	err := tx.WithContext(ctx).
		Table("invitations").
		Select(
			"invitations.id AS id",
			"events.name AS event_name",
			"users.name AS name",
			"user_invitation.invited_at AS invited_at",
			"user_invitation.rsvp_status AS rsvp_status",
		).
		Joins("JOIN events ON events.id = invitations.event_id").
		Joins("JOIN user_invitation ON user_invitation.invitation_id = invitations.id").
		Joins("JOIN users ON users.id = user_invitation.user_id").
		Where("invitations.event_id = ?", eventID).
		Scan(&resp).Error
	return resp, err // maksa dikit
}

func (r *invitationRepository) GetAllUserInvitations(ctx context.Context, tx *gorm.DB) ([]entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	var invitations []entity.Invitation
	if err := tx.WithContext(ctx).
		Preload("Users").
		Preload("Event").
		Find(&invitations).Error; err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *invitationRepository) Create(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&invitation).Error; err != nil {
		return entity.Invitation{}, err
	}

	if err := tx.WithContext(ctx).Preload("Users").Preload("Event").First(&invitation, "id = ?", invitation.ID).Error; err != nil {
		return entity.Invitation{}, err
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

func (r *invitationRepository) CheckInvitationExist(ctx context.Context, tx *gorm.DB, eventID uuid.UUID, userIDs ...uuid.UUID) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	err := tx.Model(&entity.Invitation{}).
		Joins("INNER JOIN user_invitation ui ON ui.invitation_id = invitations.id").
		Where("invitations.event_id = ? AND ui.user_id IN ?", eventID, userIDs).
		Count(&count).Error
	return count > 0, err
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
