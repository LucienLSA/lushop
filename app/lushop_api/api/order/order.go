package order

import (
	"context"
	v2base "lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	v2orderproto "lushopapi/proto/order"
	"lushopapi/utils/jwtClaims"
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

// 订单列表
func OrderList(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := v2orderproto.OrderFilterRequest{}
	//如果是管理员用户则返回所有的订单
	model := claims.(*jwtClaims.CustomClaims)
	if model.AuthorityId != 1 {
		request.UserId = int32(userId.(uint))
	}

	pages := ctx.DefaultQuery("p", "1")
	pagesInt, _ := strconv.Atoi(pages)
	//request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "10")
	perNumsInt, _ := strconv.Atoi(perNums)
	//request.PagePerNums = int32(perNumsInt)

	request.Pages = int32(pagesInt)
	request.PagePerNums = int32(perNumsInt)

	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("[OrderList] 获取【订单列表】失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	/*
		{
			"total":100,
			"data":[
				{
					"id":123,
					"id":123,
				}
			]
		}
	*/
	reMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		//tmpMap["add_time"] = item.OrderSn
		tmpMap["order_sn"] = item.OrderSn

		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)

}

// 创建订单
func OrderCreate(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}
	userId, _ := ctx.Get("userId")

	e, b := sentinel.Entry("create_order", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		zap.S().Errorw("[CreateOrder] 被限流了")
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于繁忙，请稍后重试",
		})
		return
	}

	zap.S().Debug("执行grpc层CreateOrder")
	rsp, err := global.OrderSrvClient.CreateOrder(context.WithValue(context.Background(), "ginContext", ctx), &v2orderproto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("[CreateOrder] 新建【订单】失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	zap.S().Debug("grpc层CreateOrder执行结束")
	//退出限流
	e.Exit()

	//生成支付宝的支付url
	appID := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.AppID)
	privateKey := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.PrivateKey)
	aliPublicKey := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.AliPublicKey)
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		zap.S().Errorw("[alipay] 实例化【支付宝的url】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		zap.S().Errorw("[alipay] 加载【支付宝的公钥】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	//PagePay 网页支付
	//WapPay 手机支付
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL
	p.Subject = "lushop订单 - " + rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	//p.ProductCode = "QUICK_WAP_WAY" 手机支付   ↓网页支付↓
	p.ProductCode = global.ServerConfig.AliPayInfo.ProductCode

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("[alipay] 生成【支付url】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		// "id":         rsp.Id,
		"order_sn":   rsp.OrderSn,
		"alipay_url": url.String(),
	})

}

// 订单详情
func OrderDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}
	//如果是管理员用户则返回所有的订单
	request := v2orderproto.OrderRequest{
		Id: int32(i),
	}
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*jwtClaims.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}
	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("[OrderDeteil] 获取【订单详情】失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tempMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}
		goodsList = append(goodsList, tempMap)
	}

	reMap["goods"] = goodsList

	//生成支付宝的支付url
	appID := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.AppID)
	privateKey := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.PrivateKey)
	aliPublicKey := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.AliPublicKey)
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		zap.S().Errorw("[alipay] 实例化【支付宝的url】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		zap.S().Errorw("[alipay] 加载【支付宝的公钥】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	//PagePay 网页支付
	//WapPay 手机支付
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL
	p.Subject = "lushop订单-" + rsp.OrderInfo.OrderSn
	p.OutTradeNo = rsp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64)
	//p.ProductCode = "QUICK_WAP_WAY" 手机支付   ↓网页支付↓
	p.ProductCode = global.ServerConfig.AliPayInfo.ProductCode

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("[alipay] 生成【支付url】失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	reMap["alipay_url"] = url.String()

	ctx.JSON(http.StatusOK, reMap)
}
