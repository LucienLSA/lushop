package device

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetDeviceFingerprint 生成设备唯一标识
func GetDeviceFingerprint(c *gin.Context) string {
	// 收集设备特征信息
	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()
	acceptLanguage := c.GetHeader("Accept-Language")
	acceptEncoding := c.GetHeader("Accept-Encoding")

	// 组合特征字符串
	featureStr := strings.Join([]string{userAgent, ip, acceptLanguage, acceptEncoding}, "|")

	// SHA256哈希生成唯一标识
	h := sha256.New()
	h.Write([]byte(featureStr))
	return hex.EncodeToString(h.Sum(nil))
}
