package http

import (
	"cpi-hub-api/pkg/apperror"
	"encoding/json"
	"net/http"
)

type ErrorResponseData struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type PaginationData struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PaginatedResponse struct {
	Data       interface{}    `json:"data,omitempty"`
	Pagination PaginationData `json:"pagination"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, data)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponseData{
		StatusCode: statusCode,
		Message:    message,
	}
	JSONResponse(w, statusCode, response)
}

func CreatedResponse(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusCreated, data)
}

func BadRequestResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusBadRequest, message)
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message)
}

func InternalServerErrorResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusInternalServerError, message)
}

// NewError creates an error response using the apperror package
func NewError(w http.ResponseWriter, err error) {
	statusCode, message := apperror.StatusCodeAndMessage(err)
	ErrorResponse(w, statusCode, message)
}

// PaginatedSuccessResponse creates a paginated success response
func PaginatedSuccessResponse(w http.ResponseWriter, data interface{}, page, pageSize, total int) {
	response := PaginatedResponse{
		Data: data,
		Pagination: PaginationData{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}
	JSONResponse(w, http.StatusOK, response)
}
