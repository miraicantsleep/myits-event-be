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

	// userID := ctx.MustGet("user_id").(string)

	result, err := c.bookingRequestService.CreateBookingRequest(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to create booking request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Booking request created successfully", result)
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
		// Handle not found error specifically
		// if errors.Is(err, gorm.ErrRecordNotFound) { // Or a custom error from service
		// 	res := utils.BuildResponseFailed("Booking request not found", err.Error(), nil)
		// 	ctx.JSON(http.StatusNotFound, res)
		// 	return
		// }
		res := utils.BuildResponseFailed("Failed to get booking request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Booking request retrieved successfully", result)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) GetAll(ctx *gin.Context) {
	results, err := c.bookingRequestService.GetAllBookingRequests(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed("Failed to get all booking requests", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("All booking requests retrieved successfully", results)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required for approval", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// TODO: Add role-based authorization if needed
	// role := ctx.MustGet("role").(string)
	// if role != "admin" && role != "department_admin" { // Example roles
	// 	res := utils.BuildResponseFailed("Unauthorized", "Insufficient permissions", nil)
	// 	ctx.JSON(http.StatusForbidden, res)
	// 	return
	// }

	err := c.bookingRequestService.ApproveBookingRequest(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to approve booking request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Booking request approved successfully", nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookingRequestController) Reject(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		res := utils.BuildResponseFailed("Booking request ID is required for rejection", "ID is empty", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// TODO: Add role-based authorization if needed

	err := c.bookingRequestService.RejectBookingRequest(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to reject booking request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess("Booking request rejected successfully", nil)
	ctx.JSON(http.StatusOK, res)
}
