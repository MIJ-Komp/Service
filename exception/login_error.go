package exception

type LoginError struct {
	Message string
}

func NewLoginError(message string) LoginError {
	return LoginError{Message: message}
}

func (loginError LoginError) Error() string {
	return loginError.Message
}
