package task

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/models"
)

var (
	stopChan     chan struct{}
	scanInterval time.Duration
)

// Run 启动任务服务并注册优雅退出
func Run() error {
	cfg := config.Get().Task
	if cfg == nil {
		return nil
	}
	// 检查定时任务开关
	if !cfg.Enabled {
		return nil
	}

	// 初始化系统任务
	RegisterSystemTasks()

	// 初始化任务锁
	InitLock(models.GetRedis())

	// 设置扫描间隔
	scanInterval = time.Duration(cfg.ScanInterval) * time.Second

	// 初始化停止信号通道
	stopChan = make(chan struct{})

	// 启动任务扫描
	go scan()

	log.Get().Infof("Task service started, scan interval: %v", scanInterval)

	// 注册优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		Stop()
	}()

	return nil
}

// Stop 停止任务服务
func Stop() {
	if stopChan != nil {
		close(stopChan)
		log.Get().Info("Task service stopped")
	}
}

// scan 扫描并执行任务
func scan() {
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tasks, err := models.ScanTask()
			if err != nil {
				log.Get().Errorf("Scan task error: %v", err)
				continue
			}

			for _, task := range tasks {
				go RunTask(&task)
			}
		case <-stopChan:
			return
		}
	}
}
