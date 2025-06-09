package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/miraicantsleep/myits-event-be/utils"
)

type (
	RoomController interface {
		Create(ctx *gin.Context)
		GetRoomByID(ctx *gin.Context)
		GetAllRoom(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	roomController struct {
		roomService service.RoomService
	}
)

func NewRoomController(rs service.RoomService) RoomController {
	return &roomController{
		roomService: rs,
	}
}

func (c *roomController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	role := ctx.MustGet("role").(string)
	var roomRequest dto.RoomCreateRequest
	if err := ctx.ShouldBind(&roomRequest); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.roomService.Create(ctx.Request.Context(), roomRequest, userId, role)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ROOM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ROOM, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *roomController) GetRoomByID(ctx *gin.Context) {
	roomID := ctx.Param("id")
	if roomID == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, "Room ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.roomService.GetRoomByID(ctx.Request.Context(), roomID)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ROOM_BY_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ROOM_BY_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *roomController) GetAllRoom(ctx *gin.Context) {
	result, err := c.roomService.GetAllRoom(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_ROOM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_ROOM, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *roomController) Update(ctx *gin.Context) {
	var roomRequest dto.RoomUpdateRequest
	if err := ctx.ShouldBind(&roomRequest); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	roomID := ctx.Param("id")
	if roomID == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, "Room ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.roomService.Update(ctx.Request.Context(), roomID, roomRequest)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ROOM, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ROOM, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *roomController) Delete(ctx *gin.Context) {
	roomID := ctx.Param("id")
	if roomID == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, "Room ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.roomService.Delete(ctx.Request.Context(), roomID); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ROOM, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ROOM, nil)
	ctx.JSON(http.StatusOK, res)
}
