package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/types"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"
	"github.com/zeromicro/go-zero/rest/httpx"
)

/*
获取登录二维码
结果将图片base64编码后返回
*/
func GetQrLoginCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetLoginQrCode:", r.Header)
		platform := r.URL.Query().Get("platform")
		uuid := tools.GetUUID()
		ctx := r.Context()
		_, err := svcCtx.RedisClient.Set(ctx, uuid, "not update", time.Minute*5).Result() // 生成二维码, 5分钟有效，代表此次扫码请求
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}
		// 生成二维码图片
		var png []byte
		//固定方法
		var content struct {
			Uuid     string
			CheckUrl string
		}
		content.Uuid = uuid
		switch platform {
		case "dingtalk":
			redirect_url := url.QueryEscape(os.Getenv("DINGTALKREDIRECTURI"))
			content.CheckUrl = fmt.Sprintf(`https://oapi.dingtalk.com/connect/qrconnect?appid=%s&response_type=code&scope=snsapi_login&state=STATE&redirect_uri=%s`, os.Getenv("DINGTALKAPPID"), redirect_url)
		case "github":
			redirect_url := url.QueryEscape(os.Getenv("GITHUBREDIRECTURI"))
			content.CheckUrl = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s", os.Getenv("GITHUBCLIENTID"), redirect_url, uuid)
		default:
			content.CheckUrl = "暂不支持该平台登录"
		}
		fmt.Println("content:", content)
		// content_str, _ := json.Marshal(content)
		png, err = qrcode.Encode(string(content.CheckUrl), qrcode.Medium, 256)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
			return
		}
		//文件流需要使用base64编码后才可使用
		res := base64.StdEncoding.EncodeToString(png)
		http.SetCookie(w, &http.Cookie{ //非扫码情况下第三方登录设置cookie,但可能存在伪造假站点攻击
			Name:     "state",
			Value:    uuid,
			Path:     "/usercenter/api/v1",
			HttpOnly: true,
			MaxAge:   60 * 5})
		httpx.OkJsonCtx(ctx, w, types.GetQrLoginCodeResp{Data: res})
	}
}

/*
获取登录二维码状态
1. 获取uuid, 通过uuid获取redis中的值，如果存在则返回登录成功
2. 如果key不存在，则代表二维码过期
3. 如果值为“”not update“”则代表二维码未被扫描
*/
func GetQrLoginStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetLoginQrCode:", r.Header)
		ctx := r.Context()
		uuid := r.URL.Query().Get("uuid")
		if uuid == "" {
			httpx.ErrorCtx(ctx, w, fmt.Errorf("uuid is empty"))
		}
		if value, err := svcCtx.RedisClient.Get(ctx, uuid).Result(); err != nil {
			fmt.Println("GetLoginQrCode:", err)
			httpx.OkJsonCtx(ctx, w, types.GetLoginQrCodeStatusResp{
				BaseResponse: types.BaseResponse{
					Code: tools.RedisKeyNotFoundError.Code,
					Msg:  tools.RedisKeyNotFoundError.Msg,
				},
			})
		} else {
			httpx.OkJsonCtx(ctx, w, types.GetLoginQrCodeStatusResp{
				BaseResponse: types.BaseResponse{
					Code: 0,
					Msg:  "",
				},
				Data: types.LoginQrCodeStatusResp{Status: value},
			})
		}
	}
}

/*
二维码登录
*/
// func LoginByQrCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// https://oapi.dingtalk.com/connect/qrconnect?appid=SuiteKey &response_type=code&scope=snsapi_login&state=STATE&redirect_uri=REDIRECT_URI
// fmt.Println("LoginByQrCodeHandler:", r.Header)
// ctx := r.Context()
// var req types.LoginByQrCodeReq
// if err := httpx.Parse(r, &req); err != nil {
// 	httpx.ErrorCtx(ctx, w, err)
// 	return
// }
// l := user.NewLoginByQrCodeLogic(ctx, svcCtx)
// resp, err := l.LoginByQrCode(&req)
// if err != nil {
// 	httpx.ErrorCtx(ctx, w, err)
// } else {
// 	httpx.OkJsonCtx(ctx, w, resp)
// }
// 	}
// }

