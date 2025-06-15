package types

import "api.mijkomp.com/models/enum"

type AppUser struct {
	UserId   uint
	UserType enum.UserType
}
