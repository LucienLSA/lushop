package handler

import (
	"context"
	"encoding/json"
	"lushopsrvs/goods_srv/global"
	"lushopsrvs/goods_srv/model"
	"lushopsrvs/goods_srv/proto"

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
// func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error)
// func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error)
// func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error)
// func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error)
