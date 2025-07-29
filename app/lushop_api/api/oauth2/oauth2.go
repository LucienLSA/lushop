package oauth2

import (
	"context"
	"html/template"
	"log"
	"lushopapi/config"
	"lushopapi/global"
	"lushopapi/utils/session"
	"net/http"
	"net/url"
	"time"

	v2userproto "lushopapi/proto/user"
	utilsoauth2 "lushopapi/utils/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authentication 通过 gRPC 用户服务校验用户名和密码
func Authentication(ctx context.Context, username, password string) (userID string, err error) {
	// 1. 通过手机号获取用户信息
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &v2userproto.MobileRequest{
		Mobile: username,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return "", errors.New("用户不存在")
			default:
				return "", errors.New("登录失败,内部错误")
			}
		}
	}

	//2. 上面只是查询到用户了而已，并没有检查密码
	if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &v2userproto.PasswordCheckInfo{
		PassWord:          password,     // 前端用户传入的密码
		EncryptedPassWord: rsp.PassWord, // 数据库中查询到的用户设置的密码
	}); pasErr != nil || !passRsp.Success {
		return "", errors.New("用户名或密码错误")
	}
	// 3. 返回用户ID
	return rsp.Mobile, nil
}

// 密码模式认证
func PasswordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {

	userID, err = Authentication(ctx, username, password)
	return
}

// 用户授权处理 (原生HTTP版本)
// 注意: 此函数签名必须是 (http.ResponseWriter, *http.Request) 以匹配 go-oauth2 的要求.
// 你的 session 包需要提供支持原生 http 的方法, 例如 GetRaw/SetRaw.
func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	// 假设 session 包提供了 GetRaw 方法
	v, _ := session.GetRaw(r, "LoggedInUserID")
	if v == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		// 假设 session 包提供了 SetRaw 方法
		session.SetRaw(w, r, "RequestForm", r.Form)
		log.Printf("UserAuthorizeHandler: set session.RequestForm = %#v", r.Form)

		// 重定向到登录页面
		w.Header().Set("Location", "/oauth2/login")
		http.Redirect(w, r, "/oauth2/login", http.StatusFound)
		return
	}
	userID, _ = v.(string)
	return
}

// scope处理 (原生HTTP版本)
// 注意: 此函数签名必须是 (http.ResponseWriter, *http.Request) 以匹配 go-oauth2 的要求.
func AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	s := utilsoauth2.ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		err = errors.New("无效的权限范围")
		return
	}
	scope = utilsoauth2.ScopeJoin(s)
	return
}

// 内部错误处理
func InternalErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
	return
}

// 响应错误处理
func ResponseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}

// 授权页面
func AuthorizeHandler(c *gin.Context) {
	var form url.Values
	if v, _ := session.GetRaw(c.Request, "RequestForm"); v != nil {
		c.Request.ParseForm()
		if c.Request.Form.Get("client_id") == "" {
			form = v.(url.Values)
		}
	}
	c.Request.Form = form

	if err := session.DeleteRaw(c.Writer, c.Request, "RequestForm"); err != nil {
		ErrorHandler(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := global.Srv.HandleAuthorizeRequest(c.Writer, c.Request); err != nil {
		ErrorHandler(c, err.Error(), http.StatusBadRequest)
		return
	}
}

// 登录页面
func LoginHandler(c *gin.Context) {
	form, _ := session.GetRaw(c.Request, "RequestForm")
	log.Printf("LoginHandler: session.RequestForm = %#v", form)
	if form == nil {
		ErrorHandler(c, "无效的请求", http.StatusInternalServerError)
		return
	}

	clientID := form.(url.Values).Get("client_id")
	scope := form.(url.Values).Get("scope")

	client := utilsoauth2.GetOAuth2Client(clientID)

	data := struct {
		Client config.Oauth2ClientConfig
		Scope  []config.Oauth2ScopeConfig
		Error  string
	}{
		Client: *client,
		Scope:  utilsoauth2.ScopeFilter(clientID, scope),
	}

	if data.Scope == nil {
		ErrorHandler(c, "无效的权限范围", http.StatusBadRequest)
		return
	}

	if c.Request.Method == "POST" {
		var userID string
		var err error

		if c.Request.Form == nil {
			err = c.Request.ParseForm()
			if err != nil {
				ErrorHandler(c, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if c.Request.Form.Get("type") == "password" {

			userID, err = Authentication(c.Request.Context(), c.Request.Form.Get("username"), c.Request.Form.Get("password"))
			if err != nil {
				data.Error = err.Error()
				t, _ := template.ParseFiles("template/login.html")
				t.Execute(c.Writer, data)
				return
			}
		}
		err = session.SetRaw(c.Writer, c.Request, "LoggedInUserID", userID)
		if err != nil {
			ErrorHandler(c, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Header("Location", "/authorize")
		c.Status(http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("template/login.html")
	t.Execute(c.Writer, data)
}

// 退出登录
func LogoutHandler(c *gin.Context) {
	if c.Request.Form == nil {
		if err := c.Request.ParseForm(); err != nil {
			ErrorHandler(c, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 检查redirect_uri参数
	redirectURI := c.Request.Form.Get("redirect_uri")
	if redirectURI == "" {
		ErrorHandler(c, "参数不能为空(redirect_uri)", http.StatusBadRequest)
		return
	}
	if _, err := url.Parse(redirectURI); err != nil {
		ErrorHandler(c, "参数无效(redirect_uri)", http.StatusBadRequest)
		return
	}

	// 删除公共回话
	if err := session.DeleteRaw(c.Writer, c.Request, "LoggedInUserID"); err != nil {
		ErrorHandler(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Header("Location", redirectURI)
	c.Status(http.StatusFound)
}

// token接口
func TokenHandler(c *gin.Context) {
	err := global.Srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// token校验
func VerifyHandler(c *gin.Context) {
	token, err := global.Srv.ValidationBearerToken(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cli, err := global.Mgr.GetClient(c.Request.Context(), token.GetClientID())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := gin.H{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		"domain":     cli.GetDomain(),
	}
	c.JSON(http.StatusOK, data)
}

// 404
func NotFoundHandler(c *gin.Context) {
	ErrorHandler(c, "无效的地址", http.StatusNotFound)
}

// 错误页面
func ErrorHandler(c *gin.Context, message string, status int) {
	c.Status(status)
	if status >= 400 {
		t, _ := template.ParseFiles("template/error.html")
		body := struct {
			Status  int
			Message string
		}{Status: status, Message: message}
		t.Execute(c.Writer, body)
	}
}
