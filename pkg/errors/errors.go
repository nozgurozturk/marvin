package errors

import "net/http"

type AppError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

func BadRequest(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func UnprocessableEntity(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusUnprocessableEntity,
		Error:   "unprocessable_entity",
	}
}

func InternalServer(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "server_error",
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func AlreadyExist(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusConflict,
		Error:   "already_exist",
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   "forbidden",
	}
}