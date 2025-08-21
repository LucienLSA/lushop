package initialize

import (
	"lushopapi/api/oauth2"
	"lushopapi/global"
	"time"

	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v5"
)

func OAuth2() {
	global.Mgr = manage.NewDefaultManager()
	global.Mgr.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * time.Duration(global.ServerConfig.Oauth2Info.AccessTokenExp),
		RefreshTokenExp:   time.Hour * time.Duration(global.ServerConfig.Oauth2Info.RefreshTokenExp),
		IsGenerateRefresh: true,
	})
	global.Mgr.MustTokenStorage(store.NewMemoryTokenStore())
	global.Mgr.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(global.ServerConfig.Oauth2Info.JwtSignedKey), jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	for _, v := range global.ServerConfig.Oauth2Info.Client {
		clientStore.Set(v.Id, &models.Client{
			ID:     v.Id,
			Secret: v.Secret,
			Domain: v.Domain,
			Public: false,
		})
	}
	global.Mgr.MapClientStorage(clientStore)

	global.Srv = server.NewServer(server.NewConfig(), global.Mgr)
	global.Srv.SetPasswordAuthorizationHandler(oauth2.PasswordAuthorizationHandler) // 密码认证模式
	global.Srv.SetUserAuthorizationHandler(oauth2.UserAuthorizeHandler)
	global.Srv.SetAuthorizeScopeHandler(oauth2.AuthorizeScopeHandler) // 授权码模式
	global.Srv.SetInternalErrorHandler(oauth2.InternalErrorHandler)
	global.Srv.SetResponseErrorHandler(oauth2.ResponseErrorHandler)
}
