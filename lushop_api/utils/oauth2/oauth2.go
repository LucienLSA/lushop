package oauth2

import (
	"lushopapi/config"
	"lushopapi/global"
	"strings"
)

// 获取指定 clientID 的 OAuth2 客户端配置
func GetOAuth2Client(clientID string) *config.Oauth2ClientConfig {
	for _, v := range global.ServerConfig.Oauth2Info.Client {
		if v.Id == clientID {
			return &v
		}
	}
	return nil
}

// 将 scope 列表拼接为逗号分隔字符串
func ScopeJoin(scope []config.Oauth2ScopeConfig) string {
	var s []string
	for _, sc := range scope {
		s = append(s, sc.Id)
	}
	return strings.Join(s, ",")
}

// 根据 clientID 和 scope 字符串过滤出有效的 scope 配置
func ScopeFilter(clientID string, scope string) []config.Oauth2ScopeConfig {
	client := GetOAuth2Client(clientID)

	sl := strings.Split(scope, ",")
	var result []config.Oauth2ScopeConfig
	for _, str := range sl {
		for _, sc := range client.Scope {
			if str == sc.Id {
				result = append(result, sc)
			}
		}
	}
	return result
}
