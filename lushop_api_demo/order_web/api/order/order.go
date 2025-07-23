package order

import (
	"context"
	"fmt"
	"net/http"
	"orderweb/api"
	"orderweb/forms"
	"orderweb/global"
	proto_order "orderweb/proto/gen/order"
	"orderweb/utils/jwtClaims"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	req := proto_order.OrderFilterRequest{}
	jwtC := claims.(*jwtClaims.CustomClaims)
	// 如果是管理员，则返回所有的订单
	if jwtC.AuthorityId == 1 {
		req.UserId = int32(userId.(uint))
	}
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	req.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PagePerNums = int32(perNumsInt)

	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &req)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	/*
		{
			"total":100,
			"data": [
				{
					"id":
					"status":
				}
			]

		}
	*/
	rspMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		tempMap := map[string]interface{}{}
		tempMap["id"] = item.Id
		tempMap["status"] = item.Status
		tempMap["pay_type"] = item.PayType
		tempMap["user"] = item.UserId
		tempMap["post"] = item.Post
		tempMap["address"] = item.Address
		tempMap["name"] = item.Name
		tempMap["mobile"] = item.Mobile
		tempMap["order_sn"] = item.OrderSn
		tempMap["id"] = item.Id
		tempMap["add_time"] = item.AddTime
		orderList = append(orderList, tempMap)
	}
	rspMap["data"] = orderList
	ctx.JSON(http.StatusOK, rspMap)
}

func New(ctx *gin.Context) {
	orderForms := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForms); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(context.Background(), &proto_order.OrderRequest{
		UserId:  int32(userId.(uint)),
		Address: orderForms.Address,
		Mobile:  orderForms.Mobile,
		Name:    orderForms.Name,
		Post:    orderForms.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// 生成支付宝订单url
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
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "lushop order:" + rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	p.ProductCode = global.ServerConfig.AlipayInfo.ProductCode

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付URL失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	fmt.Println(url.String())
	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "参数格式有误",
		})
		return
	}
	req := proto_order.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	jwtC := claims.(*jwtClaims.CustomClaims)
	// 如果是管理员，则返回所有的订单
	if jwtC.AuthorityId == 1 {
		req.UserId = int32(userId.(uint))
	}
	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &req)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	rspMap := gin.H{}
	rspMap["id"] = rsp.OrderInfo.Id
	rspMap["status"] = rsp.OrderInfo.Status
	rspMap["user"] = rsp.OrderInfo.UserId
	rspMap["post"] = rsp.OrderInfo.Post
	rspMap["total"] = rsp.OrderInfo.Total
	rspMap["add_time"] = rsp.OrderInfo.AddTime
	rspMap["address"] = rsp.OrderInfo.Address
	rspMap["name"] = rsp.OrderInfo.Name
	rspMap["mobile"] = rsp.OrderInfo.Mobile
	rspMap["pay_type"] = rsp.OrderInfo.PayType
	rspMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"price": item.GoodsPrice,
			"image": item.GoodsImage,
			"nums":  item.Nums,
		}
		goodsList = append(goodsList, tmpMap)
	}
	rspMap["goods"] = goodsList

	// 生成支付宝url
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
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "lushop order:" + rsp.OrderInfo.OrderSn
	p.OutTradeNo = rsp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64)
	p.ProductCode = global.ServerConfig.AlipayInfo.ProductCode

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付URL失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	fmt.Println(url.String())
	rspMap["alipay_url"] = url.String()
	ctx.JSON(http.StatusOK, rspMap)
}
