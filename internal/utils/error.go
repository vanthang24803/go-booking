package utils

import "time"

type AppError struct {
	Code      int    `json:"httpCode"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
