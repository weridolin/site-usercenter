package types

import "github.com/weridolin/site-gateway/services/users/models"

func (u UserInfo) FromUserModel(user models.User) *UserInfo {
	return &UserInfo{
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		IsSuperAdmin: user.IsSuperAdmin,
		Age:          user.Age,
	}
}

func (m Menu) FromMenuModel(menu models.Menu) *Menu {
	return &Menu{
		Name:      menu.Name,
		Url:       menu.Url,
		Icon:      menu.Icon,
		ParentId:  menu.ParentId,
		Redirect:  menu.Redirect,
		Component: menu.Component,
		Type:      menu.Type,
	}
}
