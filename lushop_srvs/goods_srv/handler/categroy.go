package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"lushopsrvs/goods_srv/global"
	"lushopsrvs/goods_srv/model"
	"lushopsrvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
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
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
	// for _, category := range categorys {
	// 	fmt.Println(category.Name)
	// }
	// return nil, nil
}

// 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 在需要的地方创建无事务的会话
	db := global.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
	// 父分类和子分类的信息都需要获取
	categoryListRsp := proto.SubCategoryListResponse{}
	var category model.Category
	// 1. 传入请求查询的Id，作为父分类的id
	result := db.First(&category, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	// 父分类
	categoryListRsp.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	// 检查数据库连接是否正常
	if err := db.Exec("SELECT 1").Error; err != nil {
		fmt.Println("DB connection error:", err)
	}
	// 检查数据库的事务是否开启
	fmt.Println("Is DB in transaction?", db.Statement.ConnPool != db.ConnPool)
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
	// 如果查询的等级不为1，根据父类的id查询子分类
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategorys)
	// global.DB.Where("parent_category_id=? AND deleted_at IS NULL", req.Id).Find(&subCategorys)
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

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	if req.Level != 1 {
		// 这里减少对父类目的查询，是否存在
		category.ParentCategoryID = req.ParentCategory
	}
	category.IsTab = req.IsTab
	global.DB.Save(&category)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Category{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	category := model.Category{}
	result := global.DB.First(&category, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
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
	global.DB.Save(&category)
	return &emptypb.Empty{}, nil
}
