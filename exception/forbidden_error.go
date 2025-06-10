package exception

type ForbiddenError struct {
	Message string
}

func NewForbiddenError(message string) ForbiddenError {
	return ForbiddenError{Message: message}
}
func (forbiddenError ForbiddenError) Error() string {
	return forbiddenError.Message
}
