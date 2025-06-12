package dto

import "github.com/google/uuid"

const (
	// Success
	MESSAGE_SUCCESS_CREATE_BOOKING_REQUEST   = "Success create booking request"
	MESSAGE_SUCCESS_GET_BOOKING_REQUEST      = "Success get booking request"
	MESSAGE_SUCCESS_GET_ALL_BOOKING_REQUESTS = "Success get all booking requests"
	MESSAGE_SUCCESS_UPDATE_BOOKING_REQUEST   = "Success update booking request"
	MESSAGE_SUCCESS_DELETE_BOOKING_REQUEST   = "Success delete booking request"
	MESSAGE_SUCCESS_APPROVE_BOOKING_REQUEST  = "Success approve booking request"
	MESSAGE_SUCCESS_REJECT_BOOKING_REQUEST   = "Success reject booking request"

	// Failed
	MESSAGE_FAILED_CREATE_BOOKING_REQUEST   = "Failed create booking request"
	MESSAGE_FAILED_GET_BOOKING_REQUEST      = "Failed get booking request"
	MESSAGE_FAILED_GET_ALL_BOOKING_REQUESTS = "Failed get all booking requests"
	MESSAGE_FAILED_UPDATE_BOOKING_REQUEST   = "Failed update booking request"
	MESSAGE_FAILED_DELETE_BOOKING_REQUEST   = "Failed delete booking request"
	MESSAGE_FAILED_APPROVE_BOOKING_REQUEST  = "Failed approve booking request"
	MESSAGE_FAILED_REJECT_BOOKING_REQUEST   = "Failed reject booking request"
)

type BookingRequestCreateRequest struct {
	EventID uuid.UUID   `json:"event_id" binding:"required"`
	RoomIDs []uuid.UUID `json:"room_ids" binding:"required,min=1"`
}

type BookingRequestUpdateRequest struct {
	RoomIDs []uuid.UUID `json:"room_ids" binding:"omitempty,min=1"`
	Status  string      `json:"status" binding:"omitempty,oneof=pending approved rejected"`
}

type BookingRequestResponse struct {
	ID          uuid.UUID      `json:"id"`
	EventID     uuid.UUID      `json:"event_id"`
	EventName   string         `json:"event_name"`
	RequestedAt string         `json:"requested_at"`
	Status      string         `json:"status"`
	Rooms       []RoomResponse `json:"rooms"`
}
