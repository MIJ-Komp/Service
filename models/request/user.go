package request

type RegisterUserPayload struct {
	FullName string `json:"fullName" validate:"required,min=1,max=128"`
	UserName string `json:"userName" validate:"required,min=5,max=16"`
	Email    string `json:"email" validate:"required,min=5,max=128"`
	Password string `json:"password" validate:"required,min=1,max=64"`
}
