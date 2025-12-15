package utils

import (
	"time"

	apierrors "ContractAnalysis/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, err *apierrors.APIError) {
	c.JSON(int(err.Code), Response{
		Code:    int(err.Code),
		Message: err.Message,
		Error: map[string]interface{}{
			"type":    err.Type,
			"details": err.Details,
		},
		Timestamp: time.Now().Unix(),
	})
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Items      interface{}        `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(items interface{}, page, limit, total int) *PaginatedResponse {
	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Items: items,
		Pagination: PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// PaginatedSuccessResponse sends a paginated success response
func PaginatedSuccessResponse(c *gin.Context, code int, message string, items interface{}, page, limit, total int) {
	data := NewPaginatedResponse(items, page, limit, total)
	SuccessResponse(c, code, message, data)
}
