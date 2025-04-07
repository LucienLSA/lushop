package shopcart

import (
	"context"
	"lushopapi/order_web/api"
	"lushopapi/order_web/global"
	"lushopapi/order_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 获取购物车商品
func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}
	// 请求商品服务获取商品信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(),
		&proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		zap.S().Errorw("[List] 批量查询 【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	rspMap := gin.H{
		"total": rsp.Total,
	}
	/*
		{
			"total":12,
			"data":[
				{
					"id":1,
					"goods_id":421,
					"goods_price":
					"goods_name":
					"goods_image":
					"nums":
					"checked":
				}
			]
		}
	*/
	goodsList := make([]interface{})
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				tempMap := map[string]interface{}{}
				tempMap["id"] = item.Id
				tempMap["goods_id"] = item.GoodsId
				tempMap["goods_name"] = good.Name
				tempMap["goods_image"] = good.GoodsFrontImage
				tempMap["goods_price"] = item.S
				tempMap["nums"] = item.Nums
				tempMap["checked"] = item.Checked
			}
		}
	}
}

func New(ctx *gin.Context) {

}

func Delete(ctx *gin.Context) {

}

func Update(ctx *gin.Context) {

}
