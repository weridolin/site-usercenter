package user

import (
	"fmt"
	"net/http"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/logic/user"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TokenRefreshHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		fmt.Println("req:", r.Header)
		l := user.NewTokenRefreshLogic(r.Context(), svcCtx)
		resp, err := l.TokenRefresh(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
