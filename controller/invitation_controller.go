package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/miraicantsleep/myits-event-be/utils"
)

type (
	InvitationController interface {
		Create(ctx *gin.Context)
		GetInvitationByID(ctx *gin.Context)
		GetInvitationByEventID(ctx *gin.Context)
		GetAllInvitations(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	invitationController struct {
		invitationService service.InvitationService
		jwtService        service.JWTService
	}
)

func NewInvitationController(
	invitationService service.InvitationService,
	jwtService service.JWTService,
) InvitationController {
	return &invitationController{
		invitationService: invitationService,
		jwtService:        jwtService,
	}
}

func (c *invitationController) Create(ctx *gin.Context) {
	var req dto.CreateInvitationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.invitationService.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVITATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_INVITATION, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *invitationController) GetInvitationByID(ctx *gin.Context) {
	invitationId := ctx.Param("id")
	if invitationId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVITATION_BY_ID, "invitation ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.invitationService.GetInvitationByID(ctx.Request.Context(), invitationId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVITATION_BY_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_INVITATION_BY_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *invitationController) GetInvitationByEventID(ctx *gin.Context) {
	eventId := ctx.Param("event_id")
	if eventId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVITATION_BY_EVENT_ID, "event ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.invitationService.GetInvitationByEventID(ctx.Request.Context(), eventId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVITATION_BY_EVENT_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_INVITATION_BY_EVENT_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *invitationController) GetAllInvitations(ctx *gin.Context) {
	result, err := c.invitationService.GetAllInvitations(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_INVITATIONS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_INVITATIONS, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *invitationController) Update(ctx *gin.Context) {
	invitationId := ctx.Param("id")
	if invitationId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVITATION, "invitation ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UpdateInvitationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.invitationService.Update(ctx.Request.Context(), invitationId, req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVITATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_INVITATION, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *invitationController) Delete(ctx *gin.Context) {
	invitationId := ctx.Param("id")
	if invitationId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_INVITATION, "invitation ID is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.invitationService.Delete(ctx.Request.Context(), invitationId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_INVITATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_INVITATION, nil)
	ctx.JSON(http.StatusOK, res)
}
