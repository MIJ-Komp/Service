package request

type LoginUserPayload struct {
	Email    string `validate:"required,max=500,min=5"`
	Password string `validate:"required,max=500,min=3"`
}

type LoginGooglePayload struct {
	OpenId string
	Name   string
	Email  string
	Image  string
}
