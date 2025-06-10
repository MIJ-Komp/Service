package exception

import (
	"fmt"
)

type ValidationError struct {
	Message string
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

func NewModelValidationError(err interface{}) {
	if err != nil {
		panic(NewValidationError(fmt.Sprintf("Please check your data. %s", err.(string))))
	}
}

func (validationError ValidationError) Error() string {
	return validationError.Message
}
