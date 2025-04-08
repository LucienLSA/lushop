package handler

import (
	"context"
	"goodssrv/global"
	"goodssrv/model"
	"goodssrv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播图
func (s *GoodsServer) BannerList(ctx context.Context, req *emptypb.Empty) (*proto.BannerListResponse, error) {
	bannerListRsp := proto.BannerListResponse{}
	var banners []model.Banner
	result := global.DB.Find(&banners)
	bannerListRsp.Total = int32(result.RowsAffected)
	var bannerRsp []*proto.BannerResponse
	for _, banner := range banners {
		bannerRsp = append(bannerRsp, &proto.BannerResponse{
			Id:    int32(banner.ID),
			Image: banner.Image,
			Index: int32(banner.Index),
			Url:   banner.Url,
		})
	}
	bannerListRsp.Data = bannerRsp
	return &bannerListRsp, nil
}
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{}
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url
	global.DB.Save(&banner)
	return &proto.BannerResponse{Id: banner.ID}, nil
}
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Banner{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	banners := model.Banner{}
	result := global.DB.First(&banners)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	if req.Url != "" {
		banners.Url = req.Url
	}
	if req.Image != "" {
		banners.Image = req.Image
	}
	if req.Index != 0 {
		banners.Index = req.Index
	}
	global.DB.Save(&banners)
	return &emptypb.Empty{}, nil
}
