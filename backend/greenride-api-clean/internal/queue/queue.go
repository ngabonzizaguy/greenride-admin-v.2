package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"

	"greenride/internal/config"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
)

// TaskType 定义任务类型
const (
	TypeSendNotification = "notification:send"
	TypeSendFCM          = "fcm:send"
	TypeProcessPayment   = "payment:process"
	TypeSendSMS          = "sms:send"
	TypeSendEmail        = "email:send"
	TypeGenerateReport   = "report:generate"
	TypeCleanupData      = "cleanup:data"
)

// NotificationPayload 通知任务载荷
type NotificationPayload struct {
	UserID  uint                   `json:"user_id"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// FCMPayload FCM推送载荷
type FCMPayload struct {
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
	ImageURL    string            `json:"image_url,omitempty"`
	ClickAction string            `json:"click_action,omitempty"`
	OrderID     string            `json:"order_id,omitempty"`     // 订单相关推送
	MessageType string            `json:"message_type,omitempty"` // 消息类型
}

// PaymentPayload 支付处理载荷
type PaymentPayload struct {
	TripID        uint    `json:"trip_id"`
	CustomerID    uint    `json:"customer_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

// SMSPayload 短信载荷
type SMSPayload struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// EmailPayload 邮件载荷
type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  bool   `json:"is_html"`
}

// QueueManager 队列管理器
type QueueManager struct {
	client *asynq.Client
	server *asynq.Server
}

// QueueManagerInterface 定义队列管理器接口
type QueueManagerInterface interface {
	EnqueueFCM(payload map[string]interface{}, priority string, delay time.Duration) error
}

// NewQueueManager 创建队列管理器
func NewQueueManager(cfg *config.Config) *QueueManager {
	// 解析Redis DSN为asynq可用的选项
	dsn := cfg.Redis.Dsn
	if dsn == "" {
		panic("Redis DSN is empty")
	}

	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse Redis DSN: %v", err))
	}

	// 构建asynq Redis连接选项
	redisOpt := asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	}

	return &QueueManager{
		client: asynq.NewClient(redisOpt),
		server: asynq.NewServer(
			redisOpt,
			asynq.Config{
				Concurrency: 10,
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
				ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
					log.Printf("Task failed: %s, Error: %v", task.Type(), err)
				}),
			},
		),
	}
}

// NewQueueClient 创建队列客户端（简化版本）
func NewQueueClient(cfg *config.Config) *QueueManager {
	return NewQueueManager(cfg)
}

// Start 启动队列处理
func (qm *QueueManager) Start() error {
	mux := asynq.NewServeMux()

	// 注册任务处理器
	mux.HandleFunc(TypeSendNotification, HandleSendNotification)
	mux.HandleFunc(TypeSendFCM, HandleSendFCM)
	mux.HandleFunc(TypeProcessPayment, HandleProcessPayment)
	mux.HandleFunc(TypeSendSMS, HandleSendSMS)
	mux.HandleFunc(TypeSendEmail, HandleSendEmail)
	mux.HandleFunc(TypeGenerateReport, HandleGenerateReport)
	mux.HandleFunc(TypeCleanupData, HandleCleanupData)

	return qm.server.Run(mux)
}

// Stop 停止队列处理
func (qm *QueueManager) Stop() {
	qm.server.Shutdown()
	qm.client.Close()
}

// EnqueueNotification 添加通知任务到队列
func (qm *QueueManager) EnqueueNotification(payload NotificationPayload, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendNotification, data)
	opts := []asynq.Option{
		asynq.Queue("default"),
		asynq.MaxRetry(3),
	}

	if delay > 0 {
		opts = append(opts, asynq.ProcessIn(delay))
	}

	_, err = qm.client.Enqueue(task, opts...)
	return err
}

// EnqueuePayment 添加支付处理任务到队列
func (qm *QueueManager) EnqueuePayment(payload PaymentPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeProcessPayment, data)
	_, err = qm.client.Enqueue(task,
		asynq.Queue("critical"),
		asynq.MaxRetry(5),
		asynq.Timeout(5*time.Minute),
	)
	return err
}

// EnqueueSMS 添加短信发送任务到队列
func (qm *QueueManager) EnqueueSMS(payload SMSPayload, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendSMS, data)
	opts := []asynq.Option{
		asynq.Queue("default"),
		asynq.MaxRetry(3),
	}

	if delay > 0 {
		opts = append(opts, asynq.ProcessIn(delay))
	}

	_, err = qm.client.Enqueue(task, opts...)
	return err
}

// EnqueueFCMPayload 添加FCM推送任务到队列（结构体版本）
func (qm *QueueManager) EnqueueFCMPayload(payload FCMPayload, priority string, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendFCM, data)

	// 根据消息类型确定队列优先级
	queue := "default"
	maxRetries := 3

	switch priority {
	case "critical":
		queue = "critical"
		maxRetries = 5
	case "high":
		queue = "default"
		maxRetries = 4
	case "low":
		queue = "low"
		maxRetries = 2
	}

	opts := []asynq.Option{
		asynq.Queue(queue),
		asynq.MaxRetry(maxRetries),
		asynq.Timeout(30 * time.Second),
	}

	if delay > 0 {
		opts = append(opts, asynq.ProcessIn(delay))
	}

	_, err = qm.client.Enqueue(task, opts...)
	return err
}

