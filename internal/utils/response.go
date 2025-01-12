package utils

import "time"

type Response struct {
	Code      int         `json:"httpCode"`
	Result    interface{} `json:"result"`
	Timestamp string      `json:"timestamp"`
}

type Pagination struct {
	Code        int         `json:"httpCode"`
	Limit       int         `json:"limit"`
	CurrentPage int         `json:"current_page"`
	Total       int         `json:"total_items"`
	TotalPages  int         `json:"total_pages"`
	Result      interface{} `json:"result"`
	Timestamp   string      `json:"timestamp"`
}

func NewResponse(code int, result interface{}) *Response {
	return &Response{
		Code:      code,
		Result:    result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func NewPaginationResponse(total int, totalPages int, page int, limit int, result interface{}) *Pagination {
	return &Pagination{
		Code:        200,
		Total:       total,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		Result:      result,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}
