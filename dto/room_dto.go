package dto

import "errors"

const (
	// sucess
	MESSAGE_SUCCESS_CREATE_ROOM    = "Success create room"
	MESSAGE_SUCCESS_GET_ROOM_BY_ID = "Success get room by id"
	MESSAGE_SUCCESS_GET_ALL_ROOM   = "Success get all rooms"
	MESSAGE_SUCCESS_UPDATE_ROOM    = "Success update room"
	MESSAGE_SUCCESS_DELETE_ROOM    = "Success delete room"

	// Failed
	MESSAGE_FAILED_CREATE_ROOM    = "Failed create room"
	MESSAGE_FAILED_GET_ROOM_BY_ID = "Failed get room by id"
	MESSAGE_FAILED_GET_ALL_ROOM   = "Failed get all rooms"
	MESSAGE_FAILED_UPDATE_ROOM    = "Failed update room"
	MESSAGE_FAILED_DELETE_ROOM    = "Failed delete room"
)

// errs
var (
	ErrCreateRoom          = errors.New("failed to create room")
	ErrGetRoomByID         = errors.New("failed to get room by id")
	ErrGetRoomByName       = errors.New("failed to get room by name")
	ErrGetAllRoom          = errors.New("failed to get all rooms")
	ErrUpdateRoom          = errors.New("failed to update room")
	ErrRoomNotFound        = errors.New("room not found")
	ErrDeleteRoom          = errors.New("failed to delete room")
	ErrRoomAlreadyExists   = errors.New("room already exists")
	ErrRoomInvalidCapacity = errors.New("room capacity must be greater than zero")
)

type (
	RoomCreateRequest struct {
		Name         string `json:"name" binding:"required"`
		Capacity     int    `json:"capacity" binding:"required,gt=0"`
		DepartmentID string `json:"department_id" binding:"omitempty"` // required for admin, optional for department role
	}

	RoomUpdateRequest struct {
		Name         string `json:"name" binding:"omitempty"`
		Capacity     int    `json:"capacity" binding:"omitempty"`
		DepartmentID string `json:"department_id" binding:"omitempty"`
	}

	RoomResponse struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Department string `json:"department"`
		Capacity   int    `json:"capacity"`
	}
)
