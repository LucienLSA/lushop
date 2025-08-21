package handler

import (
	"context"
	"fmt"
	"usersrv/global"
	"usersrv/model"
	proto "usersrv/proto"

	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user *model.User) *proto.UserInfoResponse {
	// grpc中的message中字段有默认值，不能随意赋值nil，会错
	// 哪些字段是有默认值的
	userInfoRsp := &proto.UserInfoResponse{
		Id:       int32(user.ID),
		PassWord: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.UTC().Unix())
	}
	return userInfoRsp
}

// func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if page == 0 {
// 			page = 1
// 		}
// 		switch {
// 		case pageSize > 100:
// 			pageSize = 100
// 		case pageSize <= 100:
// 			pageSize = 10
// 		}
// 		offset := (page - 1) * pageSize
// 		return db.Offset(offset).Limit(pageSize)
// 	}
// }

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// db = global.DB
		if pageNum < 1 {
			pageNum = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 1:
			pageSize = 10
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表
	fmt.Println("获取用户列表")
	var users []model.User
	// db := global.NewDBClient(ctx)
	// result := global.DB.Find(&users)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }
	// rsp := &proto.UserListResponse{}
	// // rsp.Total = int32(result.RowsAffected)
	// // global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	// var total int64
	// err := global.DB.Model(&model.User{}).Count(&total).Error
	// if err != nil {
	// 	return nil, status.Errorf(codes.NotFound, "用户不存在")
	// }
	// rsp.Total = int32(total)
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := ModelToResponse(&user)
		rsp.Data = append(rsp.Data, userInfoRsp)
	}
	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 手机号码查询用户
	user := model.User{} // 每次都新建
	// db := global.NewDBClient(ctx)
	fmt.Println(req.Mobile)
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(&user)
	return userInfoRsp, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// id查询用户
	user := model.User{} // 每次都新建
	// db := global.DB
	// db := global.NewDBClient(ctx)
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(&user)
	return userInfoRsp, nil
}

// 更新用户密码
func (s *UserServer) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordInfo) (*empty.Empty, error) {
	// 1. 查询用户是否存在
	var user model.User
	if result := global.DB.First(&user, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// 2. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "密码加密失败")
	}

	// 3. 更新密码
	user.Password = string(hashedPassword)
	if result := global.DB.Save(&user); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新密码失败")
	}
	return &empty.Empty{}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 新建用户
	user := model.User{} // 每次都新建
	var count int64
	err := global.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).Count(&count).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "数据库错误: %v", err)
	}
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	// result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	// if result.RowsAffected == 1 {
	// 	return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	// }
	// if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	return nil, status.Errorf(codes.Internal, "数据库错误: %v", result.Error)
	// }
	// 构造新用户
	user.Mobile = req.Mobile
	user.NickName = req.NickName
	if err := user.SetPassword(req.PassWord); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// 保存到数据库
	if err := global.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "创建用户失败: %v", err)
	}
	userInfoRsp := ModelToResponse(&user)
	return userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	user := model.User{} // 每次都新建
	// db := global.NewDBClient(ctx)
	// db := global.DB
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected < 1 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

// 检验密码，传入的是请求中的原密码
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	user := model.User{} // 每次都新建
	// 将请求中输入用户的密码代入user模型中对比
	user.Password = req.EncryptedPassWord
	// fmt.Println(user.Password)
	// fmt.Println(req.PassWord)
	ok := user.CheckPassword(req.PassWord)
	// if !ok {
	// 	return nil, status.Errorf(codes.Internal, "用户密码错误")
	// }
	return &proto.CheckResponse{Success: ok}, nil
}
