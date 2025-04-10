package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goodssrv/global"
	"goodssrv/model"
	"goodssrv/proto"

	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}

// 商品接口
// 关键词搜索、查询新品、查询热门商品、通过价格区间筛选、通过商品分类筛选
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 各种条件查询
	goodsListRsp := &proto.GoodsListResponse{}
	// match bool 复合查询
	q := elastic.NewBoolQuery()
	// localDB := global.DB.Model(&model.Goods{}).Session(&gorm.Session{SkipDefaultTransaction: true})
	// localDB := global.DB.Model(&model.Goods{}).Session(&gorm.Session{})
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		// localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		// localDB.Where("is_hot=true")
		// localDB = localDB.Where(model.Goods{IsHot: true})
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		// localDB = localDB.Where("is_new=true")
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.PriceMin > 0 {
		// localDB = localDB.Where("shop_price >= ?", req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		// localDB = localDB.Where("shop_price <= ?", req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		// localDB = localDB.Where("brand_id=?", req.Brand)
		q = q.Filter(elastic.NewTermQuery("brand_id", req.Brand))
	}
	fmt.Println(req.TopCategory)
	// var sqlQuery string
	// 通过category查询商品
	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		result := global.DB.Where("id=?", req.TopCategory).First(&category)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		subQuery := global.DB.Model(&model.Category{}).Select("id")
		if category.Level == 1 {
			// sqlQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id IN (SELECT id FROM category WHERE parent_category_id=%d)",
			// 	req.TopCategory)
			// 查询二级分类ID（一级分类的子分类）
			subQuery = subQuery.Where("parent_category_id IN (?)",
				global.DB.Model(&model.Category{}).
					Select("id").
					Where("parent_category_id = ?", req.TopCategory))
		} else if category.Level == 2 {
			// sqlQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id =%d",
			// 	req.TopCategory)
			// 查询三级分类ID（二级分类的子分类）
			subQuery = subQuery.Where("parent_category_id = ?", req.TopCategory)
			// localDB = localDB.Joins(
			// 	"JOIN category c3 ON goods.category_id = c3.id "+
			// 		"JOIN category c2 ON c3.parent_category_id = c2.id "+
			// 		"WHERE c2.id = ?", req.TopCategory)
		} else if category.Level == 3 {
			// sqlQuery = fmt.Sprintf("SELECT id FROM category WHERE id =%d", req.TopCategory)
			subQuery = subQuery.Where("id = ?", req.TopCategory)
		}
		// localDB = localDB.Where(fmt.Sprintf("category_id IN (%s)", sqlQuery))
		// localDB = localDB.Where("category_id IN (?)", subQuery)
		type Result struct {
			ID int32
		}
		var results []Result
		if err := subQuery.Scan(&results).Error; err != nil {
			return nil, status.Errorf(codes.Internal, "获取商品总数失败: %v", err)
		}
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		//生成terms查询
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	// 分页
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 10:
		req.PagePerNums = 10
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).
		Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	goodsIds := make([]int32, 0)
	goodsListRsp.Total = int32(result.Hits.TotalHits.Value)
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	// 查询id在某个数组中的值
	var goods []model.Goods
	res := localDB.Preload("Category").Preload("Brand").Find(&goods, goodsIds)
	if res.Error != nil {
		return nil, res.Error
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

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods
	result := global.DB.First(&goods, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoRsp := ModelToResponse(goods)
	return &goodsInfoRsp, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	result := global.DB.First(&category, req.CategoryId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brand
	result = global.DB.First(&brand, req.BrandId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	goods := model.Goods{
		Brand:           brand,
		BrandID:         brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}
	tx := global.DB.Begin()
	result = tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	result := global.DB.First(&model.Goods{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var category model.Category
	result := global.DB.First(&category, req.CategoryId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brand
	result = global.DB.First(&brand, req.BrandId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	var goods model.Goods
	goods.Brand = brand
	goods.BrandID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale
	global.DB.Save(&goods)
	return &emptypb.Empty{}, nil
}
