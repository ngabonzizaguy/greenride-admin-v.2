package services

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/task"
)

const (
	// 订单超时取消任务常量
	TaskOrderTimeoutCancel       = "order_timeout_cancel"
	TaskStuckOrderCancel         = "stuck_order_cancel"
	OrderTimeoutDuration         = 1 * time.Hour   // 1小时超时时间
	StuckOrderDuration           = 2 * time.Hour  // 2小时未更新视为卡单
)

// InitOrderTaskHandlers 初始化订单相关任务处理器
func InitOrderTaskHandlers() {

	// 注册订单超时取消任务处理器
	task.RegisterHandler(TaskOrderTimeoutCancel, OrderTimeoutCancelHandler)
	task.RegisterHandler(TaskStuckOrderCancel, StuckOrderCancelHandler)

	// 注册定时任务
	RegisterOrderTasks()
}

// RegisterOrderTasks 注册订单相关定时任务
func RegisterOrderTasks() {
	// 订单超时取消任务 - 每5分钟执行一次
	orderTimeoutCancelTask := &models.Task{
		TaskID:     "order_timeout_cancel_scheduler",
		Name:       "订单超时自动取消",
		Type:       "order",
		HandlerKey: TaskOrderTimeoutCancel,
		Cron:       "*/5 * * * *", // 每5分钟执行一次
		Status:     protocol.TaskStatusEnabled,
		MaxRetries: 3,
		Timeout:    300, // 5分钟超时
		Params: protocol.MapData{
			"timeout_minutes": 60, // 60分钟超时
		},
		Remark: "自动取消超过预定/预约时间1小时且未被司机接单的订单",
	}

	// 卡单自动取消：driver_arrived 或 in_progress 超过 2 小时未更新则取消，释放司机
	stuckOrderCancelTask := &models.Task{
		TaskID:     "stuck_order_cancel_scheduler",
		Name:       "卡单自动取消",
		Type:       "order",
		HandlerKey: TaskStuckOrderCancel,
		Cron:       "*/15 * * * *", // 每15分钟执行一次
		Status:     protocol.TaskStatusEnabled,
		MaxRetries: 3,
		Timeout:    300,
		Params: protocol.MapData{
			"stuck_hours": 2.0, // 超过 2 小时未更新视为卡单
		},
		Remark: "自动取消 driver_arrived 或 in_progress 超过 2 小时未更新的订单，释放司机接单能力",
	}

	// 添加任务到系统中
	tasks := []*models.Task{orderTimeoutCancelTask, stuckOrderCancelTask}
	task.InitTasks(tasks)
}

// StuckOrderCancelHandler 处理卡单自动取消（driver_arrived / in_progress 长时间未更新）
func StuckOrderCancelHandler(ctx context.Context, params protocol.MapData) error {
	log.Get().Info("开始执行卡单自动取消任务...")

	hours := 2.0
	if v := params.GetFloat64("stuck_hours"); v > 0 {
		hours = v
	}
	cutoff := time.Now().Add(-time.Duration(hours * float64(time.Hour))).UnixMilli()

	var orderIDs []string
	err := models.DB.WithContext(ctx).
		Model(&models.Order{}).
		Select("order_id").
		Where("order_type = ?", protocol.RideOrder).
		Where("status IN ?", []string{protocol.StatusDriverArrived, protocol.StatusInProgress}).
		Where("updated_at < ?", cutoff).
		Pluck("order_id", &orderIDs).Error
	if err != nil {
		log.Get().Errorf("查询卡单订单失败: %v", err)
		return fmt.Errorf("查询卡单订单失败: %v", err)
	}
	if len(orderIDs) == 0 {
		log.Get().Info("没有找到需要取消的卡单订单")
		return nil
	}
	log.Get().Infof("找到 %d 个卡单订单", len(orderIDs))
	_, cancelErr := asyncCancelTimeoutOrders(ctx, orderIDs, "Order auto-cancelled: stuck in driver_arrived or in_progress for over "+fmt.Sprintf("%.0f", hours)+" hours")
	return cancelErr
}