// EnqueueFCM 实现接口版本（接受 map[string]interface{}）
func (qm *QueueManager) EnqueueFCM(payload map[string]interface{}, priority string, delay time.Duration) error {
	// 转换 payload 到 FCMPayload 结构
	fcmPayload := FCMPayload{}

	if userID, ok := payload["user_id"].(string); ok {
		fcmPayload.UserID = userID
	}
	if title, ok := payload["title"].(string); ok {
		fcmPayload.Title = title
	}
	if body, ok := payload["body"].(string); ok {
		fcmPayload.Body = body
	}
	if data, ok := payload["data"].(map[string]string); ok {
		fcmPayload.Data = data
	}
	if messageType, ok := payload["message_type"].(string); ok {
		fcmPayload.MessageType = messageType
	}

	return qm.EnqueueFCMPayload(fcmPayload, priority, delay)
}

// EnqueueEmail 添加邮件发送任务到队列
func (qm *QueueManager) EnqueueEmail(payload EmailPayload, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendEmail, data)
	opts := []asynq.Option{
		asynq.Queue("default"),
		asynq.MaxRetry(3),
	}

	if delay > 0 {
		opts = append(opts, asynq.ProcessIn(delay))
	}

	_, err = qm.client.Enqueue(task, opts...)
	return err
}

// HandleSendFCM 处理FCM推送任务
// HandleSendFCM 处理发送FCM任务（调用Firebase服务）
func HandleSendFCM(ctx context.Context, t *asynq.Task) error {
	var payload FCMPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal FCM payload: %v", err)
	}

	log.Printf("Processing FCM task for user %s: %s", payload.UserID, payload.Title)

	// 直接调用Firebase服务
	firebaseService := services.GetFirebaseService()
	if firebaseService == nil {
		return fmt.Errorf("firebase service not available")
	}

	// 从数据库获取用户的FCM tokens
	var tokens []string
	err := models.GetDB().Model(&models.FCMToken{}).
		Where("user_id = ? AND status = ?", payload.UserID, protocol.StatusActive).
		Pluck("token", &tokens).Error

	if err != nil {
		return fmt.Errorf("failed to get FCM tokens for user %s: %v", payload.UserID, err)
	}

	if len(tokens) == 0 {
		log.Printf("No valid FCM tokens found for user %s", payload.UserID)
		return nil // 不视为错误，只是没有设备可推送
	}

	// 准备FCM消息
	fcmMsg := &protocol.FCMMessage{
		UserID: payload.UserID,
		Title:  payload.Title,
		Body:   payload.Body,
		Data:   payload.Data,
	}

	// 根据token数量选择发送方式
	var sendErr error
	if len(tokens) == 1 {
		// 单个token，使用SendFCMMessage
		fcmMsg.Token = tokens[0]
		_, sendErr = firebaseService.SendFCMMessage(ctx, fcmMsg)
	} else {
		// 多个token，使用SendMulticastFCMMessage
		_, sendErr = firebaseService.SendMulticastFCMMessage(ctx, fcmMsg, tokens)
	}

	if sendErr != nil {
		log.Printf("Failed to send FCM message to user %s: %v", payload.UserID, sendErr)
		return sendErr
	}

	log.Printf("FCM message sent successfully to user %s", payload.UserID)
	return nil
}

// HandleSendNotification 处理发送通知任务
func HandleSendNotification(ctx context.Context, t *asynq.Task) error {
	var payload NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Sending notification to user %d: %s", payload.UserID, payload.Message)

	// 这里实现实际的通知发送逻辑
	// 例如：FCM推送、WebSocket推送等

	return nil
}

// HandleProcessPayment 处理支付任务
func HandleProcessPayment(ctx context.Context, t *asynq.Task) error {
	var payload PaymentPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Processing payment for trip %d, amount: %.2f", payload.TripID, payload.Amount)

	// 这里实现实际的支付处理逻辑
	// 例如：调用支付网关API

	return nil
}

// HandleSendSMS 处理发送短信任务
func HandleSendSMS(ctx context.Context, t *asynq.Task) error {
	var payload SMSPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Sending SMS to %s: %s", payload.Phone, payload.Message)

	// 这里实现实际的短信发送逻辑
	// 例如：调用阿里云短信服务API

	return nil
}

// HandleSendEmail 处理发送邮件任务
func HandleSendEmail(ctx context.Context, t *asynq.Task) error {
	var payload EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Sending email to %s: %s", payload.To, payload.Subject)

	// 这里实现实际的邮件发送逻辑
	// 例如：使用SMTP或第三方邮件服务

	return nil
}

// HandleGenerateReport 处理生成报告任务
func HandleGenerateReport(ctx context.Context, t *asynq.Task) error {
	log.Printf("Generating report...")

	// 这里实现报告生成逻辑
	// 例如：从数据库查询数据，生成PDF或Excel

	return nil
}

// HandleCleanupData 处理数据清理任务
func HandleCleanupData(ctx context.Context, t *asynq.Task) error {
	log.Printf("Cleaning up data...")

	// 这里实现数据清理逻辑
	// 例如：删除过期的临时文件、日志等

	return nil
}

// NewTaskProcessor 创建任务处理器
func NewTaskProcessor(cfg *config.Config, db interface{}, cache interface{}) *QueueManager {
	// 简化实现，返回标准的QueueManager
	return NewQueueManager(cfg)
}

// StartQueueServer 启动队列服务器
func StartQueueServer(cfg *config.Config, processor *QueueManager) {
	log.Println("Queue server started")
	if err := processor.Start(); err != nil {
		log.Fatal("Failed to start queue server:", err)
	}
}
