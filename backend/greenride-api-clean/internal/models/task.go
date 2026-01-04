package models

import (
	"time"

	"greenride/internal/protocol"
)

// Task 定时任务表
type Task struct {
	ID         int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskID     string           `json:"task_id" gorm:"column:task_id;type:varchar(64);uniqueIndex"` // 任务唯一标识
	Name       string           `json:"name" gorm:"column:name;type:varchar(256)"`                  // 任务名称
	Type       string           `json:"type" gorm:"column:type;type:varchar(64)"`                   // 任务类型
	HandlerKey string           `json:"handler_key" gorm:"column:handler_key;type:varchar(64)"`     // 处理器标识
	Cron       string           `json:"cron" gorm:"column:cron;type:varchar(32)"`                   // cron表达式
	Status     string           `json:"status" gorm:"column:status;type:varchar(32)"`               // 任务状态
	LastTime   int64            `json:"last_time" gorm:"column:last_time"`                          // 上次执行时间
	NextTime   int64            `json:"next_time" gorm:"column:next_time"`                          // 下次执行时间
	LastResult string           `json:"last_result" gorm:"column:last_result;type:varchar(64)"`     // 上次执行结果
	RetryCount int              `json:"retry_count" gorm:"column:retry_count"`                      // 重试次数
	MaxRetries int              `json:"max_retries" gorm:"column:max_retries"`                      // 最大重试次数
	Timeout    int              `json:"timeout" gorm:"column:timeout"`                              // 超时时间(秒)
	Params     protocol.MapData `json:"params" gorm:"column:params;type:json;serializer:json"`      // 任务参数
	Remark     string           `json:"remark" gorm:"column:remark;type:varchar(512)"`              // 备注说明
	CreatedAt  int64            `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

// TableName 表名
func (Task) TableName() string {
	return "t_task"
}

// ScanTask 扫描可执行的任务
func ScanTask() ([]Task, error) {
	var tasks []Task
	err := DB.Where("status = ? AND next_time>=0 and next_time <= ?",
		protocol.TaskStatusEnabled,
		time.Now().UnixMilli(),
	).Find(&tasks).Error

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// UpdateTaskExecution 更新任务执行信息
func UpdateTaskExecution(task *Task, result string, nextTime time.Time) error {
	now := time.Now()
	task.LastTime = now.UnixMilli()
	task.NextTime = nextTime.UnixMilli()
	task.LastResult = result
	return DB.Save(task).Error
}

func CheckTaskExist(taskID string) (isExist bool) {
	err := DB.Model(&Task{}).Where("task_id = ?", taskID).Select("task_id").First(&taskID).Error
	return err == nil
}

// GetTaskByTaskID 根据TaskID获取任务
func GetTaskByTaskID(taskID string) (*Task, error) {
	var task Task
	err := DB.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask 创建任务
func CreateTask(task *Task) error {
	return DB.Create(task).Error
}

// UpdateTask 更新任务
func UpdateTask(task *Task) error {
	return DB.Save(task).Error
}
