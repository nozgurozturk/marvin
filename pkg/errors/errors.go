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
		Error:   http.StatusText(http.StatusBadRequest),
	}
}

func UnprocessableEntity(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusUnprocessableEntity,
		Error:   http.StatusText(http.StatusUnprocessableEntity),
	}
}

func InternalServer(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   http.StatusText(http.StatusInternalServerError),
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   http.StatusText(http.StatusNotFound),
	}
}

func AlreadyExist(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusConflict,
		Error:   http.StatusText(http.StatusConflict),
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   http.StatusText(http.StatusUnauthorized),
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   http.StatusText(http.StatusForbidden),
	}
}