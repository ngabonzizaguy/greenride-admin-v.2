package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

var (
	firebaseServiceInstance *FirebaseService
	firebaseServiceOnce     sync.Once
)

type FirebaseService struct {
	app             *firebase.App
	messagingClient *messaging.Client
	db              *gorm.DB
}

// GetFirebaseService 获取Firebase服务单例
func GetFirebaseService() *FirebaseService {
	firebaseServiceOnce.Do(func() {
		SetupFirebaseService()
	})
	return firebaseServiceInstance
}

// SetupFirebaseService 设置Firebase服务
func SetupFirebaseService() {
	cfg := config.Get()
	if cfg.Firebase == nil {
		log.Get().Warn("Firebase configuration is nil, skipping Firebase service setup")
		return
	}

	// 使用json.Marshal将map转换为JSON字节数组的做法
	credentialsJSON, err := json.Marshal(cfg.Firebase)
	if err != nil {
		log.Get().Errorf("failed to marshal firebase config: %v", err)
		return
	}

	// 初始化Firebase应用
	opts := option.WithCredentialsJSON(credentialsJSON)
	app, err := firebase.NewApp(context.Background(), nil, opts)
	if err != nil {
		log.Get().Errorf("failed to initialize firebase app: %v", err)
		return
	}

	// 初始化Messaging客户端
	messagingClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Get().Errorf("failed to initialize firebase messaging client: %v", err)
		return
	}

	firebaseServiceInstance = &FirebaseService{
		app:             app,
		messagingClient: messagingClient,
		db:              models.GetDB(),
	}

	log.Get().Info("Firebase service initialized successfully")
}

// ValidateFCMToken 校验FCM token有效性
func (s *FirebaseService) ValidateFCMToken(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("FCM token is empty")
	}

	// 使用DryRun模式发送消息来验证token有效性
	// 这种方式不会实际发送消息，只进行验证
	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"type": "validation",
		},
	}

	// 使用SendDryRun进行验证，不会实际发送消息
	_, err := s.messagingClient.SendDryRun(ctx, message)
	return err
}

// ValidateMultipleFCMTokens 批量验证FCM token有效性
func (s *FirebaseService) ValidateMultipleFCMTokens(ctx context.Context, tokens []string) map[string]error {
	if len(tokens) == 0 {
		return nil
	}

	results := make(map[string]error)

	// 使用DryRun模式批量验证tokens
	// 为每个token创建消息并使用SendDryRun验证
	for _, token := range tokens {
		err := s.ValidateFCMToken(ctx, token)
		results[token] = err
	}

	return results
}

// GetValidTokens 从token列表中获取有效的token（串行验证）
func (s *FirebaseService) GetValidTokens(ctx context.Context, tokens []string) []string {
	if len(tokens) == 0 {
		return nil
	}

	validTokens := make([]string, 0)
	results := s.ValidateMultipleFCMTokens(ctx, tokens)

	for token, err := range results {
		if err == nil {
			validTokens = append(validTokens, token)
		}
	}

	return validTokens
}

// RegisterFCMToken 注册FCM Token
func (s *FirebaseService) RegisterFCMToken(userID, token, deviceID, platform, appID string) {
	if token == "" || userID == "" {
		return
	}

	// 先验证token的有效性
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.ValidateFCMToken(ctx, token); err != nil {
		log.Get().Warnf("invalid FCM token for user %s: %v", userID, err)
		return
	}

	fcmToken := models.GetFcmTokenByUserIDAndToken(userID, token)
	if fcmToken != nil {
		if fcmToken.Status != protocol.StatusActive {
			values := &models.FCMTokenValues{}
			values.SetDeviceID(deviceID).
				SetPlatform(platform).
				SetAppID(appID).
				SetActiveStatus()
			err := models.GetDB().Model(fcmToken).UpdateColumns(values).Error
			if err != nil {
				log.Get().Errorf("failed to update fcm token: %v", err)
			}
		}
		return
	}
	fcmToken = &models.FCMToken{
		UserID:   userID,
		Token:    token,
		DeviceID: deviceID,
		Platform: platform,
		AppID:    appID,
		Status:   protocol.StatusActive,
	}

	err := models.CreateFCMToken(fcmToken)
	if err != nil {
		log.Get().Errorf("failed to create fcm token: %v", err)
	}
}

