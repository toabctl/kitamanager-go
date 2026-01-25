package models

import "fmt"

// PaginationParams holds pagination request parameters
type PaginationParams struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// Validate checks for explicitly invalid pagination values.
// Returns an error if page or limit are negative, or if limit exceeds 100.
// Zero values are allowed as they indicate "not provided" and will be set by SetDefaults.
func (p *PaginationParams) Validate() error {
	if p.Page < 0 {
		return fmt.Errorf("page must be positive")
	}
	if p.Limit < 0 {
		return fmt.Errorf("limit must be positive")
	}
	if p.Limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}
	return nil
}

// SetDefaults sets default values for pagination parameters
func (p *PaginationParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
}

// Offset calculates the offset for database queries
func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginatedResponse wraps list responses with pagination metadata
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"20"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"5"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, page, limit int, total int64) PaginatedResponse[T] {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	return PaginatedResponse[T]{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
