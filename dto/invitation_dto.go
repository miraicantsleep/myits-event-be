package dto

import (
	"errors"
)

const (
	// Success messages
	MESSAGE_SUCCESS_CREATE_INVITATION          = "Success create invitation"
	MESSAGE_SUCCESS_GET_INVITATION_BY_ID       = "Success get invitation by id"
	MESSAGE_SUCCESS_GET_ALL_INVITATIONS        = "Success get all invitations"
	MESSAGE_SUCCESS_UPDATE_INVITATION          = "Success update invitation"
	MESSAGE_SUCCESS_DELETE_INVITATION          = "Success delete invitation"
	MESSAGE_SUCCESS_GET_INVITATION_BY_EVENT_ID = "Success get invitation by event id"

	// Failed messages
	MESSAGE_FAILED_CREATE_INVITATION          = "Failed create invitation"
	MESSAGE_FAILED_GET_INVITATION_BY_ID       = "Failed get invitation by id"
	MESSAGE_FAILED_GET_ALL_INVITATIONS        = "Failed get all invitations"
	MESSAGE_FAILED_UPDATE_INVITATION          = "Failed update invitation"
	MESSAGE_FAILED_DELETE_INVITATION          = "Failed delete invitation"
	MESSAGE_FAILED_GET_INVITATION_BY_EVENT_ID = "Failed get invitation by event id"
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
	EventID string   `json:"event_id" binding:"required"`
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
}

type CreateInvitationResponse struct {
	EventName  string   `json:"event_name"`
	Names      []string `json:"names,omitempty"`
	InvitedAt  string   `json:"invited_at"`
	RSVPStatus string   `json:"rsvp_status"`
}

type UpdateInvitationRequest struct {
	RSVPStatus string `json:"rsvp_status" binding:"required,oneof=accepted declined pending"`
	RsvpAt     string `json:"rsvp_at,omitempty"`
	AttendedAt string `json:"attended_at,omitempty"`
}

type InvitationResponse struct {
	ID         string `json:"id"`
	EventName  string `json:"event_name"`
	Name       string `json:"name,omitempty"`
	InvitedAt  string `json:"invited_at"`
	RSVPStatus string `json:"rsvp_status"`
	RsvpAt     string `json:"rsvp_at,omitempty"`
	AttendedAt string `json:"attended_at,omitempty"`
}