/*********************        第三方登录        **************************/
func GetThirdLoginUriHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		platform := r.URL.Query().Get("platform")
		// 生成一个token,用来进行防伪造和过期限制
		uuid := tools.GetUUID()
		token := models.GenStateToken(svcCtx.Config.JwtAuth.AccessSecret, time.Second*60*5)
		_, err := svcCtx.RedisClient.Set(r.Context(), tools.OauthStateCacheKey(uuid), token, 0).Result()
		if err != nil {
			fmt.Println("get third platform login uri cache state key error ->", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		var uri string
		switch platform {
		case "github":
			// redirect_url := url.QueryEscape(os.Getenv("GITHUBREDIRECTURI"))
			// redirect_url := os.Getenv("GITHUBREDIRECTURI")
			// client_id := os.Getenv("GITHUBCLIENTID")
			redirect_url := "http://127.0.0.1:8080/usercenter/api/v1/third-platform/github-login"
			client_id := "Ov23liUpzHCPVymdPW2R"

			uri = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s", client_id, redirect_url, uuid)
			http.SetCookie(w, &http.Cookie{
				Name:     "state", //用于防伪造校验
				Value:    token,
				Path:     "/usercenter/api/v1/third-platform/github-login",
				HttpOnly: true,
				MaxAge:   60 * 5})
		case "gitee":
			redirect_url := os.Getenv("GITEEREDIRECTURI")
			client_id:= os.Getenv("GITEECLIENTID")
			// redirect_url := "http://127.0.0.1:8080/usercenter/api/v1/third-platform/gitee-login"
			// client_id := "c45c44c778f0ad54f97c797555760ae2f9ae46b3504d68fd14c5fc091962bd89"

			uri = fmt.Sprintf("https://gitee.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s", client_id, redirect_url, uuid)
			http.SetCookie(w, &http.Cookie{
				Name:     "state", //用于防伪造校验
				Value:    token,
				Path:     "/usercenter/api/v1/third-platform/gitee-login",
				HttpOnly: true,
				MaxAge:   60 * 5})
		default:
			http.Error(w, "not support platform", http.StatusNotImplemented)
			return
		}
		http.Redirect(w, r, uri, http.StatusFound)
	}
}

func LoginByDingTalkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #DOC:https://open.dingtalk.com/document/orgapp/scan-qr-code-to-log-on-to-third-party-websites
		fmt.Println("LoginByDingTalkHandler:", r.Header)
		code, state := r.URL.Query().Get("code"), r.URL.Query().Get("state")
		fmt.Println("login in by ding ding talk", "code:", code, "state:", state)
		// fmt.Println("LoginByDingTalkHandler:", r.Header)
		// ctx := r.Context()
		// var req types.LoginByDingTalkReq
		// if err := httpx.Parse(r, &req); err != nil {}
	}
}

func LoginByQQHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// #DOC:https://open.dingtalk.com/document/orgapp/scan-qr-code-to-log-on-to-third-party-websites
		fmt.Println("LoginByQQHandler:", r.Header)
		// code, state := r.URL.Query().Get("code"), r.URL.Query().Get()
	}
}

