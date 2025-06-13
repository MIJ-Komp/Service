package response

type User struct {
	Id       uint   `json:"id"`
	UserName string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

type AuditTrail struct {
	Id       uint   `json:"id"`
	UserName string `json:"userName"`
}
