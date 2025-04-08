package handler

import (
	"context"
	"useropsrv/global"
	"useropsrv/model"
	proto_message "useropsrv/proto/gen/message"
)

// 获取留言列表
func (*UserOpServer) MessageList(ctx context.Context, req *proto_message.MessageRequest) (*proto_message.MessageListResponse, error) {
	var rsp proto_message.MessageListResponse
	var messages []model.LeavingMessages
	var messageList []*proto_message.MessageResponse

	result := global.DB.Where(&model.LeavingMessages{User: req.UserId}).Find(&messages)
	rsp.Total = int32(result.RowsAffected)

	for _, message := range messages {
		messageList = append(messageList, &proto_message.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	rsp.Data = messageList
	return &rsp, nil
}

func (*UserOpServer) CreateMessage(ctx context.Context, req *proto_message.MessageRequest) (*proto_message.MessageResponse, error) {
	var message model.LeavingMessages

	message.User = req.UserId
	message.MessageType = req.MessageType
	message.Subject = req.Subject
	message.Message = req.Message
	message.File = req.File

	global.DB.Save(&message)

	return &proto_message.MessageResponse{Id: message.ID}, nil
}
