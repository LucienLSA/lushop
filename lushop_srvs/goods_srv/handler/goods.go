package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goodssrv/global"
	"goodssrv/model"
	proto "goodssrv/proto"

	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// 行代码的作用是编译期接口实现检查。
// &GoodsServer{} 必须实现 proto.GoodsServer 接口，否则编译报错。
// 保证你 GoodsServer 结构体确实实现了 proto 生成的所有 gRPC 方法，防止漏写。
var _ proto.GoodsServer = &GoodsServer{}

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
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//关键词搜索，查询新品，查询热门商品，通过价格区间筛选，通过商品分类筛选
	//使用es的目的是搜索出商品的id来，通过id拿到具体的字段信息是通过mysql来完成的
	//我们使用es是用来做搜索的，是否应该将所有的mysql字段全部在es中保存一份？
	//es用来做搜索，这个时候我们一般只会把搜索和过滤的字段信息保存到es中
	//es可以用来当作mysql使用，但是实际上mysql和es之间是互补的关系 ，  一般mysql用来做存储使用，es用来做搜索使用
	//es想要提高性能，就要将es的内存设置的够大（有最大限制）  占内存1k  2k  没必要的字段不要存
	goodsListResponse := &proto.GoodsListResponse{}
	//match bool 复合查询
	//q = q.Must(NewTermQuery("tag", "wow"))
	//q = q.Filter(NewTermQuery("account", "1"))
	//q := elastic.NewBoolQuery()
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		//搜索
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
		//q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
		//Filter不会算分  Must会参数得分
		//q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		localDB = localDB.Where(model.Goods{IsNew: true})
		//q = q.Filter(elastic.NewTermQuery("is_new", req.IsHot))
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price>=?", req.PriceMin)
		//q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price<=?", req.PriceMax)
		//q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id=?", req.Brand)
		//q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}
	//通过category去查询商品
	//用mysql查询取id放到categoryIds用es查询
	var subQuery string
	categoryIds := make([]interface{}, 0)
	var goods []model.Goods
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id IN (SELECT id FROM category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("SELECT id FROM category WHERE id=%d", req.TopCategory)
		}
		type Result struct {
			ID int32
		}
		var results []Result
		global.DB.Model(&model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		//生成terms查询
		//基础知识：函数参数是...interface{}的时候传递多个值没有问题   但就是不能传[]int{1,2,3}...   必须得是[]interface{}{1,2,3}...
		//q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
		//localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery)).Find(&goods)
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}
	//动词 名词 执行条件 确定执行
	//.Query(q).From().Size() - 分页
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	//result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	//if err != nil {
	//	return nil, err
	//}
	//goodsIds := make([]int32, 0)
	//goodsListResponse.Total = int32(result.Hits.TotalHits.Value)
	//for _, value := range result.Hits.Hits {
	//	goods := model.EsGoods{}
	//	_ = json.Unmarshal(value.Source, &goods)
	//	goodsIds = append(goodsIds, goods.ID)
	//}
	//if len(goodsIds) == 0 {
	//	return &proto.GoodsListResponse{}, nil
	//}
	//var count int64
	//localDB.Count(&count)
	//goodsListResponse.Total = int32(count)

	//if result2 := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))); result2.Error != nil {
	//	return nil, result2.Error
	//}
	//查询id在某个数组中的值

	re := localDB.Preload("Category").Preload("Brand").Find(&goods)
	if re.Error != nil {
		return nil, re.Error
	}
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	return goodsListResponse, nil
}

// 商品接口
// 关键词搜索、查询新品、查询热门商品、通过价格区间筛选、通过商品分类筛选
func (s *GoodsServer) GoodsListES(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
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
	// 总数
	goodsListRsp.Total = int32(result.Hits.TotalHits.Value)
	// 解析json数据
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
	// 查询的结果保持到返回响应的结构体
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
	if result.Error != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数无效")
	}
	for _, good := range goods {
		goodsInfoRsp := ModelToResponse(good)
		goodsListRsp.Data = append(goodsListRsp.Data, &goodsInfoRsp)
	}
	goodsListRsp.Total = int32(result.RowsAffected)
	return goodsListRsp, nil
}

// 商品详情
func (s *GoodsServer) GetGoodsDetail1(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods
	result := global.DB.First(&goods, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoRsp := ModelToResponse(goods)
	return &goodsInfoRsp, nil
}

// 商品详情
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Brand").Preload("Category").Where(req.Id).First(&goods); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "参数无效")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

// 新建商品
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品已存在")
	}
	goods = model.Goods{
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
	if result := tx.Save(&goods); result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "新增商品失败")
	}
	tx.Commit()
	goodsDetailResponse := ModelToResponse(goods)
	return &goodsDetailResponse, nil
}

// 删除商品
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新商品以及部分更新
func (s *GoodsServer) UpdateGoods1(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
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
	if result := global.DB.Save(&goods); result.Error != nil {
		return nil, status.Error(codes.Internal, "更新商品失败")
	}
	return &emptypb.Empty{}, nil
}

// 事务更新商品
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods
	// 如果请求中没有传递 CategoryId 和 BrandId，只更新商品的 IsNew、IsHot、OnSale 这三个布尔字段。
	if req.CategoryId == 0 && req.BrandId == 0 {
		if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
		}
		goods.IsNew = req.IsNew
		goods.IsHot = req.IsHot
		goods.OnSale = req.OnSale
		tx := global.DB.Begin()
		if result := tx.Save(&goods); result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
		tx.Commit()
		return &emptypb.Empty{}, nil
	}
	// 如果要更新商品的分类或品牌，先检查数据库中是否存在对应的分类和品牌。
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}
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
	// 所有数据库操作都在事务中进行，确保数据一致性。
	tx := global.DB.Begin()
	if result := tx.Save(&goods); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
