package data_mapper

import (
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/response"
)

func MapAuditTrail(user entity.User) response.AuditTrail {
	return response.AuditTrail{
		Id:       user.Id,
		UserName: user.UserName,
	}
}
