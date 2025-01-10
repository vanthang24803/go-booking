package utils

import "time"

type Response struct {
	Code      int         `json:"httpCode"`
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
