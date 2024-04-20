package apperror

import (
	"fmt"
)

type AppError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppErrors []AppError

func NewError(field, message string) AppError {
	return AppError{
		Field:   field,
		Message: message,
	}
}

func NewErrors(errors ...AppError) AppErrors {
	return errors
}

func EmailAlreadyExists() AppErrors {
	return NewErrors(NewError("email", "The email has already been taken"))
}

func InvalidCredentials() AppErrors {
	return NewErrors(NewError("credentials", "The provided credentials are invalid"))
}

func OldPasswordIncorrect() AppErrors {
	return NewErrors(NewError("old_password", "The old password is incorrect"))
}

func PostNotFound() AppErrors {
	return NewErrors(NewError("post", "The post could not be found"))
}

func UserNotFound() AppErrors {
	return NewErrors(NewError("user", "The user could not be found"))
}

func Unauthorized() AppErrors {
	return NewErrors(NewError("authorization", "You are not authorized to perform this action"))
}

func OwnershipRequired() AppErrors {
	return NewErrors(NewError("ownership", "You don't have permission to modify this resource"))
}

func InvalidID() AppErrors {
	return NewErrors(NewError("id", "The provided ID is invalid"))
}

func DatabaseError(err error) AppErrors {
	return NewErrors(NewError("database", fmt.Sprintf("A database error occurred: %v", err)))
}
