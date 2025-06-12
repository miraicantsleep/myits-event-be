package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/miraicantsleep/myits-event-be/utils"
)

type BookingRequestController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Approve(ctx *gin.Context)
	Reject(ctx *gin.Context)
}

type bookingRequestController struct {
	bookingRequestService service.BookingRequestService
}

func NewBookingRequestController(brService service.BookingRequestService) BookingRequestController {
	return &bookingRequestController{
		bookingRequestService: brService,
	}
}

func (c *bookingRequestController) Create(ctx *gin.Context) {
	var req dto.BookingRequestCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.bookingRequestService.CreateBookingRequest(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_BOOKING_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_BOOKING_REQUEST, result)
	ctx.JSON(http.StatusCreated, res)
}

func (c *bookingRequestController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.bookingRequestService.GetBookingRequestByID(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_BOOKING_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_BOOKING_REQUEST, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) GetAll(ctx *gin.Context) {
	results, err := c.bookingRequestService.GetAllBookingRequests(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_BOOKING_REQUESTS, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_BOOKING_REQUESTS, results)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req dto.BookingRequestUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// ADD THIS LINE
	role := ctx.MustGet("role").(string)

	// MODIFY THIS LINE to pass the role
	result, err := c.bookingRequestService.UpdateBookingRequest(ctx.Request.Context(), id, req, role)

	// MODIFY this error handling block
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_BOOKING_REQUEST, err.Error(), nil)
		// Use StatusForbidden for permission errors
		if err.Error() == "ormawa users are not permitted to change the booking status" {
			ctx.JSON(http.StatusForbidden, res)
			return
		}
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_BOOKING_REQUEST, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err := c.bookingRequestService.DeleteBookingRequest(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_BOOKING_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_BOOKING_REQUEST, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required for approval", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err := c.bookingRequestService.ApproveBookingRequest(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_APPROVE_BOOKING_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_APPROVE_BOOKING_REQUEST, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) Reject(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required for rejection", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err := c.bookingRequestService.RejectBookingRequest(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REJECT_BOOKING_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REJECT_BOOKING_REQUEST, nil)
	ctx.JSON(http.StatusOK, res)
}
