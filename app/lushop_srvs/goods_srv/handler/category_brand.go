package handler

import (
	"context"
	"goodssrv/global"
	"goodssrv/model"
	proto "goodssrv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 获取品牌与商品分类列表
func (s *GoodsServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var categoryBrands []model.GoodsCategoryBrand
	categoryBrandListRsp := proto.CategoryBrandListResponse{}

	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total)
	categoryBrandListRsp.Total = int32(total)

	global.DB.Preload("Category").Preload("Brand").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrands)

	var categoryBrandRsp []*proto.CategoryBrandResponse
	for _, categoryBrand := range categoryBrands {
		categoryBrandRsp = append(categoryBrandRsp, &proto.CategoryBrandResponse{
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrand.Category.ID,
				Name:           categoryBrand.Category.Name,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
			},
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.Brand.ID,
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.Logo,
			},
			Id: categoryBrand.ID,
		})
	}
	categoryBrandListRsp.Data = categoryBrandRsp
	return &categoryBrandListRsp, nil
}

// 通过商品分类获取品牌
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	// // 在需要的地方创建无事务的会话
	// db := global.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	brandListRsp := proto.BrandListResponse{}
	var category model.Category
	// 查询品牌分类是否存在
	// result := db.First(&category, req.Id)
	// if result.RowsAffected == 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	// }
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	// 根据商品分类查询品牌商品分类表
	var categoryBrands []model.GoodsCategoryBrand
	// result = db.Preload("Brands").Where(&model.GoodsCategoryBrand{CategoryID: category.ID}).Find(&categoryBrands)
	// if result.RowsAffected > 0 {
	// 	brandListRsp.Total = int32(result.RowsAffected)
	// }
	if result := global.DB.Preload("Brand").Where(&model.GoodsCategoryBrand{CategoryID: req.Id}).Find(&categoryBrands); result.RowsAffected > 0 {
		brandListRsp.Total = int32(result.RowsAffected)
	}
	// 查询到的品牌结果与响应返回绑定
	var brandInfoRsp []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoRsp = append(brandInfoRsp, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brand.ID,
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.Logo,
		})
	}
	brandListRsp.Data = brandInfoRsp
	return &brandListRsp, nil
}

func (s *GoodsServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
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
	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		BrandID:    req.BrandId,
	}
	if result := global.DB.Save(&categoryBrand); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "创建商品分类品牌失败")
	}
	return &proto.CategoryBrandResponse{
		Id: categoryBrand.ID,
	}, nil
}

func (s *GoodsServer) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.GoodsCategoryBrand{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	var categoryBrand model.GoodsCategoryBrand
	result := global.DB.First(&categoryBrand, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类品牌不存在")
	}

	var category model.Category
	result = global.DB.First(&category, req.CategoryId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	var brand model.Brand
	result = global.DB.First(&brand, req.BrandId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	categoryBrand.CategoryID = req.CategoryId
	categoryBrand.BrandID = req.BrandId
	if result := global.DB.Save(&categoryBrand); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新商品分类品牌失败")
	}
	return &emptypb.Empty{}, nil
}
