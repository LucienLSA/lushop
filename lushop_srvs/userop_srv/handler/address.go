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

// 获取地址列表
func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addressInfo []model.Address
	var addressRsp proto.AddressListResponse
	var addressResponses []*proto.AddressResponse

	if result := global.DB.Where(&model.Address{User: req.UserId}).Find(&addressInfo); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "不存在收货地址")
	} else {
		addressRsp.Total = int32(result.RowsAffected)
	}

	for _, address := range addressInfo {
		addressResponses = append(addressResponses, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	addressRsp.Data = addressResponses
	return &addressRsp, nil
}

func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	var address model.Address

	address.User = req.UserId
	address.Province = req.Province
	address.City = req.City
	address.District = req.District
	address.Address = req.Address
	address.SignerName = req.SignerName
	address.SignerMobile = req.SignerMobile

	if result := global.DB.Create(&address); result.Error != nil {
		return nil, result.Error
	}

	return &proto.AddressResponse{Id: address.ID}, nil
}

func (*UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where(&model.Address{BaseModel: model.BaseModel{ID: req.Id}}).First(&model.Address{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "收货地址不存在")
	}
	if result := global.DB.Where(&model.Address{BaseModel: model.BaseModel{ID: req.Id}}).Delete(&model.Address{}); result.Error != nil {
		return nil, result.Error
	}
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	var address model.Address

	if result := global.DB.Where("id=? and user=?", req.Id, req.UserId).First(&address); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if address.Province != "" {
		address.Province = req.Province
	}

	if address.City != "" {
		address.City = req.City
	}

	if address.District != "" {
		address.District = req.District
	}

	if address.Address != "" {
		address.Address = req.Address
	}

	if address.SignerName != "" {
		address.SignerName = req.SignerName
	}

	if address.SignerMobile != "" {
		address.SignerMobile = req.SignerMobile
	}

	if result := global.DB.Save(&address); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
