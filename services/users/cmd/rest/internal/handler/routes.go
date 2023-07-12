// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	user "github.com/weridolin/site-gateway/services/users/cmd/rest/internal/handler/user"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/logout",
				Handler: user.LogoutHandler(serverCtx),
			},
		},
		rest.WithPrefix("/usercenter/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/token/refresh",
				Handler: user.TokenRefreshHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/token/validate",
				Handler: user.TokenValidateHandler(serverCtx),
			},
		},
		// rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret), // todo改成接口里面实现验证逻辑
	)
}
