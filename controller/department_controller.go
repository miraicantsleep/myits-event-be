package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/miraicantsleep/myits-event-be/utils"
)

type (
	DepartmentController interface {
		Create(ctx *gin.Context)
		GetDepartmentByID(ctx *gin.Context)
		GetAllDepartment(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	departmentController struct {
		departmentService service.DepartmentService
		userService       service.UserService
	}
)

func NewDepartmentController(ds service.DepartmentService, us service.UserService) DepartmentController {
	return &departmentController{
		departmentService: ds,
		userService:       us,
	}
}

func (c *departmentController) Create(ctx *gin.Context) {
	var department dto.DepartmentCreateRequest
	if err := ctx.ShouldBind(&department); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.departmentService.Create(ctx.Request.Context(), department)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_DEPARTMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_DEPARTMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *departmentController) GetAllDepartment(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.departmentService.GetAllDepartmentWithPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_DEPARTMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	resp := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_LIST_DEPARTMENT,
		Data:    result.Data,
		Meta:    result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *departmentController) GetDepartmentByID(ctx *gin.Context) {
	departmentId := ctx.Param("id")

	result, err := c.departmentService.GetDepartmentById(ctx.Request.Context(), departmentId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DEPARTMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DEPARTMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *departmentController) Update(ctx *gin.Context) {
	departmentId := ctx.Param("id")
	if departmentId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, "failed to get department id", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.DepartmentUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.departmentService.Update(ctx.Request.Context(), req, departmentId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_DEPARTMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_DEPARTMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *departmentController) Delete(ctx *gin.Context) {
	departmentId := ctx.Param("id")
	if departmentId == "" {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, "failed to get department id", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.departmentService.Delete(ctx.Request.Context(), departmentId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_DEPARTMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_DEPARTMENT, nil)
	ctx.JSON(http.StatusOK, res)
}
