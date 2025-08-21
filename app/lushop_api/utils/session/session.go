package session

import (
	"lushopapi/global"
	"net/http"

	"github.com/gorilla/sessions"
)

// ========== 原生 HTTP 版本的方法 ==========

// GetSessionRaw 获取原生HTTP会话
func GetSessionRaw(r *http.Request) (*sessions.Session, error) {
	return global.Redistore.Get(r, global.ServerConfig.SessionInfo.Name)
}

// GetRaw 从原生HTTP会话中获取值
func GetRaw(r *http.Request, key string) (interface{}, error) {
	session, err := GetSessionRaw(r)
	if err != nil {
		return nil, err
	}
	if value, ok := session.Values[key]; ok {
		return value, nil
	}
	return nil, nil
}

// SetRaw 设置原生HTTP会话值
func SetRaw(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, err := GetSessionRaw(r)
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(r, w)
}

// DeleteRaw 删除原生HTTP会话中的值
func DeleteRaw(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := GetSessionRaw(r)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	return session.Save(r, w)
}

// ClearRaw 清除原生HTTP会话
func ClearRaw(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSessionRaw(r)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1 // 设置为-1使cookie立即过期
	return session.Save(r, w)
}

// GetStringRaw 获取原生HTTP会话中的字符串值
func GetStringRaw(r *http.Request, key string) (string, error) {
	value, err := GetRaw(r, key)
	if err != nil {
		return "", err
	}
	if str, ok := value.(string); ok {
		return str, nil
	}
	return "", nil
}

// GetIntRaw 获取原生HTTP会话中的整数值
func GetIntRaw(r *http.Request, key string) (int, error) {
	value, err := GetRaw(r, key)
	if err != nil {
		return 0, err
	}
	if num, ok := value.(int); ok {
		return num, nil
	}
	return 0, nil
}

// GetBoolRaw 获取原生HTTP会话中的布尔值
func GetBoolRaw(r *http.Request, key string) (bool, error) {
	value, err := GetRaw(r, key)
	if err != nil {
		return false, err
	}
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, nil
}
