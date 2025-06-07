package errors

import (
	"net/http"

	"github.com/robatipoor/task-scheduler/internal/master/dto"
)

type ResponseError interface {
	Response() (int, dto.ErrorResponse)
}

type UnauthorizedError struct {
	Message string
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		Message: message,
	}
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

func (e UnauthorizedError) Response() (int, dto.ErrorResponse) {
	return http.StatusUnauthorized, RespondWithError(e.Message)
}

type UniqueConstraintViolationError struct {
	Message string
}

func NewUniqueConstraintViolationError(message string) UniqueConstraintViolationError {
	return UniqueConstraintViolationError{
		Message: message,
	}
}

func (e UniqueConstraintViolationError) Error() string {
	return e.Message
}

func (e UniqueConstraintViolationError) Response() (int, dto.ErrorResponse) {
	return http.StatusBadRequest, RespondWithError(e.Message)
}

func RespondWithError(message string) dto.ErrorResponse {
	return dto.ErrorResponse{Error: message}
}
