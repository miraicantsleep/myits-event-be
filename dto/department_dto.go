package dto

import (
	"errors"
)

const (
	// Failed
	MESSAGE_FAILED_CREATE_DEPARTMENT   = "failed create department"
	MESSAGE_FAILED_GET_LIST_DEPARTMENT = "failed get list department"
	MESSAGE_FAILED_GET_DEPARTMENT      = "failed get department"
	MESSAGE_FAILED_UPDATE_DEPARTMENT   = "failed update department"
	MESSAGE_FAILED_DELETE_DEPARTMENT   = "failed delete department"

	// Success
	MESSAGE_SUCCESS_CREATE_DEPARTMENT   = "success create department"
	MESSAGE_SUCCESS_GET_LIST_DEPARTMENT = "success get list department"
	MESSAGE_SUCCESS_GET_DEPARTMENT      = "success get department"
	MESSAGE_SUCCESS_UPDATE_DEPARTMENT   = "success update department"
	MESSAGE_SUCCESS_DELETE_DEPARTMENT   = "success delete department"
)

var (
	ErrCreateDepartment     = errors.New("failed to create department")
	ErrGetDepartmentById    = errors.New("failed to get department by id")
	ErrGetDepartmentByEmail = errors.New("failed to get department by email")
	ErrUpdateDepartment     = errors.New("failed to update department")
	ErrDepartmentNotFound   = errors.New("department not found")
	ErrDeleteDepartment     = errors.New("failed to delete department")
)

type (
	DepartmentCreateRequest struct {
		Name     string `json:"name" form:"name" binding:"required,min=2,max=100"`
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required,min=8"`
		Faculty  string `json:"faculty" form:"faculty" binding:"required,min=2,max=100"`
	}

	DepartmentResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Faculty string `json:"faculty"`
		Email   string `json:"email"`
	}

	DepartmentPaginationResponse struct {
		Data []DepartmentResponse `json:"data"`
		PaginationResponse
	}

	GetAllDepartmentRepositoryResponse struct {
		Departments []DepartmentResponse `json:"Departments"`
		PaginationResponse
	}

	DepartmentUpdateRequest struct {
		Name    string `json:"name" form:"name" binding:"omitempty,min=2,max=100"`
		Faculty string `json:"faculty" form:"faculty" binding:"omitempty,min=2,max=100"`
		Email   string `json:"email" form:"email" binding:"required,email"`
	}

	DepartmentUpdateResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Faculty string `json:"faculty"`
		Email   string `json:"email"`
	}
)
