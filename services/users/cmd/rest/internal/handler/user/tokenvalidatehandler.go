package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
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
		var src_method = r.Header.Get("X-Original-Method")
		// 获取权限API权限表达式
		permsRequired := tools.FormatPermissionFromUri(src_uri, src_method)
		fmt.Println("permsRequired:", permsRequired)
		// 判断api是否需要鉴权
		// fmt.Println(tools.ResourceAuthenticatedCacheKey(permsRequired), ">>>")
		val, err := svcCtx.RedisClient.Get(r.Context(), tools.ResourceAuthenticatedCacheKey).Result()
		if err == redis.Nil {
			// TODO 查不到接口权限，在查下数据库
			w.WriteHeader(http.StatusUnauthorized)
			httpx.ErrorCtx(r.Context(), w, err)
		} else if err != nil {
			fmt.Println("redis get authentication error:", err)
			panic(err)
		} else {
			var permission []tools.ResourceAuthenticatedItem
			json.Unmarshal([]byte(val), &permission)
			for _, p := range permission {
				if tools.MatchRegex(p.Resource, permsRequired) {
					if !p.Authenticated {
						fmt.Println("api不需要鉴权 -> ", src_uri, ":", src_method)
						w.WriteHeader(http.StatusOK)
						return
					}

				}
			}
			// if val == "0" {
			// 	// 不需要鉴权
			// 	//TODO 正则匹配
			// 	fmt.Println("api不需要鉴权 -> ", src_uri, ":", src_method)
			// 	w.WriteHeader(http.StatusOK)
			// 	return
			// }
		}
		fmt.Println("api需要鉴权 -> ", src_uri, ":", src_method)
		var req types.ValidateTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("Parse token error:", err)
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
			w.Header().Set("X-Super-Admin", strconv.FormatBool(resp.IsSuperAdmin))
			// w.Header().Set("X-Content-Type-Options", "")
			// w.Header().Set("X-Frame-Options", "")
			fmt.Println("response header", w.Header(), "Response:", resp, r.Context())
			w.WriteHeader(http.StatusOK)
			httpx.Ok(w)
		}
	}
}
