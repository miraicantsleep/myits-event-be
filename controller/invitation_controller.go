package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/entity" // Added for RSVPStatus constants
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
		ScanQRCode(ctx *gin.Context)
		AcceptRSVP(ctx *gin.Context)  // New method
		DeclineRSVP(ctx *gin.Context) // New method
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

func (c *invitationController) ScanQRCode(ctx *gin.Context) {
	qrCode := ctx.Param("qr_code") // Get QR code from path parameter
	if qrCode == "" {
		res := utils.BuildResponseFailed("QR code is required", "Missing QR code in path", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.invitationService.ScanQRCode(ctx.Request.Context(), qrCode)
	if err != nil {
		// Consider different error types for specific responses
		// For example, if err.Error() is "QR code not found" or "QR code already used"
		// you might want to return a specific HTTP status code like http.StatusNotFound or http.StatusConflict
		// For now, a generic bad request or a more specific error message will do.
		// The service layer returns descriptive errors.
		res := utils.BuildResponseFailed("Failed to process QR code", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res) // Or potentially http.StatusNotFound / http.StatusConflict based on err
		return
	}

	res := utils.BuildResponseSuccess("QR code scanned successfully", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *invitationController) AcceptRSVP(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.HTML(http.StatusBadRequest, "rsvp_error.html", gin.H{
			"title":   "RSVP Error",
			"message": "Invalid RSVP link. No token provided.",
		})
		return
	}

	// Call the service method (to be created in the next plan step)
	// For now, assume it's called ProcessRSVP and it's part of invitationService interface
	err := c.invitationService.ProcessRSVP(ctx.Request.Context(), token, entity.RSVPStatusAccepted) // entity.RSVPStatusAccepted = "accepted"

	if err != nil {
		// Determine message based on error type if possible
		// For now, a generic error message.
		// Specific errors like "already RSVP'd" or "token not found" will be handled by the service.
		ctx.HTML(http.StatusOK, "rsvp_error.html", gin.H{ // Using StatusOK for error page display from email link
			"title":   "RSVP Problem",
			"message": err.Error(), // Display service error message directly
		})
		return
	}

	ctx.HTML(http.StatusOK, "rsvp_success.html", gin.H{
		"title":   "RSVP Confirmed",
		"message": "Thank you! Your RSVP has been successfully recorded as ACCEPTED.",
		"status":  "Accepted",
	})
}

func (c *invitationController) DeclineRSVP(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.HTML(http.StatusBadRequest, "rsvp_error.html", gin.H{
			"title":   "RSVP Error",
			"message": "Invalid RSVP link. No token provided.",
		})
		return
	}

	err := c.invitationService.ProcessRSVP(ctx.Request.Context(), token, entity.RSVPStatusDeclined) // entity.RSVPStatusDeclined = "declined"

	if err != nil {
		ctx.HTML(http.StatusOK, "rsvp_error.html", gin.H{
			"title":   "RSVP Problem",
			"message": err.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "rsvp_success.html", gin.H{
		"title":   "RSVP Updated",
		"message": "Your RSVP has been successfully recorded as DECLINED.",
		"status":  "Declined",
	})
}
