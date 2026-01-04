package services

import (
	"context"

	"greenride/internal/models"
	"greenride/internal/protocol"
)

// SendSignal 发送信号到Redis（通用方法）
func SendSignal(signalType string, bizID string) error {
	return models.GetRedis().SAdd(
		context.Background(),
		signalType,
		bizID,
	).Err()
}

// SendUserRegisteredSignal 发送用户注册信号
func SendUserRegisteredSignal(userID string) error {
	return SendSignal(protocol.SignalUserRegistered, userID)
}
