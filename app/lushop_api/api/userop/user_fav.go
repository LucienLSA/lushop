package userop

import (
	"context"
	"lushopapi/api/base"
	v2base "lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	v2goodsproto "lushopapi/proto/goods"
	v2useropproto "lushopapi/proto/userop"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 收藏列表
func FavList(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	userFavRsp, err := global.UserOpSrvClient.GetFavList(context.Background(), &v2useropproto.UserFavRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("获取收藏列表失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range userFavRsp.Data {
		ids = append(ids, item.GoodsId)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	//请求商品服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &v2goodsproto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": userFavRsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range userFavRsp.Data {
		data := gin.H{
			"id": item.GoodsId,
		}

		for _, good := range goods.Data {
			if item.GoodsId == good.Id {
				data["name"] = good.Name
				data["shop_price"] = good.ShopPrice
			}
		}

		goodsList = append(goodsList, data)
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

// 新建收藏
func FavCreate(ctx *gin.Context) {
	userFavForm := forms.UserFavForm{}
	if err := ctx.ShouldBindJSON(&userFavForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err := global.UserOpSrvClient.AddUserFav(context.Background(), &v2useropproto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: userFavForm.GoodsId,
	})

	if err != nil {
		zap.S().Errorw("添加收藏记录失败")
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "添加收藏记录成功",
	})
}

// 删除收藏
func FavDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserOpSrvClient.DeleteUserFav(context.Background(), &v2useropproto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除收藏记录失败")
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// 收藏详情
func FavDetail(ctx *gin.Context) {
	goodsId := ctx.Param("id")
	goodsIdInt, err := strconv.ParseInt(goodsId, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	userId, _ := ctx.Get("userId")

	// 查询收藏状态
	favRsp, err := global.UserOpSrvClient.GetUserFavDetail(context.Background(), &v2useropproto.UserFavRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
	})
	if err != nil {
		zap.S().Errorw("查询收藏状态失败") //未收藏
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 查询商品详情
	goodsRsp, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &v2goodsproto.GoodInfoRequest{
		Id: int32(goodsIdInt),
	})
	if err != nil {
		zap.S().Errorw("查询商品详情失败")
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 组合返回数据
	result := gin.H{
		"favorite": gin.H{
			"user_id":  favRsp.UserId,
			"goods_id": favRsp.GoodsId,
		},
		"goods": gin.H{
			"id":                goodsRsp.Id,
			"name":              goodsRsp.Name,
			"goods_sn":          goodsRsp.GoodsSn,
			"click_num":         goodsRsp.ClickNum,
			"sold_num":          goodsRsp.SoldNum,
			"fav_num":           goodsRsp.FavNum,
			"market_price":      goodsRsp.MarketPrice,
			"shop_price":        goodsRsp.ShopPrice,
			"goods_brief":       goodsRsp.GoodsBrief,
			"ship_free":         goodsRsp.ShipFree,
			"images":            goodsRsp.Images,
			"desc_images":       goodsRsp.DescImages,
			"goods_front_image": goodsRsp.GoodsFrontImage,
			"is_new":            goodsRsp.IsNew,
			"is_hot":            goodsRsp.IsHot,
			"on_sale":           goodsRsp.OnSale,
			"category": gin.H{
				"id":   goodsRsp.Category.Id,
				"name": goodsRsp.Category.Name,
			},
			"brand": gin.H{
				"id":   goodsRsp.Brand.Id,
				"name": goodsRsp.Brand.Name,
				"logo": goodsRsp.Brand.Logo,
			},
		},
	}

	ctx.JSON(http.StatusOK, result)
}
