package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	model_pagination "github.com/roto17/zeus/lib/models/paginations"
)

func GetPaginationParams(c *gin.Context) (limitVal int, offsetVal int) {
	// Get page and limit from query, with default values
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Ensure page and limit have valid values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	return limit, offset
}

func GetPaginationMetadata(c *gin.Context, totalItems int64, limit int) model_pagination.PaginationMetadata {

	// Generate pagination metadata
	baseURL := fmt.Sprintf("http://%s%s", c.Request.Host, c.Request.URL.Path)

	// Get page from query or default to 1
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	// Limit validation
	if limit < 1 {
		limit = 10
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	// Generate page numbers
	pageNumbers := make([]model_pagination.PageInfo, totalPages)
	for i := 0; i < totalPages; i++ {
		pageNumbers[i] = model_pagination.PageInfo{
			Index:    i + 1,
			Page_url: fmt.Sprintf("%s?page=%d&limit=%d", baseURL, i+1, limit),
		}
	}

	// Determine next and previous pages
	var nextPageIndex *int
	if page < totalPages {
		next := page + 1
		nextPageIndex = &next
	}

	var previousPageIndex *int
	if page > 1 && page <= totalPages {
		prev := page - 1
		previousPageIndex = &prev
	}

	// Build URLs
	nextPageURL := ""
	if nextPageIndex != nil && *nextPageIndex > 0 && *nextPageIndex <= totalPages {
		nextPageURL = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, *nextPageIndex, limit)
	}

	previousPageURL := ""
	if previousPageIndex != nil && *previousPageIndex > 0 && *previousPageIndex <= totalPages {
		previousPageURL = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, *previousPageIndex, limit)
	}

	return model_pagination.PaginationMetadata{
		A_CurrentPage: page,
		B_Limit:       limit,
		C_TotalItems:  totalItems,
		D_TotalPages:  totalPages,
		// G_PreviousPageIndex: previousPageIndex,
		H_PreviousPageURL: previousPageURL,
		// E_NextPageIndex:     nextPageIndex,
		F_NextPageURL: nextPageURL,
		I_Pages:       pageNumbers,
	}
}
