package pay

import (
	"context"
	"net/http"
	"orderweb/global"
	proto_order "orderweb/proto/gen/order"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

// 支付宝的回调通知
func Notify(ctx *gin.Context) {
	// 验证是否为支付宝发送
	appID := global.GetEnvInfoStr(global.ServerConfig.AlipayInfo.AppID)
	privateKey := global.GetEnvInfoStr(global.ServerConfig.AlipayInfo.PrivateKey)
	aliPublicKey := global.GetEnvInfoStr(global.ServerConfig.AlipayInfo.AliPublicKey)
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		zap.S().Errorw("生成支付宝URL失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	noti, err := client.DecodeNotification(ctx.Request.Form)
	if err != nil {
		zap.S().Errorw("调用 VerifySign 方法验证签名失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 业务处理
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto_order.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "不合法的通知",
		})
		return
	}
	// 如果通知消息没有问题，我们需要确认收到通知消息，不然支付宝后续会继续推送相同的消息
	// alipay.ACKNotification(writer)
	ctx.String(http.StatusOK, "success")
}
