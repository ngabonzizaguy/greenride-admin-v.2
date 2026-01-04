package services

import (
	"context"
	"errors"
	"fmt"
	"greenride/internal/models"
	"greenride/internal/protocol"

	"github.com/spf13/cast"
)

// 发送FCM消息的业务逻辑
func (m *MessageService) SendFcmMessage(message *Message) error {
	if m.FirebaseService == nil {
		return errors.New("firebase service not available")
	}

	// 使用统一的PrepareMessage方法准备FCM消息
	fcmMessage := m.PrepareMessage(message, protocol.MsgChannelFcm)
	if fcmMessage == nil {
		return fmt.Errorf("failed to prepare FCM message for type: %s, language: %s", message.Type, message.Language)
	}

	// 构建参数
	params := protocol.MapData{}
	params.Copy(DefaultParams)
	params.Copy(fcmMessage.Params)

	// 参数转换为字符串，过滤掉不适合FCM的字段
	stringParams := make(map[string]string)
	for k, v := range params {
		valueStr := cast.ToString(v)
		// 过滤掉空值的保留字段和不适合FCM数据的字段
		if (k == "from" || k == "to") && valueStr == "" {
			continue
		}
		// 只保留真正需要的数据字段，排除控制字段
		if k != "title" && k != "content" && k != "token" {
			stringParams[k] = valueStr
		}
	}

	// 获取标题和内容
	titleVal := params["title"]
	contentVal := params["content"]
	titleStr := cast.ToString(titleVal)
	contentStr := cast.ToString(contentVal)

	// 内容校验
	if contentStr == "" {
		return errors.New("empty content for FCM message")
	}

	// 构建FCM消息
	fcmMsg := &protocol.FCMMessage{
		Title: titleStr,
		Body:  contentStr,
		Data:  stringParams,
	}

	// 根据参数决定发送方式
	if tokenVal, ok := params["token"]; ok && tokenVal != nil {
		// 单个token
		token := cast.ToString(tokenVal)
		if token == "" {
			return errors.New("empty FCM token provided")
		}

		fcmMsg.Token = token

		// 通过token查找对应的用户ID
		if fcmToken := models.GetFcmTokenByToken(token); fcmToken != nil {
			fcmMsg.UserID = fcmToken.UserID
		}

		ctx := context.Background()
		_, err := m.FirebaseService.SendFCMMessage(ctx, fcmMsg)
		if err != nil {
			return fmt.Errorf("failed to send FCM message: %v", err)
		}
	} else if userIDVal, ok := params["to"]; ok && userIDVal != nil {
		// 使用用户ID获取所有token
		userID := cast.ToString(userIDVal)
		if userID == "" {
			return errors.New("empty user ID provided")
		}

		// 设置用户ID到FCM消息中
		fcmMsg.UserID = userID

		tokens := m.FirebaseService.GetUserFCMTokens(userID)
		if len(tokens) == 0 {
			return fmt.Errorf("no FCM tokens found for user: %s", userID)
		}

		// 多个token
		ctx := context.Background()
		_, err := m.FirebaseService.SendMulticastFCMMessage(ctx, fcmMsg, tokens)
		if err != nil {
			return fmt.Errorf("failed to send multicast FCM message: %v", err)
		}
	} else {
		return errors.New("neither token nor userID provided for FCM message")
	}

	return nil
}
