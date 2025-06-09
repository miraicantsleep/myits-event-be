package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	// Success messages
	MESSAGE_SUCCESS_CREATE_INVITATION    = "Success create invitation"
	MESSAGE_SUCCESS_GET_INVITATION_BY_ID = "Success get invitation by id"
	MESSAGE_SUCCESS_GET_ALL_INVITATIONS  = "Success get all invitations"
	MESSAGE_SUCCESS_UPDATE_INVITATION    = "Success update invitation"
	MESSAGE_SUCCESS_DELETE_INVITATION    = "Success delete invitation"

	// Failed messages
	MESSAGE_FAILED_CREATE_INVITATION    = "Failed create invitation"
	MESSAGE_FAILED_GET_INVITATION_BY_ID = "Failed get invitation by id"
	MESSAGE_FAILED_GET_ALL_INVITATIONS  = "Failed get all invitations"
	MESSAGE_FAILED_UPDATE_INVITATION    = "Failed update invitation"
	MESSAGE_FAILED_DELETE_INVITATION    = "Failed delete invitation"
)

var (
	ErrCreateInvitation            = errors.New("failed to create invitation")
	ErrGetInvitationByID           = errors.New("failed to get invitation by id")
	ErrGetInvitationByEventID      = errors.New("failed to get invitation by event id")
	ErrGetAllInvitations           = errors.New("failed to get all invitations")
	ErrUpdateInvitation            = errors.New("failed to update invitation")
	ErrInvitationNotFound          = errors.New("invitation not found")
	ErrDeleteInvitation            = errors.New("failed to delete invitation")
	ErrInvitationAlreadyExists     = errors.New("invitation already exists")
	ErrInvitationInvalidRSVPStatus = errors.New("invalid RSVP status, must be one of accepted, declined, pending")
)

type CreateInvitationRequest struct {
	EventID uuid.UUID   `json:"event_id" binding:"required"`
	UserIDs []uuid.UUID `json:"user_ids" binding:"required,min=1"`
}

type UpdateInvitationRequest struct {
	RSVPStatus string     `json:"rsvp_status" binding:"required,oneof=accepted declined pending"`
	RsvpAt     *time.Time `json:"rsvp_at,omitempty"`
}

type InvitationResponse struct {
	ID         uuid.UUID   `json:"id"`
	EventID    uuid.UUID   `json:"event_id"`
	InvitedAt  time.Time   `json:"invited_at"`
	RSVPStatus string      `json:"rsvp_status"`
	RsvpAt     *time.Time  `json:"rsvp_at,omitempty"`
	AttendedAt *time.Time  `json:"attended_at,omitempty"`
	UserIDs    []uuid.UUID `json:"user_ids"`
}
