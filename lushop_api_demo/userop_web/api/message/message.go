package message

import (
	"context"
	"useropweb/api"
	"useropweb/forms"
	"useropweb/global"
	proto_message "useropweb/proto/gen/message"
	"useropweb/utils/jwtClaims"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"net/http"
)

func List(ctx *gin.Context) {
	request := &proto_message.MessageRequest{}

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*jwtClaims.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.MessageSrvClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["type"] = value.MessageType
		reMap["subject"] = value.Subject
		reMap["message"] = value.Message
		reMap["file"] = value.File

		result = append(result, reMap)
	}
	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.MessageSrvClient.CreateMessage(context.Background(), &proto_message.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