/*
github oauth登录第三方回调
*/
func LoginByGithubHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		stateValue, _ := svcCtx.RedisClient.Get(r.Context(), tools.OauthStateCacheKey(state)).Result()

		// 校验cookie中的 state字段,防止伪造攻击
		if cookieStateValue, err := r.Cookie("state"); err != nil || cookieStateValue.Value != stateValue {
			fmt.Println("state is invalid,refuse request")
			http.Error(w, "state is invalid", http.StatusForbidden)
			return
		}
		if _, err := models.ParseStateToken(stateValue, svcCtx.Config.JwtAuth.AccessSecret); err != nil {
			// 校验下state是否过期
			fmt.Println("parse state token error ->", err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		fmt.Println("Login by github handler:", r.Header, code, state)
		if code == "" {
			http.Error(w, "code is empty", http.StatusBadRequest)
			return
		} else {
			// 获取access token
			// github_client_id := os.Getenv("GITHUBCLIENTID")
			// github_client_secret := os.Getenv("GITHUBCLIENTSECRET")
			github_client_id := "Ov23liUpzHCPVymdPW2R"
			github_client_secret := "9b966488d33ab61541b6bcbdbfca3fa9e6d01cdb"
			url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", github_client_id, github_client_secret, code)
			header := map[string]string{"Accept": "application/json"}
			resp, StatusCode, err := tools.HttpPost(url, header, nil)
			fmt.Println("get github access token ->", string(resp), " status code -> ", StatusCode, "er ->", err)
			if err != nil || StatusCode != 200 {
				// fmt.Println("get github access token error ,status code->", StatusCode, " err -> ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var response struct {
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
				ErrorURI         string `json:"error_uri"`
				AccessToken      string `json:"access_token"`
				Scope            string `json:"scope"`
				TokenType        string `json:"token_type"`
			}
			err = json.Unmarshal(resp, &response)
			if err != nil || response.Error != "" {
				fmt.Println("unmarshal github access token error ->", err, " resp -> ", response)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("github access token ->", response)

			// 获取GitHub用户信息
			url = "https://api.github.com/user"
			header = map[string]string{
				"Authorization": "Bearer " + response.AccessToken,
				"Accept":        "application/json",
			}
			resp, StatusCode, err = tools.HttpGet(url, nil, header)
			if err != nil || StatusCode != 200 {
				fmt.Println("get github user info error ,status code->", StatusCode, " err -> ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Println("github user info ->", string(resp))
			var githubUserInfo struct {
				AvatarUrl string `json:"avatar_url"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				Location  string `json:"location"`
				Id        int    `json:"id"`
			}
			err = json.Unmarshal(resp, &githubUserInfo)
			if err != nil {
				fmt.Println("unmarshal github user info error ->", err, " resp -> ", githubUserInfo)
			}
			id_str := strconv.FormatInt(int64(githubUserInfo.Id), 10)
			// 直接生成一个默认的用户名并插入到用户表
			user, _ := svcCtx.UserModel.CreateUserByThirdPlatform("github", id_str, svcCtx.DB)
			if user.IsBind {
				fmt.Println("has login by github user and bind local account -> user", user.Username)
				accessToken := models.GenToken(*user, svcCtx.Config.JwtAuth.AccessSecret)
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    accessToken,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   60 * 5 * 24})
				httpx.OkJsonCtx(r.Context(), w, &types.LoginResp{
					BaseResponse: types.BaseResponse{
						Code: 0,
						Msg:  "登录成功",
					},
					Data: types.UserInfoWithToken{
						AccessToken: accessToken,
						UserInfo: types.UserInfo{
							Avatar:       user.Avatar,
							Email:        user.Email,
							Phone:        user.Phone,
							Age:          user.Age,
							Gender:       user.Gender,
							IsSuperAdmin: user.IsSuperAdmin,
							Username:     user.Username,
						},
					},
				})
			} else {
				fmt.Println("has login by github user but not bind local account")
				// 返回绑定的链接地址,通过state来防止伪造请求
				redirect_bind_url := fmt.Sprintf("/usercenter/third-platform/bind/%d?state=%s", user.ID, state)
				http.Redirect(w, r, redirect_bind_url, http.StatusFound)
			}
		}
	}
}

/*
qq oauth登录第三方回调
*/
/*
gitee oauth登录第三方回调
*/
func LoginByGiteeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		stateValue, _ := svcCtx.RedisClient.Get(r.Context(), tools.OauthStateCacheKey(state)).Result()

		// 校验cookie中的 state字段,防止伪造攻击
		if cookieStateValue, err := r.Cookie("state"); err != nil || cookieStateValue.Value != stateValue {
			fmt.Println("state is invalid,refuse request")
			http.Error(w, "state is invalid", http.StatusForbidden)
			return
		}
		if _, err := models.ParseStateToken(stateValue, svcCtx.Config.JwtAuth.AccessSecret); err != nil {
			// 校验下state是否过期
			fmt.Println("parse state token error ->", err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		fmt.Println("login by gitee code ->", code, "state ->", state)
		if code == "" {
			http.Error(w, "code is empty", http.StatusBadRequest)
			return
		} else {
			// 获取access token
			client_id := os.Getenv("GITEECLIENTID")
			client_secret := os.Getenv("GITEECLIENTSECRET")
			redirect_url := os.Getenv("GITEEREDIRECTURL")
			// redirect_url := "http://127.0.0.1:8080/usercenter/api/v1/third-platform/gitee-login"
			// client_id := "c45c44c778f0ad54f97c797555760ae2f9ae46b3504d68fd14c5fc091962bd89"
			// client_secret := "74e4760f72e1e4ce100d60a2cbaa0515b8b6c4e487b7dd9a1387272cf9e192cc"
			url := fmt.Sprintf("https://gitee.com/oauth/token?grant_type=authorization_code&code=%s&client_id=%s&redirect_uri=%s&client_secret=%s", code, client_id, redirect_url, client_secret)
			header := map[string]string{"Accept": "application/json"}
			resp, StatusCode, err := tools.HttpPost(url, header, nil)
			fmt.Println("get gitee access token ->", string(resp), " status code -> ", StatusCode, "er ->", err)
			if err != nil || StatusCode != 200 {
				// fmt.Println("get github access token error ,status code->", StatusCode, " err -> ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var response struct {
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
				ErrorURI         string `json:"error_uri"`
				AccessToken      string `json:"access_token"`
				RefreshToken     string `json:"refresh_token"`
				ExpiresIn        int    `json:"expires_in"`
				Scope            string `json:"scope"`
				TokenType        string `json:"token_type"`
			}
			err = json.Unmarshal(resp, &response)
			if err != nil || response.Error != "" {
				fmt.Println("unmarshal github access token error ->", err, " resp -> ", response)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("get gitee access token ->", response)
			// 获取GitEE用户信息
			url = "https://gitee.com/api/v5/user"
			header = map[string]string{
				// "Authorization": response.TokenType + " " + response.AccessToken,
				"Accept": "application/json",
			}
			resp, StatusCode, err = tools.HttpGet(url, map[string]string{"access_token": response.AccessToken}, header)
			if err != nil || StatusCode != 200 {
				fmt.Println("get gitee user info error ,status code->", StatusCode, " err -> ", err, " resp -> ", string(resp))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("gitee user info ->", string(resp))
			var giteeUserInfo struct {
				AvatarUrl string `json:"avatar_url"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				// Location  string `json:"location"`
				Id int `json:"id"`
			}
			err = json.Unmarshal(resp, &giteeUserInfo)
			if err != nil {
				fmt.Println("unmarshal gitee user info error ->", err, " resp -> ", giteeUserInfo)
			}
			id_str := strconv.FormatInt(int64(giteeUserInfo.Id), 10)
			// 直接生成一个默认的用户名并插入到用户表
			user, _ := svcCtx.UserModel.CreateUserByThirdPlatform("gitee", id_str, svcCtx.DB)
			if user.IsBind {
				fmt.Println("has login by gitee user and bind local account -> user", user.Username)
				accessToken := models.GenToken(*user, svcCtx.Config.JwtAuth.AccessSecret)
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    accessToken,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   60 * 5 * 24})
				httpx.OkJsonCtx(r.Context(), w, &types.LoginResp{
					BaseResponse: types.BaseResponse{
						Code: 0,
						Msg:  "登录成功",
					},
					Data: types.UserInfoWithToken{
						AccessToken: accessToken,
						UserInfo: types.UserInfo{
							Avatar:       user.Avatar,
							Email:        user.Email,
							Phone:        user.Phone,
							Age:          user.Age,
							Gender:       user.Gender,
							IsSuperAdmin: user.IsSuperAdmin,
							Username:     user.Username,
						},
					},
				})
			} else {
				fmt.Println("has login by gitee user but not bind local account")
				// 返回绑定的链接地址,通过state来防止伪造请求
				redirect_bind_url := fmt.Sprintf("/usercenter/third-platform/bind/%d?state=%s", user.ID, state)
				http.Redirect(w, r, redirect_bind_url, http.StatusFound)
			}
		}
	}
}

/*
绑定第三方平台账号到本地账号
*/
func BindLocalAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BindLocalAccountReq
		fmt.Println("bind third account to local account", req)
		if err := httpx.Parse(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 校验下 state,防止伪造请求
		if cookie, err := r.Cookie("state"); err != nil || cookie.Value != req.State {
			http.Error(w, "state is invalid", http.StatusBadRequest)
			return
		}

		// 绑定账号到现有的账号
		// 1。先查询账号名称或者邮箱是否存在
		user, err := svcCtx.UserModel.Create(req.UserName, req.Email, req.Password, svcCtx.DB)
		if err != nil {
			fmt.Println("bind user error ->", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			// 2。更新用户信息
			// user.ThirdPlatform = "github"
			// user.ThirdPlatformUserId = req.ThirdUserId
			user.Username = req.UserName
			user.Password = req.Password
			user.IsBind = true
			user.Email = req.Email
			fmt.Println("bind user success -> user", user.Username)
			accessToken := models.GenToken(*user, svcCtx.Config.JwtAuth.AccessSecret)
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    accessToken,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   60 * 5 * 24})
			httpx.OkJsonCtx(r.Context(), w, &types.LoginResp{
				BaseResponse: types.BaseResponse{
					Code: 0,
					Msg:  "登录成功",
				},
				Data: types.UserInfoWithToken{
					AccessToken: accessToken,
					UserInfo: types.UserInfo{
						Avatar:       user.Avatar,
						Email:        user.Email,
						Phone:        user.Phone,
						Age:          user.Age,
						Gender:       user.Gender,
						IsSuperAdmin: user.IsSuperAdmin,
						Username:     user.Username,
					},
				},
			})
		}
	}
}
