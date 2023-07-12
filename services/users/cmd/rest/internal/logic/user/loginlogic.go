package user

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	user, err := l.svcCtx.UserModel.CheckPWd(req.Count, req.Password, l.svcCtx.DB)
	if err != nil {
		return &types.LoginResp{
			BaseResponse: types.BaseResponse{
				Code: tools.ModelRecordNotFound.Code,
				Msg:  "用户不存在或密码错误",
			},
		}, nil
	}
	// roleArr := make([]string, 0)
	// for _, role := range user.Roles {
	// 	// fmt.Println(role.Menus, role.Resources)
	// 	// for _, menu := range role.Menus {
	// 	// 	menuArr = append(menuArr, types.Menu{}.FromMenuModel(*menu))
	// 	// }
	// 	roleArr = append(roleArr, role.Name)
	// }
	accessToken := models.GenToken(*user, l.svcCtx.Config.JwtAuth.AccessSecret)
	return &types.LoginResp{
		BaseResponse: types.BaseResponse{
			Code: 0,
			Msg:  "登录成功",
		},
		Data: types.UserInfoWithToken{
			AccessToken: accessToken,
			UserInfo: types.UserInfo{
				Avatar: user.Avatar,
				Email:  user.Email,
				Phone:  user.Phone,
				Age:    user.Age,
				// Role:   roleArr,
				Gender: user.Gender,
			},
			// Menus: menuArr,
		},
	}, nil
}
