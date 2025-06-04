package dto

import "errors"

const (
	EVENT_TYPE_ONLINE  = "online"
	EVENT_TYPE_OFFLINE = "offline"

	// FAILED
	MESSAGE_FAILED_CREATE_EVENT   = "failed create event"
	MESSAGE_FAILED_GET_EVENT      = "failed get event"
	MESSAGE_FAILED_DELETE_EVENT   = "failed delete event"
	MESSAGE_FAILED_UPDATE_EVENT   = "failed update event"
	MESSAGE_FAILED_GET_LIST_EVENT = "failed get list event"

	// SUCCESS
	MESSAGE_SUCCESS_CREATE_EVENT   = "success create event"
	MESSAGE_SUCCESS_GET_EVENT      = "success get event"
	MESSAGE_SUCCESS_DELETE_EVENT   = "success delete event"
	MESSAGE_SUCCESS_UPDATE_EVENT   = "success update event"
	MESSAGE_SUCCESS_GET_LIST_EVENT = "success get list event"
)

var (
	ErrCreateEvent   = errors.New("failed to create event")
	ErrGetEventById  = errors.New("failed to get event by id")
	ErrUpdateEvent   = errors.New("failed to update event")
	ErrDeleteEvent   = errors.New("failed to delete event")
	ErrGetAllEvent   = errors.New("failed to get all events")
	ErrEventNotFound = errors.New("event not found")
)

type (
	EventCreateRequest struct {
		Name        string `json:"name" form:"name" binding:"required,min=2,max=100"`
		Description string `json:"description" form:"description" binding:"required,min=10,max=500"`
		Start_Time  string `json:"start_time" form:"start_time" binding:"required"`
		End_Time    string `json:"end_time" form:"end_time" binding:"required"`
		Event_Type  string `json:"event_type" form:"event_type" binding:"required,oneof=online offline"`
	}

	EventUpdateRequest struct {
		Name        string `json:"name" form:"name" binding:"omitempty,min=2,max=100"`
		Description string `json:"description" form:"description" binding:"omitempty,min=10,max=500"`
		Start_Time  string `json:"start_time" form:"start_time" binding:"omitempty"`
		End_Time    string `json:"end_time" form:"end_time" binding:"omitempty"`
		Event_Type  string `json:"event_type" form:"event_type" binding:"omitempty,oneof=online offline"`
	}

	GetAllEventRepositoryResponse struct {
		Events             []EventResponse `json:"Events"`
		PaginationResponse PaginationResponse
	}

	EventResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Start_Time  string `json:"start_time"`
		End_Time    string `json:"end_time"`
		Created_By  string `json:"created_by"`
		Event_Type  string `json:"event_type"`
	}

	EventPaginationResponse struct {
		Data []EventResponse `json:"data"`
		PaginationResponse
	}
)
