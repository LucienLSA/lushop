package handler

import (
	"context"

	"useropsrv/global"
	"useropsrv/model"
	proto "useropsrv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (*UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var rsp proto.UserFavListResponse
	var userFavs []model.UserFav
	var userFavList []*proto.UserFavResponse
	//查询用户的收藏记录
	//查询某件商品被哪些用户收藏了
	result := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Find(&userFavs)
	rsp.Total = int32(result.RowsAffected)

	for _, userFav := range userFavs {
		userFavList = append(userFavList, &proto.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}
	rsp.Data = userFavList
	return &rsp, nil
}

func (*UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	var userFav model.UserFav

	userFav.User = req.UserId
	userFav.Goods = req.GoodsId

	if result := global.DB.Save(&userFav); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (*UserOpServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	if result := global.DB.Unscoped().Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.UserFav{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavResponse, error) {
	var userfav model.UserFav
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&userfav); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &proto.UserFavResponse{
		UserId:  userfav.User,
		GoodsId: userfav.Goods,
	}, nil
}