// CreateFCMMessageRecord 创建FCM消息记录
func (s *FirebaseService) CreateFCMMessageRecord(msg *protocol.FCMMessage, firebaseMessageID string, status string, err error) *models.FCMMessageLog {
	nowtime := utils.TimeNowMilli()

	// message_id始终是我们自己生成的唯一标识符
	messageID := utils.GenerateFCMMessageID()

	// 将Data转换为JSON字符串
	var dataStr *string
	if msg.Data != nil {
		if dataJSON, err := json.Marshal(msg.Data); err == nil {
			dataJSONStr := string(dataJSON)
			dataStr = &dataJSONStr
		}
	}

	record := &models.FCMMessageLog{
		MessageID: messageID, // 我们自己生成的ID
		Salt:      fmt.Sprintf("%d", nowtime),
		FCMMessageLogValues: &models.FCMMessageLogValues{
			UserID:       &msg.UserID,
			TokenID:      &msg.Token,
			FCMToken:     &msg.Token,
			FCMMessageID: &firebaseMessageID, // Firebase返回的messageID（可能为空）
			Title:        &msg.Title,
			Body:         &msg.Body,
			ImageURL:     &msg.ImageURL,
			ClickAction:  &msg.ClickAction,
			Data:         dataStr,
			Status:       &status,
			SentAt:       &nowtime,
		},
	}

	if err != nil {
		errMsg := err.Error()
		record.ResponseMessage = &errMsg
	} // 设置消息类型
	if msg.Data != nil {
		if msgType, ok := msg.Data["msg_type"]; ok {
			record.MessageType = &msgType
		}
	}

	if err := models.GetDB().Create(record).Error; err != nil {
		log.Get().Errorf("Failed to create FCM message record: %v", err)
		return nil
	}
	return record
}

// SendFCMMessage 发送FCM消息到单个token
func (s *FirebaseService) SendFCMMessage(ctx context.Context, msg *protocol.FCMMessage) (string, error) {
	if msg == nil {
		return "", errors.New("message is nil")
	}

	if msg.Token == "" {
		return "", errors.New("token must be specified")
	}

	// 构建Firebase消息
	message := &messaging.Message{
		Token: msg.Token,
		Notification: &messaging.Notification{
			Title: msg.Title,
			Body:  msg.Body,
		},
		Data: msg.Data,
	}

	// 设置图片
	if msg.ImageURL != "" {
		message.Notification.ImageURL = msg.ImageURL
	}

	// 设置Android配置
	if msg.ClickAction != "" {
		message.Android = &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ClickAction: msg.ClickAction,
			},
		}
	}

	// 发送消息
	messageID, err := s.messagingClient.Send(ctx, message)
	result := protocol.StatusSuccess
	// 创建消息记录
	if err != nil {
		result = protocol.StatusFailed
	}
	s.CreateFCMMessageRecord(msg, messageID, result, err)

	return messageID, err
}

// SendMulticastFCMMessage 发送多播FCM消息
func (s *FirebaseService) SendMulticastFCMMessage(ctx context.Context, msg *protocol.FCMMessage, tokens []string) (int, error) {
	if msg == nil || len(tokens) == 0 {
		return 0, errors.New("message or tokens are empty")
	}

	// 生成批次ID
	batchID := fmt.Sprintf("FCM_BATCH_%s", utils.GenerateID())

	// 限制并发数量
	maxConcurrency := min(len(tokens), 10)

	// 创建信号量控制并发数
	semaphore := make(chan struct{}, maxConcurrency)
	results := make(chan struct {
		token     string
		messageID string
		err       error
	}, len(tokens))

	// 并发发送到所有token，复用SendFCMMessage
	for _, token := range tokens {
		go func(t string) {
			semaphore <- struct{}{} // 获取信号量
			defer func() {
				<-semaphore // 释放信号量
			}()

			// 创建token特定的消息
			msgCopy := *msg
			msgCopy.Token = t

			messageID, err := s.SendFCMMessage(ctx, &msgCopy)
			results <- struct {
				token     string
				messageID string
				err       error
			}{
				token:     t,
				messageID: messageID,
				err:       err,
			}
		}(token)
	}

	// 收集结果
	successCount := 0
	failureCount := 0
	for range tokens {
		result := <-results

		if result.err != nil {
			failureCount++
		} else {
			successCount++
		}
	}

	// 记录批次统计
	log.Get().Infof("FCM Send Batch[%s]: total %d, success: %d, failure: %d",
		batchID, len(tokens), successCount, failureCount)

	close(results)
	return successCount, nil
}

// DeactivateToken 停用特定设备的FCM token（
func (s *FirebaseService) DeactivateToken(userID, token string) error {
	if userID == "" {
		return errors.New("userID is required")
	}
	values := &models.FCMTokenValues{}
	values.SetStatus(protocol.StatusInactive)
	// 只停用特定设备的token
	err := s.db.Model(&models.FCMToken{}).
		Where("user_id = ? ", userID).
		Where("token = ?", token).
		UpdateColumns(values).Error

	if err != nil {
		log.Get().Errorf("failed to deactivate device token: %v", err)
		return err
	}

	log.Get().Infof("deactivated token for user %s, token: %s", userID, token)
	return nil
}

// CleanupInactiveTokens 清理非活跃的FCM tokens
func (s *FirebaseService) CleanupInactiveTokens() error {
	return models.CleanupInactiveFCMTokens()
}

// GetUserFCMTokens 获取用户的有效FCM tokens
func (s *FirebaseService) GetUserFCMTokens(userID string) []string {
	if userID == "" {
		return nil
	}

	// 从数据库获取用户的所有活跃FCM tokens
	var tokens []string
	err := models.GetDB().Model(&models.FCMToken{}).
		Where("user_id = ? AND status = ?", userID, protocol.StatusActive).
		Pluck("token", &tokens).Error

	if err != nil {
		log.Get().Errorf("Failed to get FCM tokens for user %s: %v", userID, err)
		return nil
	}

	return tokens
}
