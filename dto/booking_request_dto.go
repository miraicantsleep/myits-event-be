package dto

import "github.com/google/uuid"

type BookingRequestCreateRequest struct {
	EventID uuid.UUID   `json:"event_id" binding:"required"`
	RoomIDs []uuid.UUID `json:"room_ids" binding:"required,min=1"`
}

type BookingRequestResponse struct {
	ID          uuid.UUID      `json:"id"`
	EventID     uuid.UUID      `json:"event_id"`
	EventName   string         `json:"event_name"` // Assuming we'll fetch event name
	RequestedAt string         `json:"requested_at"`
	Status      string         `json:"status"`
	Rooms       []RoomResponse `json:"rooms"` // Reusing existing RoomResponse if suitable
}

// MinimalRoomResponse for embedding in BookingRequestResponse to avoid circular dependencies if RoomResponse imports Booking-related DTOs.
// Or, ensure RoomResponse is simple enough. For now, let's assume RoomResponse is suitable.
// If not, we might need a specific RoomInfoDTO here.
