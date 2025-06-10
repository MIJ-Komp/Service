package exception

type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}

func (notFoundError NotFoundError) Error() string {
	return notFoundError.Message
}
