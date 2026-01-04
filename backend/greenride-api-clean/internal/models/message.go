package models

import (
	"encoding/json"
	"greenride/internal/log"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// Message 消息实体
type Message struct {
	// 自增ID
	ID int64 `json:"id" gorm:"primaryKey;column:id;autoIncrement"`
	// 业务ID字段，唯一索引
	BizID string `json:"biz_id" gorm:"uniqueIndex:idx_message_biz_id;column:biz_id;type:varchar(64)"`
	// 会话ID，与checkpoint_id组合唯一
	ConversationID string `json:"conversation_id" gorm:"index:idx_conversation;uniqueIndex:idx_conv_checkpoint;column:conversation_id;type:varchar(64)"`
	// 用户ID
	UserID string `json:"user_id" gorm:"index:idx_message_user;column:user_id;type:varchar(64)"`
	// 线程ID，与checkpoint_id组合唯一
	ThreadID string `json:"thread_id" gorm:"index:idx_message_thread;uniqueIndex:idx_thread_checkpoint;column:thread_id;type:varchar(64)"`
	// 检查点ID，与conversation_id和thread_id组合唯一
	CheckpointID string `json:"checkpoint_id" gorm:"index:idx_checkpoint;uniqueIndex:idx_conv_checkpoint;uniqueIndex:idx_thread_checkpoint;column:checkpoint_id;type:varchar(64)"`
	// 消息来源
	Source string `json:"source" gorm:"column:source;type:varchar(32)"`
	// 步骤序号
	Step int `json:"step" gorm:"column:step;index:idx_message_step"`
	// 角色
	Role string `json:"role" gorm:"column:role;type:varchar(32)"`
	// 内容（文本数组）
	Content []string `json:"content" gorm:"column:content;type:json;serializer:json"`
	// 创建时间
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime;index:idx_message_created"`
	// 使用时间
	UseTime int64 `json:"use_time" gorm:"column:use_time"`
	// Token计数
	TokenCount int `json:"token_count" gorm:"column:token_count;default:0"`
	// 元数据
	Metadata map[string]interface{} `json:"metadata" gorm:"column:metadata;type:json;serializer:json"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "t_messages"
}

func NewMessage(userId, thread_id, checkpoint_id string) *Message {
	return &Message{
		BizID:          utils.GenerateUUID(),
		ConversationID: thread_id,
		UserID:         userId,
		ThreadID:       thread_id,
		CheckpointID:   checkpoint_id,
		Content:        []string{},
		Metadata:       make(map[string]interface{}),
	}
}

// InitMessageModel 初始化消息模型并进行自动迁移
func InitMessageModel() error {
	log.Get().Info("开始迁移 Message 表结构...")

	// 执行自动迁移，让GORM自动处理表结构
	if err := DB.AutoMigrate(&Message{}); err != nil {
		log.Get().Errorf("迁移 Message 表结构失败: %v", err)
		return err
	}

	// 创建必要的组合索引
	indexSQL := `
    -- 创建会话ID、步骤的组合索引(按消息顺序查询)
    CREATE INDEX IF NOT EXISTS idx_messages_conversation_step ON t_messages(conversation_id, step);
    
    -- 创建线程ID、创建时间的组合索引(按时间顺序查询)
    CREATE INDEX IF NOT EXISTS idx_messages_thread_created ON t_messages(thread_id, created_at DESC);
    
    -- 创建角色、会话ID的组合索引(按角色筛选消息)
    CREATE INDEX IF NOT EXISTS idx_messages_role_conversation ON t_messages(role, conversation_id);
    `

	if err := DB.Exec(indexSQL).Error; err != nil {
		log.Get().Warnf("创建组合索引失败(可能已存在): %v", err)
	}

	log.Get().Info("Message 表结构迁移完成")
	return nil
}

type Messages []*Message

func (t *Messages) ToInfo() []*protocol.Message {
	var infos []*protocol.Message
	for _, m := range *t {
		infos = append(infos, m.ToInfo())
	}
	return infos
}

func CheckMessageByThreadIDAndCheckpointerID(thread_id, checkpoint_id string) int64 {
	count := int64(0)
	err := DB.Where("thread_id = ? AND checkpoint_id = ?", thread_id, checkpoint_id).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}
func GetMessageByThreadIDAndCheckpointerID(thread_id, checkpoint_id string) *Message {
	var message Message
	err := DB.Where("thread_id = ? AND checkpoint_id = ?", thread_id, checkpoint_id).First(&message).Error
	if err != nil {
		return nil
	}
	return &message
}

// ToInfo 将 Message 对象转换为 MessageInfo
func (m *Message) ToInfo() *protocol.Message {
	metadataBytes, _ := json.Marshal(m.Metadata)
	metadataStr := string(metadataBytes)

	return &protocol.Message{
		ID:             m.ID,
		MessageID:      m.BizID,
		ConversationID: m.ConversationID,
		UserID:         m.UserID,
		ThreadID:       m.ThreadID,
		CheckpointID:   m.CheckpointID,
		Source:         m.Source,
		Step:           m.Step,
		Role:           m.Role,
		Content:        []string(m.Content),
		CreatedAt:      m.CreatedAt,
		UseTime:        m.UseTime,
		TokenCount:     m.TokenCount,
		Metadata:       metadataStr, // 序列化为JSON字符串
	}
}

// GetMessagesByConversationID retrieves messages for a given conversation ID with optional limit.
func GetMessagesByConversationID(conversationID string, limit int) (messages Messages) {
	db := DB.Model(&Message{}).
		Where("conversation_id = ?", conversationID).
		Order("step ASC")

	if limit > 0 {
		db = db.Limit(limit)
	}
	err := db.Find(&messages).Error
	if err != nil {
		log.Get().Errorf("获取消息失败: %v", err)
	}
	return
}

// GetMessageCountByConversationID retrieves the count of messages for a given conversation ID.
func GetMessageCountByConversationID(conversationID string) int64 {
	var messageCount int64
	if err := DB.Model(&Message{}).
		Where("conversation_id = ?", conversationID).
		Count(&messageCount).Error; err != nil {
		log.Get().Errorf("Failed to count messages: %v", err)
		return 0
	}
	return messageCount
}
