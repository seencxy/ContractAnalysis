package utils

import (
	"strconv"

	apierrors "ContractAnalysis/pkg/errors"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultPage is the default page number
	DefaultPage = 1

	// DefaultLimit is the default page size
	DefaultLimit = 20

	// MaxLimit is the maximum page size
	MaxLimit = 100
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

// ParsePaginationParams parses pagination parameters from query string
func ParsePaginationParams(c *gin.Context) (*PaginationParams, *apierrors.APIError) {
	page := DefaultPage
	limit := DefaultLimit

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return nil, apierrors.NewBadRequestError("Invalid page parameter", "page must be a positive integer")
		}
		page = p
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			return nil, apierrors.NewBadRequestError("Invalid limit parameter", "limit must be a positive integer")
		}
		if l > MaxLimit {
			l = MaxLimit
		}
		limit = l
	}

	offset := (page - 1) * limit

	return &PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// CalculateOffset calculates the offset for pagination
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calculates the total number of pages
func CalculateTotalPages(total, limit int) int {
	if limit == 0 {
		return 0
	}
	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}
	return totalPages
}
