package goods

import (
	"context"
	"encoding/json"
	v2base "lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	v2goodsproto "lushopapi/proto/goods"
	v2inventoryproto "lushopapi/proto/inventory"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var goodsDetailGroup singleflight.Group

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 商品的列表 pmin=abc
func GoodsList(ctx *gin.Context) {
	//fmt.Println("商品列表")
	//商品的列表 pmin=abc, spring cloud, go-micro
	request := &v2goodsproto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)
	//parent, _ := ctx.Get("parentSpan")
	//opentracing.ContextWithSpan(context.Background(), parent.(opentracing.Span))

	//对商品列表进行限流
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "请求过于频繁,请稍后重试"})
		return
	}
	r, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//只管list查询的逻辑限速
	e.Exit()

	reMap := map[string]interface{}{
		"total": r.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"ctegory": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

func GoodsListES(ctx *gin.Context) {
	//fmt.Println("商品列表")
	//商品的列表 pmin=abc, spring cloud, go-micro
	request := &v2goodsproto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)
	//parent, _ := ctx.Get("parentSpan")
	//opentracing.ContextWithSpan(context.Background(), parent.(opentracing.Span))

	//对商品列表进行限流
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "请求过于频繁,请稍后重试"})
		return
	}
	r, err := global.GoodsSrvClient.GoodsListES(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//只管list查询的逻辑限速
	e.Exit()

	reMap := map[string]interface{}{
		"total": r.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"ctegory": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

// 创建商品
func GoodsCreate(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &v2goodsproto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//如何设置库存
	// 设置初始库存（调用inventory_srv）
	inventoryClient := global.InventorySrvClient
	_, err = inventoryClient.SetInv(context.Background(), &v2inventoryproto.GoodsInvInfo{
		GoodsId: rsp.Id,
		Num:     goodsForm.Stocks,
	})
	if err != nil {
		// 库存设置失败时回滚商品创建（实际业务中可能需要更复杂的补偿机制）
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// 不在创建时预写详情缓存，交由详情接口统一处理缓存
	ctx.JSON(http.StatusOK, rsp)
}

// 商品详情
func GoodsDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	// 商品详情 - 热点缓存 + singleflight 请求合并
	reqCtx := ctx.Request.Context()
	cacheKey := "goods_detail_" + id

	// 1) 查缓存
	if cachedBytes, err := global.RedisClient.Get(reqCtx, cacheKey).Bytes(); err == nil {
		var cached v2goodsproto.GoodsInfoResponse
		if err := json.Unmarshal(cachedBytes, &cached); err == nil {
			ctx.JSON(http.StatusOK, &cached)
			return
		}
		zap.S().Warnf("缓存反序列化失败: %v", err)
	} else if err != redis.Nil {
		zap.S().Warnf("缓存查询失败: %v", err)
	}

	// 2) singleflight 合并请求
	v, err, _ := goodsDetailGroup.Do(cacheKey, func() (interface{}, error) {
		r, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &v2goodsproto.GoodInfoRequest{
			Id: int32(i),
		})
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(r)
		ttl := 24*time.Hour + time.Duration(rand.Intn(600))*time.Second
		if err := global.RedisClient.Set(reqCtx, cacheKey, data, ttl).Err(); err != nil {
			zap.S().Warnf("商品缓存设置失败: %v", err)
		}
		return r, nil
	})
	if err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	resp := v.(*v2goodsproto.GoodsInfoResponse)
	ctx.JSON(http.StatusOK, resp)

	// rsp := map[string]interface{}{
	// 	"id":          r.Id,
	// 	"name":        r.Name,
	// 	"goods_brief": r.GoodsBrief,
	// 	"desc":        r.GoodsDesc,
	// 	"ship_free":   r.ShipFree,
	// 	"images":      r.Images,
	// 	"desc_images": r.DescImages,
	// 	"front_image": r.GoodsFrontImage,
	// 	"shop_price":  r.ShopPrice,
	// 	"ctegory": map[string]interface{}{
	// 		"id":   r.Category.Id,
	// 		"name": r.Category.Name,
	// 	},
	// 	"brand": map[string]interface{}{
	// 		"id":   r.Brand.Id,
	// 		"name": r.Brand.Name,
	// 		"logo": r.Brand.Logo,
	// 	},
	// 	"is_hot":  r.IsHot,
	// 	"is_new":  r.IsNew,
	// 	"on_sale": r.OnSale,
	// }
	// ctx.JSON(http.StatusOK, rsp)
}

// 删除商品
func GoodsDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &v2goodsproto.DeleteGoodsInfo{
		Id: int32(i),
	})
	if err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
	}
	// 删除商品详情缓存
	cacheKey := "goods_detail_" + id
	global.RedisClient.Del(ctx.Request.Context(), cacheKey)
	ctx.Status(http.StatusOK)
}

// 查询商品的库存 库存用单独的url来请求   异步请求
func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	GoodsId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	var invInfo *v2inventoryproto.GoodsInvInfo
	if invInfo, err = global.InventorySrvClient.InvDetail(context.Background(), &v2inventoryproto.GoodsInvInfo{
		GoodsId: int32(GoodsId),
	}); err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"goodsId": invInfo.GoodsId, "num": invInfo.Num})
}

// 更新商品状态
func GoodsUpdateStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, _ := strconv.ParseInt(id, 10, 32)
	var err error
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &v2goodsproto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// 添加缓存删除逻辑
	cacheKey := "goods_detail_" + id
	global.RedisClient.Del(ctx.Request.Context(), cacheKey)
	ctx.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

// 更新商品信息
func GoodsUpdate(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}
	id := ctx.Param("id")
	i, _ := strconv.ParseInt(id, 10, 32)
	var err error
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &v2goodsproto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// 更新商品详情缓存（删除旧缓存）
	cacheKey := "goods_detail_" + id
	global.RedisClient.Del(ctx.Request.Context(), cacheKey)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}