// OrderTimeoutCancelHandler 处理订单超时取消任务
func OrderTimeoutCancelHandler(ctx context.Context, params protocol.MapData) error {
	log.Get().Info("开始执行订单超时取消任务...")

	// 获取超时时间配置（分钟）
	timeoutHours := 60.0
	if _v := params.GetFloat64("timeout_minutes"); _v > 0 {
		timeoutHours = _v
	}

	timeoutDuration := time.Duration(timeoutHours * float64(time.Minute))
	log.Get().Infof("使用超时时间: %v", timeoutDuration)

	// 处理预约订单超时
	scheduledCancelledCount, err := processTimeoutOrders(ctx, timeoutDuration)
	if err != nil {
		log.Get().Errorf("处理订单取消超时失败: %v", err)
		return fmt.Errorf("处理订单取消超时失败: %v", err)
	}

	log.Get().Infof("订单超时取消任务完成订单 %d 个", scheduledCancelledCount)

	return nil
}

// processTimeoutOrders 处理超时的预约订单
func processTimeoutOrders(ctx context.Context, timeoutDuration time.Duration) (int, error) {
	// 计算超时时间点
	timeoutTimestamp := time.Now().Add(-timeoutDuration).UnixMilli()

	// 查询超时的预约订单ID：只获取ID，提高查询效率
	var orderIDs []string
	err := models.DB.WithContext(ctx).
		Model(&models.Order{}).
		Select("order_id").
		Where("order_type = ?", protocol.RideOrder).
		Where("status = ?", protocol.StatusRequested).
		Where("scheduled_at < ?", timeoutTimestamp).
		Where("(provider_id IS NULL OR provider_id = '')").
		Pluck("order_id", &orderIDs).Error

	if err != nil {
		return 0, fmt.Errorf("查询超时预约订单失败: %v", err)
	}

	if len(orderIDs) == 0 {
		log.Get().Info("没有找到超时的预约订单")
		return 0, nil
	}

	log.Get().Infof("找到 %d 个超时的预约订单", len(orderIDs))

	// 异步轮询取消订单
	return asyncCancelTimeoutOrders(ctx, orderIDs, "Order timeout: no driver accepted within 1 hour after scheduled time")
}

// asyncCancelTimeoutOrders 异步轮询取消超时订单
func asyncCancelTimeoutOrders(ctx context.Context, orderIDs []string, reason string) (int, error) {
	if len(orderIDs) == 0 {
		return 0, nil
	}

	log.Get().Infof("开始异步处理 %d 个超时订单的取消操作", len(orderIDs))

	var wg sync.WaitGroup
	var cancelledCount int64
	orderService := GetOrderService()
	historyService := GetOrderHistoryService()

	for _, orderID := range orderIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			errCode := orderService.CancelOrder(id, "system", reason)
			// 取消订单
			if errCode != protocol.Success {
				log.Get().Errorf("取消超时订单失败: %s, 错误码: %v", id, errCode)
			}
			atomic.AddInt64(&cancelledCount, 1)
			log.Get().Infof("成功取消超时订单: %s", id)

			order := models.GetOrderByID(id)
			// 获取订单信息并处理通知
			if order != nil {
				// 记录取消历史
				if err := historyService.RecordOrderCancelled(order, order, "system", "system", reason); err != nil {
					log.Get().Errorf("记录订单取消历史失败，订单ID: %s, 错误: %v", id, err)
				}
				// 发送取消通知
				if err := orderService.NotifyOrderCancelled(id); err != nil {
					log.Get().Errorf("发送订单取消通知失败，订单ID: %s, 错误: %v", id, err)
				}
			}
		}(orderID)
	}

	wg.Wait()
	log.Get().Infof("成功取消 %v/%v 个订单,理由:%v", cancelledCount, len(orderIDs), reason)
	return int(cancelledCount), nil
}
