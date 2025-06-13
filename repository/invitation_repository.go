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
		GetInvitationByID(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) ([]dto.InvitationDetailResponse, error)
		GetInvitationByEvent(ctx context.Context, tx *gorm.DB, eventID uuid.UUID) ([]dto.InvitationResponse, error)
		GetInvitationByUserId(ctx context.Context, tx *gorm.DB, userID uuid.UUID) ([]dto.InvitationResponse, error)
		GetAllUserInvitations(ctx context.Context, tx *gorm.DB) ([]entity.Invitation, error)
		Update(ctx context.Context, tx *gorm.DB, invitation entity.Invitation) (entity.Invitation, error)
		Delete(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) error
		CheckInvitationExist(ctx context.Context, tx *gorm.DB, eventID uuid.UUID, userID ...uuid.UUID) (bool, error)
		CheckInvitationRSVPStatus(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID, rsvpStatus string) (bool, error)
		GetUserInvitationByQRCode(ctx context.Context, tx *gorm.DB, qrCode string) (entity.UserInvitation, error)
		UpdateUserInvitation(ctx context.Context, tx *gorm.DB, userInvitation entity.UserInvitation) (entity.UserInvitation, error)
		GetUserInvitation(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID, userID uuid.UUID) (entity.UserInvitation, error)
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

func (r *invitationRepository) GetInvitationByID(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID) ([]dto.InvitationDetailResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var results []dto.InvitationDetailResponse
	sql := `
		SELECT
			i.id AS invitation_id,
			e.id AS event_id,
			e.name AS event_name,
			u.id AS user_id,
			u.name AS user_name,
			u.email AS user_email,
			ui.invited_at,
			ui.rsvp_status,
			ui.rsvp_at,
			ui.attended_at,
			ui.qr_code,
			creator.name AS creator_name
		FROM
			invitations i
		JOIN
			user_invitation ui ON i.id = ui.invitation_id
		JOIN
			users u ON ui.user_id = u.id
		JOIN
			events e ON i.event_id = e.id
		JOIN
			users creator ON e.created_by = creator.id
		WHERE
			i.id = ?;
	`

	if err := tx.WithContext(ctx).Raw(sql, invitationID).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
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

	// Fetch the created invitation with its related Event and Users using Joins
	var createdInvitation entity.Invitation
	if err := tx.WithContext(ctx).
		Table("invitations").
		Select("invitations.*, events.*").
		Joins("LEFT JOIN events ON events.id = invitations.event_id").
		Where("invitations.id = ?", invitation.ID).
		First(&createdInvitation).Error; err != nil {
		return entity.Invitation{}, err
	}

	// Now fetch the associated users for the invitation
	var users []entity.User
	if err := tx.WithContext(ctx).
		Table("users").
		Joins("JOIN user_invitation ON user_invitation.user_id = users.id").
		Where("user_invitation.invitation_id = ?", invitation.ID).
		Find(&users).Error; err != nil {
		return entity.Invitation{}, err
	}
	createdInvitation.Users = users

	return createdInvitation, nil
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

// GetUserInvitationByQRCode retrieves a UserInvitation by its QRCode
func (r *invitationRepository) GetUserInvitationByQRCode(ctx context.Context, tx *gorm.DB, qrCode string) (entity.UserInvitation, error) {
	if tx == nil {
		tx = r.db
	}
	var userInvitation entity.UserInvitation
	// It's important to query the user_invitation table directly.
	// We also need to join with Users and Invitations (and Events through Invitations)
	// if we want to return more comprehensive information, but for now, just the UserInvitation record is fine.
	// However, the prompt is to update UserInvitation, so preloading User and Event is not strictly necessary for the update itself.
	// Let's keep it simple first and fetch only UserInvitation.
	if err := tx.WithContext(ctx).Where("qr_code = ?", qrCode).First(&userInvitation).Error; err != nil {
		return entity.UserInvitation{}, err
	}
	return userInvitation, nil
}

func (r *invitationRepository) GetUserInvitation(ctx context.Context, tx *gorm.DB, invitationID uuid.UUID, userID uuid.UUID) (entity.UserInvitation, error) {
	if tx == nil {
		tx = r.db
	}
	var userInvitation entity.UserInvitation
	if err := tx.WithContext(ctx).Where("invitation_id = ? AND user_id = ?", invitationID, userID).First(&userInvitation).Error; err != nil {
		return entity.UserInvitation{}, err
	}
	return userInvitation, nil
}

// UpdateUserInvitation updates an existing UserInvitation record
func (r *invitationRepository) UpdateUserInvitation(ctx context.Context, tx *gorm.DB, userInvitation entity.UserInvitation) (entity.UserInvitation, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&userInvitation).Error; err != nil {
		return entity.UserInvitation{}, err
	}
	return userInvitation, nil
}

// GetInvitationByUserId retrieves all invitations for a specific user
func (r *invitationRepository) GetInvitationByUserId(ctx context.Context, tx *gorm.DB, userID uuid.UUID) ([]dto.InvitationResponse, error) {
	if tx == nil {
		tx = r.db
	}
	var resp []dto.InvitationResponse
	err := tx.WithContext(ctx).
		Table("full_invitation_details").
		Where("user_id = ?", userID).
		Scan(&resp).Error

	return resp, err
}
