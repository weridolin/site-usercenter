package types

import "github.com/weridolin/site-gateway/services/users/models"

func (u UserInfo) FromUserModel(user models.User) *UserInfo {
	return &UserInfo{
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		Role:         user.Role,
		IsSuperAdmin: user.IsSuperAdmin,
		Age:          user.Age,
	}
}
