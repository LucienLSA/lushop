package handler

import (
	"context"
	"fmt"
	"lushopsrvs/goods_srv/global"
	"lushopsrvs/goods_srv/model"
	"lushopsrvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

// 商品接口
// 关键词搜索、查询新品、查询热门商品、通过价格区间筛选、通过商品分类筛选
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 各种条件查询
	goodsListRsp := &proto.GoodsListResponse{}
	var goods []model.Goods
	// localDB := global.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	localDB := global.DB.Model(model.Goods{}).Session(&gorm.Session{})
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.IsHot {
		// localDB.Where("is_hot=true")
		localDB = localDB.Where(model.Goods{IsHot: true})
	}
	if req.IsNew {
		localDB = localDB.Where("is_new=true")
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id=?", req.Brand)
	}
	fmt.Println(req)
	// 通过category查询商品
	if req.TopCategory > 0 {
		var category model.Category
		result := global.DB.Model(&model.Category{}).First(&category, req.TopCategory)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		// subQuery := global.DB.Model(&model.Category{})
		if category.Level == 1 {
			// sqlQuery = `SELECT id FROM category WHERE parent_category_id IN (SELECT id FROM category WHERE parent_category_id
			// 	=?)`
			// 查询二级分类ID（一级分类的子分类）
			// subQuery = subQuery.Where(
			// 	"parent_category_id IN (?)",
			// 	global.DB.Model(&model.Category{}).Select("id").Where("parent_category_id = ?", req.TopCategory))
			localDB = localDB.Joins(
				"JOIN category c3 ON goods.category_id = c3.id "+
					"JOIN category c2 ON c3.parent_category_id = c2.id "+
					"JOIN category c1 ON c2.parent_category_id = c1.id "+
					"WHERE c1.id = ?", req.TopCategory)
		} else if category.Level == 2 {
			// sqlQuery = `SELECT id FROM category WHERE parent_category_id =?`
			// 查询三级分类ID（二级分类的子分类）
			// subQuery = subQuery.Where("parent_category_id = ?", req.TopCategory)
			localDB = localDB.Joins(
				"JOIN category c3 ON goods.category_id = c3.id "+
					"JOIN category c2 ON c3.parent_category_id = c2.id "+
					"WHERE c2.id = ?", req.TopCategory)
		} else if category.Level == 3 {
			// sqlQuery = `SELECT id FROM category WHERE id =?`
			// subQuery = subQuery.Where("id = ?", req.TopCategory)
			localDB = localDB.Where("goods.category_id = ?", req.TopCategory)
		}
		// localDB = localDB.Where("category_id IN (?)", subQuery.Select("id"))
	}
	var count int64
	if err := localDB.Count(&count).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "获取商品总数失败")
	}
	goodsListRsp.Total = int32(count)
	fmt.Println(count)
	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, good := range goods {
		goodsInfoRsp := ModelToResponse(good)
		goodsListRsp.Data = append(goodsListRsp.Data, &goodsInfoRsp)
	}
	return goodsListRsp, nil
}

// 现在用户提交订单有多个商品，批量查询商品的信息
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListRsp := &proto.GoodsListResponse{}
	var goods []model.Goods
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoRsp := ModelToResponse(good)
		goodsListRsp.Data = append(goodsListRsp.Data, &goodsInfoRsp)
	}
	goodsListRsp.Total = int32(result.RowsAffected)
	return goodsListRsp, nil
}

// func (s *GoodsServer) CreateGoods(ctx context.Context,req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
// }
// func (s *GoodsServer) DeleteGoods(ctx context.Context,req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {}
// func (s *GoodsServer) UpdateGoods(ctx context.Context,req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods
	result := global.DB.First(&goods, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoRsp := ModelToResponse(goods)
	return &goodsInfoRsp, nil
}
