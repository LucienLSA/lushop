package userop

import (
	"context"
	"net/http"

	v2base "lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	v2useropproto "lushopapi/proto/userop"
	"lushopapi/utils/jwtClaims"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func MsgList(ctx *gin.Context) {
	request := &v2useropproto.MessageRequest{}

	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*jwtClaims.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.UserOpSrvClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		v2base.HandleGrpcErrorToHttp(err, ctx)
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

func MsgCreate(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.UserOpSrvClient.CreateMessage(context.Background(), &v2useropproto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorf("添加留言失败,err:%s", err.Error())
		v2base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
