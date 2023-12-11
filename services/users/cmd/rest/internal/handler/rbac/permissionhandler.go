package rbac

import (
	"net/http"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/logic/rbac"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

//重新加载资源权限s
func ReLoadResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReLoadResourceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rbac.NewPermissionLogic(r.Context(), svcCtx)
		resp, err := l.ReLoadResourceLogic(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
