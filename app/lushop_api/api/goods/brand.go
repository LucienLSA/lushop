package goods

import (
	"context"
	"lushopapi/api/base"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lushopapi/forms"
	"lushopapi/global"
	v2goodsproto "lushopapi/proto/goods"
)

func BrandList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.GoodsSrvClient.BrandList(context.Background(), &v2goodsproto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})

	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	reMap := make(map[string]interface{})
	reMap["total"] = rsp.Total
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func BrandCreate(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateBrand(context.Background(), &v2goodsproto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["logo"] = rsp.Logo

	ctx.JSON(http.StatusOK, request)
}

func BrandDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteBrand(context.Background(), &v2goodsproto.BrandRequest{Id: int32(i)})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

// 更新品牌
func BrandUpdate(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBrand(context.Background(), &v2goodsproto.BrandRequest{
		Id:   int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

// 获取品牌分类下的品牌
func CateGetBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &v2goodsproto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

// 更新品牌分类
func CateBrandUpdate(ctx *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id") // 商品品牌id
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &v2goodsproto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

// 删除分类下的品牌
func CateBrandDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &v2goodsproto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// 新建品牌分类
func CateBrandCreate(ctx *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategoryBrand(context.Background(), &v2goodsproto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id

	ctx.JSON(http.StatusOK, response)
}

// 获取品牌分类列表
func CateBrandList(ctx *gin.Context) {
	//所有的list返回的数据结构
	/*
		{
			"total": 100,
			"data":[{},{}]
		}
	*/
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)

	rsp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &v2goodsproto.CategoryBrandFilterRequest{
		Pages:       int32(pagesInt),
		PagePerNums: int32(perNumsInt),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["category"] = map[string]interface{}{
			"id":   value.Category.Id,
			"name": value.Category.Name,
		}
		reMap["brand"] = map[string]interface{}{
			"id":   value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		}

		result = append(result, reMap)
	}

	reMap["data"] = result
	ctx.JSON(http.StatusOK, reMap)
}
