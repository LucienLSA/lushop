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

// 获取品牌列表
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	// brandListResponse := proto.BrandListResponse{}
	var brandListResponse proto.BrandListResponse
	var brands []model.Brand
	// 查询数据库
	// result := global.DB.Find(&brands)
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	global.DB.Model(&model.Brand{}).Count(&total)
	var brandRsps []*proto.BrandInfoResponse
	// 数据库查询结果绑定响应
	for _, brand := range brands {
		brandRsp := &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		}
		brandRsps = append(brandRsps, brandRsp)
	}
	brandListResponse.Data = brandRsps
	brandListResponse.Total = int32(total)
	return &brandListResponse, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 处理重名情况
	// result := global.DB.First(&model.Brand{})
	// if result.RowsAffected == 1 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	// }
	brand := &model.Brand{Name: req.Name}
	if result := global.DB.Where(&brand).First(brand); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	// brand := model.Brands{}
	// brand.Name = req.Name
	// brand.Logo = req.Logo
	// global.DB.Create(model.Brands{
	// 	Name: req.Name,
	// 	Logo: req.Logo,
	// })
	// brand := &model.Brand{
	// 	Name: req.Name,
	// 	Logo: req.Logo,
	// }
	brand.Logo = req.Logo
	if result := global.DB.Save(&brand); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建品牌失败")
	}
	return &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}, nil
}

func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Brand{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brand{}
	if result := global.DB.Where(&model.Brand{BaseModel: model.BaseModel{ID: req.Id}}).First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}
	if result := global.DB.Save(&brands); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建品牌失败")
	}
	return &emptypb.Empty{}, nil
}
