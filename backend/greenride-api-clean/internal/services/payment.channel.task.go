package services

import (
	"context"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/task"
	"log"
)

// SyncPaymentChannels 更新渠道服务映射
func SyncPaymentChannels(channels []*models.PaymentChannels) {
	// 创建新的渠道映射
	newChannels := make(map[string]PaymentChannel)

	// 加载渠道服务
	for _, item := range channels {
		if item.Status != protocol.StatusActive {
			continue
		}
		var svc PaymentChannel
		switch item.ChannelCode {
		case protocol.PaymentChannelStripe:
			svc = NewStripeService(item)
		case protocol.PaymentChannelMoMo:
			svc = NewMoMoService(item)
		case protocol.PaymentChannelKPay:
			svc = NewKPayService(item)
		default:
			log.Printf("Unsupported payment channel: %s", item.ChannelCode)
			continue
		}

		if svc != nil {
			newChannels[item.AccountID] = svc
		}
	}

	// 更新全局渠道映射
	PaymentChannels = newChannels
}

func InitPaymentChannelHandlers() {

	// 注册渠道同步任务处理器
	task.RegisterHandler(protocol.PaymentChannelSyncHandler, syncPaymentChannelsHandler)

	// 创建定时任务 - 每5分钟执行一次支付渠道同步
	syncTask := &models.Task{
		TaskID:     protocol.PaymentChannelSyncHandler,
		HandlerKey: protocol.PaymentChannelSyncHandler,
		Name:       "支付渠道同步任务",
		Cron:       "*/5 * * * *", // 每5分钟执行一次
		Timeout:    300,           // 5分钟超时
		Status:     protocol.TaskStatusEnabled,
		Params:     protocol.MapData{},
	}

	// 初始化任务
	if err := task.InitTask(syncTask); err != nil {
		log.Printf("Failed to initialize payment channel sync task: %v", err)
	} else {
		log.Printf("Payment channel sync task initialized successfully")
	}
	go func() {
		//马上执行同步
		syncPaymentChannelsHandler(context.Background(), protocol.MapData{})
	}()
}

// syncPaymentChannelsHandler 渠道同步处理器
func syncPaymentChannelsHandler(ctx context.Context, params protocol.MapData) error {
	// 查询活跃的支付渠道
	channels, err := models.GetActiveChannels()
	if err != nil {
		log.Printf("Failed to query payment channels: %v", err)
		return err
	}

	// 更新渠道服务
	SyncPaymentChannels(channels)
	return nil
}
