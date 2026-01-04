package task

import (
	"context"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// InitTaskHandlers 初始化任务处理器
func InitTaskHandlers() {
	log.Get().Info("初始化任务处理器...")

	// 注册信号处理器
	RegisterHandler(protocol.SignalProcessorHandler, func(ctx context.Context, params protocol.MapData) error {
		// 解析处理器逻辑
		return ProcessSignals(ctx, params)
	})
}

// RegisterSystemTasks 初始化系统任务
func RegisterSystemTasks() error {
	InitTaskHandlers()

	// 初始化或更新系统任务
	tasks := []*models.Task{
		{
			TaskID:     protocol.SignalProcessorHandler,
			HandlerKey: protocol.SignalProcessorHandler,
			Name:       "信号处理器",
			Cron:       "*/3 * * * *", // 每3秒扫描一次
			Timeout:    300,
			Status:     protocol.TaskStatusEnabled,
			Params:     map[string]any{},
		},
	}

	InitTasks(tasks)

	return nil
}

func InitTasks(tasks []*models.Task) {
	for _, task := range tasks {
		InitTask(task)
	}
}

// InitTask 初始化或更新系统任务
func InitTask(task *models.Task) error {
	// 查询任务是否存在
	if !models.CheckTaskExist(task.TaskID) {
		if err := models.CreateTask(task); err != nil {
			log.Get().Errorf("Create Task %v failed: %v \n", task.TaskID, err)
			return err
		}
		log.Get().Infof("Task %v created \n", task.TaskID)
	}

	return nil
}
