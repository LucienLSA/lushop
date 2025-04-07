package order

import (
	"context"
	"net/http"
	"orderweb/api"
	"orderweb/global"
	proto_order "orderweb/proto/gen/order"
	"orderweb/utils/jwtClaims"
	"strconv"

	"github.com/gin-gonic/gin"
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

}

func Detail(ctx *gin.Context) {

}
