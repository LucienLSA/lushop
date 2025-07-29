package order

import (
	"context"
	"lushopapi/global"
	v2orderproto "lushopapi/proto/order"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

// 支付宝回调通知
func Notify(ctx *gin.Context) {
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("[alipay] 实例化【支付宝的url】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("[alipay] 加载【支付宝的公钥】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		zap.S().Errorw("[alipay] VerifySign 方法验证签名失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &v2orderproto.OrderStatus{
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
	//alipay.AckNotification(rep)
	ctx.String(http.StatusOK, "success")
}
