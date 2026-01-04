package task

import (
	"context"
	"fmt"
	"time"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/utils"
)

// RunTask 执行单个任务
func RunTask(task *models.Task) error {
	// 获取任务锁
	mutex, err := Lock(task.TaskID)
	if err != nil {
		return fmt.Errorf("failed to acquire task lock: %v", err)
	}
	defer Unlock(mutex)

	// 查找任务处理器
	handler, exists := GetHandler(task.HandlerKey)
	if !exists {
		return fmt.Errorf("handler not found for key: %s", task.HandlerKey)
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.Timeout)*time.Second)
	defer cancel()

	// 执行任务并记录日志
	startTime := time.Now()
	//log.Get().Infof("Task started: %s", task.TaskID)

	err = handler(ctx, task.Params)

	// 计算下次执行时间
	nextTime := utils.CalculateNextTime(task.Cron)

	// 准备执行结果
	result := "success"
	if err != nil {
		result = fmt.Sprintf("failed: %v", err)
		log.Get().Errorf("Task failed: %s, error: %v", task.TaskID, err)
	}

	// 更新任务执行信息
	if err := models.UpdateTaskExecution(task, result, nextTime); err != nil {
		log.Get().Errorf("Failed to update task execution: %v", err)
	}

	log.Get().Infof("Task completed: %s, duration: %v", task.TaskID, time.Since(startTime))
	return err
}
