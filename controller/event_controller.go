package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/miraicantsleep/myits-event-be/utils"
)

type (
	EventController interface {
		Create(ctx *gin.Context)
		GetAllEvent(ctx *gin.Context)
		GetEventByID(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GetEventAttendees(ctx *gin.Context)
	}

	eventController struct {
		eventService service.EventService
	}
)

func NewEventController(es service.EventService) EventController {
	return &eventController{
		eventService: es,
	}
}

func (c *eventController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	var req dto.EventCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.eventService.Create(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_EVENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_EVENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *eventController) GetAllEvent(ctx *gin.Context) {
	user_role := ctx.MustGet("role").(string)
	user_id := ctx.MustGet("user_id").(string)
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.eventService.GetAllEventWithPagination(ctx.Request.Context(), req, user_role, user_id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_EVENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_EVENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *eventController) GetEventByID(ctx *gin.Context) {
	eventId := ctx.Param("id")
	if eventId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, "event ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.eventService.GetEventById(ctx.Request.Context(), eventId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_EVENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_EVENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *eventController) Update(ctx *gin.Context) {
	eventId := ctx.Param("id")
	if eventId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, "event ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.EventUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.eventService.Update(ctx.Request.Context(), req, eventId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_EVENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_EVENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *eventController) Delete(ctx *gin.Context) {
	eventId := ctx.Param("id")
	if eventId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_TOKEN_NOT_VALID, "event ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.eventService.Delete(ctx.Request.Context(), eventId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_EVENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_EVENT, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *eventController) GetEventAttendees(ctx *gin.Context) {
	eventId := ctx.Param("id")
	if eventId == "" {
		res := utils.BuildResponseFailed("Event ID is required", "Event ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	attendees, err := c.eventService.GetEventAttendees(ctx.Request.Context(), eventId)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to get event attendees", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Successfully fetched event attendees", attendees)
	ctx.JSON(http.StatusOK, res)
}
