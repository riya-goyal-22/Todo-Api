package error

import (
	"net/http"
	"todo/model"
)

func NewNotFoundError(message string) *model.Error {
	return &model.Error{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewInternalServerError(message string) *model.Error {
	return &model.Error{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *model.Error {
	return &model.Error{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *model.Error {
	return &model.Error{
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}
