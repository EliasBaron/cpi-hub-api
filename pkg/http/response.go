package http

import (
	"cpi-hub-api/pkg/apperror"
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationData struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data,omitempty"`
	Pagination PaginationData `json:"pagination"`
	Error      string         `json:"error,omitempty"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSONResponse(w, http.StatusOK, response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := Response{
		Success: false,
		Error:   message,
	}
	JSONResponse(w, statusCode, response)
}

func CreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSONResponse(w, http.StatusCreated, response)
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
func PaginatedSuccessResponse(w http.ResponseWriter, message string, data interface{}, page, pageSize, total int) {
	response := PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: PaginationData{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}
	JSONResponse(w, http.StatusOK, response)
}
