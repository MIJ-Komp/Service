package service

import (
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"github.com/google/uuid"
)

type UserService interface {
	Create(user request.RegisterUserPayload) response.User
	GetByEmail(email string) response.User
	LoginUser(user request.LoginUserPayload) string
	GetById(userId uint) response.User
	HasToken(userId uint, userToken uuid.UUID) bool
}
