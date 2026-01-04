package task

import (
	"context"
	"fmt"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// 信号类型与Handler的映射关系
var signalHandlerMap = map[string]string{
	protocol.SignalOrderSubmit:    protocol.OrderSubmitHandler,
	protocol.SignalOrderComplete:  protocol.OrderCompleteHandler,
	protocol.SignalOrderCancel:    protocol.OrderCancelHandler,
	protocol.SignalUserRegistered: protocol.UserRegisteredHandler,
}

// ProcessSignals 统一信号处理Task Handler
func ProcessSignals(ctx context.Context, params map[string]any) error {
	// 遍历所有信号类型
	for signalType, handlerKey := range signalHandlerMap {
		if err := processSignalType(ctx, signalType, handlerKey); err != nil {
			log.Get().Errorf("Failed to process signal type %s: %v", signalType, err)
			// 继续处理其他信号类型，不因一个失败而停止
		}
	}
	return nil
}

// processSignalType 处理单个信号类型
func processSignalType(ctx context.Context, signalType, handlerKey string) error {
	// 原子操作：读取并清空Redis中的业务ID数组
	redisCtx := context.Background()
	pipe := models.GetRedis().TxPipeline()
	membersCmd := pipe.SMembers(redisCtx, signalType)
	pipe.Del(redisCtx, signalType)
	_, err := pipe.Exec(redisCtx)
	if err != nil {
		return err
	}

	bizIDs := membersCmd.Val()
	if len(bizIDs) == 0 {
		return nil // 没有待处理信号
	}

	log.Get().Infof("Processing %d signals for type: %s", len(bizIDs), signalType)

	// 获取对应的Task Handler
	handler, exists := GetHandler(handlerKey)
	if !exists {
		return fmt.Errorf("handler not found for key: %s", handlerKey)
	}

	// 为每个业务ID处理信号
	for _, bizID := range bizIDs {
		lockKey := fmt.Sprintf("%s:%s", signalType, bizID)

		// 使用Task的分布式锁机制
		mutex, err := Lock(lockKey)
		if err != nil {
			log.Get().Infof("Signal %s already processing, skipped", lockKey)
			continue
		}

		// 构造符合Task Handler要求的参数
		handlerParams := map[string]any{
			"signal_type": signalType,
			"biz_id":      bizID,
		}

		// 调用Task Handler
		if err := handler(ctx, handlerParams); err != nil {
			log.Get().Errorf("Failed to handle signal %s:%s, error: %v", signalType, bizID, err)
		} else {
			log.Get().Infof("Successfully handled signal %s:%s", signalType, bizID)
		}

		// 释放锁
		Unlock(mutex)
	}

	return nil
}
