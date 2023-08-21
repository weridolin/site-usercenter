package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/logic/user"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(r.Header.Get("X-User"))
		var req types.GetMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("register parse params user error", err, r.Header)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRbacLogic(r.Context(), svcCtx)
		resp, err := l.GetMenu(userId, &req)
		if err != nil {
			fmt.Println("register user error", err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
