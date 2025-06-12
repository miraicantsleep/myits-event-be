package dto

import "github.com/google/uuid"

type BookingRequestCreateRequest struct {
	EventID uuid.UUID   `json:"event_id" binding:"required"`
	RoomIDs []uuid.UUID `json:"room_ids" binding:"required,min=1"`
}

type BookingRequestResponse struct {
	ID          uuid.UUID      `json:"id"`
	EventID     uuid.UUID      `json:"event_id"`
	EventName   string         `json:"event_name"`
	RequestedAt string         `json:"requested_at"`
	Status      string         `json:"status"`
	Rooms       []RoomResponse `json:"rooms"`
}
