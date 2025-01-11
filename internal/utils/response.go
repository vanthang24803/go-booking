package utils

import "time"

type Response struct {
	Code      int         `json:"httpCode"`
	Result    interface{} `json:"result"`
	Timestamp string      `json:"timestamp"`
}

type Pagination struct {
	Code      int         `json:"httpCode"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Result    interface{} `json:"result"`
	Timestamp string      `json:"timestamp"`
}

func NewResponse(code int, result interface{}) *Response {
	return &Response{
		Code:      code,
		Result:    result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func NewPaginationResponse(total int, page int, limit int, result interface{}) *Pagination {
	return &Pagination{
		Code:      200,
		Total:     total,
		Page:      page,
		Limit:     limit,
		Result:    result,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
