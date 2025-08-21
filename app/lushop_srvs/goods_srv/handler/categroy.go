package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goodssrv/global"
	"goodssrv/model"
	proto "goodssrv/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 商品分类
func (s *GoodsServer) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*
		[
			{
				"id": xx,
				"name": "",
				"level": 1,
				"is_tap": false,
				"parent": xx,
				"sub_category":[
					"id": xx,
					"name": "",
					"level": 1,
					"is_tap": false,
					"sub_category":[]
				]
			}
		]
	*/
	// 将 Total Data  JsonData 分别获取存入到结构体中，得到嵌套的输出结构体
	categoryRsp := proto.CategoryListResponse{}
	var categorys []model.Category
	if result := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys); result.Error != nil {
		return nil, result.Error
	}
	b, _ := json.Marshal(&categorys)
	categoryRsp.JsonData = string(b)

	var categorys_proto []model.Category
	result := global.DB.Find(&categorys_proto)
	if result.Error != nil {
		return nil, result.Error
	}

	categoryRsp.Total = int32(result.RowsAffected)
	for _, category := range categorys_proto {
		categoryInfo := proto.CategoryInfoResponse{}
		categoryInfo.Id = category.ID
		categoryInfo.Name = category.Name
		categoryInfo.ParentCategory = category.ParentCategoryID
		categoryInfo.Level = category.Level
		categoryInfo.IsTab = category.IsTab
		categoryRsp.Data = append(categoryRsp.Data, &categoryInfo)
	}
	return &categoryRsp, nil

}

// 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 在需要的地方创建无事务的会话
	// db := global.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	// 父分类和子分类的信息都需要获取
	categoryListRsp := proto.SubCategoryListResponse{}
	var category model.Category
	// 1. 传入请求查询的Id，作为父分类的id
	result := global.DB.First(&category, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	// 父分类返回响应
	categoryListRsp.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	// // 2. 单独查询子分类（不带预加载）
	// if err := db.Where("parent_category_id = ?", req.Id).Find(&subCategorys).Error; err != nil {
	// 	return nil, err
	// }
	// // 3. 如果一级分类，再尝试预加载多级
	// if category.Level == 1 {
	// 	if err := global.DB.Preload("SubCategory.SubCategory").Find(&subCategorys).Error; err != nil {
	// 		return nil, err
	// 	}
	// }
	var subCategorys []model.Category
	var subCategoryRsp []*proto.CategoryInfoResponse
	preloads := "SubCategory"
	// 如果查询的等级为1，获取所有分类级别，这里由于获取的是子分类，不需要三级分类
	if category.Level == 1 {
		preloads = "SubCategory.SubCategory"
	}
	fmt.Println(preloads)
	// 但本实现只查直接子分类。
	// 如果查询的等级不为1，根据父类的id查询子分类
	// global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategorys)
	if result := global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys); result.Error != nil {
		return nil, result.Error
	}
	// 子分类信息存入响应
	for _, subCategory := range subCategorys {
		subCategoryRsp = append(subCategoryRsp, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}
	categoryListRsp.SubCategorys = subCategoryRsp
	return &categoryListRsp, nil
}

func (s *GoodsServer) CreateCategory1(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	if req.Level != 1 {
		// 这里减少对父类目的查询，是否存在
		category.ParentCategoryID = req.ParentCategory
	}
	category.IsTab = req.IsTab
	if result := global.DB.Save(&category); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "创建商品分类失败")
	}
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

// 新建商品分类
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	// 验证父分类是否存在
	if req.Level != 1 {
		var parentCategory model.Category
		result := global.DB.First(&parentCategory, req.ParentCategory)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "父分类不存在")
		}
	}

	category := model.Category{
		Name:  req.Name,
		Level: req.Level,
		IsTab: req.IsTab,
	}
	if req.Level != 1 {
		category.ParentCategoryID = req.ParentCategory
	}

	if result := global.DB.Create(&category); result.Error != nil {
		zap.S().Error("新建商品分类失败！")
		return nil, status.Errorf(codes.Internal, "创建商品分类失败")
	}
	zap.S().Infof("category ID:%d", category.ID)
	return &proto.CategoryInfoResponse{
		Id:             int32(category.ID),
		Name:           category.Name,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
	}, nil
}

// 删除商品分类
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Category{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新商品分类
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	category := model.Category{}
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	if result := global.DB.Save(&category); result.Error != nil {
		zap.S().Error("更新商品分类失败", result.Error)
		return nil, status.Errorf(codes.Internal, "更新商品分类失败")
	}
	return &emptypb.Empty{}, nil
}
