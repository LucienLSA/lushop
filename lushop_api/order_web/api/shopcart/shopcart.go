package shopCart

import (
	"context"
	"orderweb/api"
	"orderweb/forms"
	"orderweb/global"
	proto_goods "orderweb/proto/gen/goods"
	proto_inventory "orderweb/proto/gen/inventory"
	"strconv"

	"net/http"
	proto_order "orderweb/proto/gen/order"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 获取购物车商品
func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto_order.UserInfo{
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
		&proto_goods.BatchGoodsIdInfo{Id: ids})
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
	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				tempMap := map[string]interface{}{}
				tempMap["id"] = item.Id
				tempMap["goods_id"] = item.GoodsId
				tempMap["goods_name"] = good.Name
				tempMap["goods_image"] = good.GoodsFrontImage
				tempMap["goods_price"] = good.ShopPrice
				tempMap["nums"] = item.Nums
				tempMap["checked"] = item.Checked
				goodsList = append(goodsList, tempMap)
			}
		}
	}
	rspMap["data"] = goodsList
	ctx.JSON(http.StatusOK, rspMap)
}

// 添加商品到购物车
func New(ctx *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	// 添加商品到购物车之前，检查商品信息是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto_goods.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//如果添加到购物车的数量和库存的数量不一致
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto_inventory.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto_order.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	rspMap := make(map[string]interface{}, 0)
	rspMap["id"] = rsp.Id
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "参数有误",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto_order.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

// 更新购物车记录
func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "参数请求格式出错",
		})
		return
	}
	itemForm := forms.ShopCartItemUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	userId, _ := ctx.Get("userId")
	req := proto_order.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		req.Checked = *itemForm.Checked
	}
	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &req)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
