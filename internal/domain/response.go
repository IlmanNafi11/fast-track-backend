package domain

import "time"

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	BaseResponse
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type ErrorResponse struct {
	BaseResponse
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
	PerPage      int `json:"per_page"`
}

type PaginatedResponse struct {
	SuccessResponse
	Meta PaginationMeta `json:"meta"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
