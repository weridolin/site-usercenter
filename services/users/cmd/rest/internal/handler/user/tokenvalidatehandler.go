package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/logic/user"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/tools"
	"github.com/zeromicro/go-zero/rest/httpx"
	xerrors "github.com/zeromicro/x/errors"
)

func TokenValidateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var src_uri = r.Header.Get("X-Original-Request-Uri")
		// 获取权限API权限表达式
		permsRequired := tools.FormatPermissionFromUri(src_uri, strings.ToLower(r.Method))
		fmt.Println("permsRequired:", permsRequired, "src_uri:", src_uri)
		if src_uri == "/usercenter/api/v1/login" || src_uri == "/usercenter/api/v1/register" {
			w.WriteHeader(http.StatusOK)
			return
		}
		var req types.ValidateTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := user.NewTokenValidateLogic(r.Context(), svcCtx)
		resp, err := l.TokenValidate(&req, permsRequired)
		if err != nil {
			w.WriteHeader(err.(*xerrors.CodeMsg).Code)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Header().Set("X-Forwarded-User", resp.UserId)
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
