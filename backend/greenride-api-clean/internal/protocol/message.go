package protocol

import "html/template"

// Message 消息信息对象，用于API返回
type Message struct {
	ID                int64    `json:"id"`                           // 数据库ID
	MessageID         string   `json:"message_id"`                   // 消息ID(对应model中的biz_id)
	ConversationID    string   `json:"conversation_id"`              // 对话ID
	ConversationTitle string   `json:"conversation_title,omitempty"` // 对话标题(非必须)
	UserID            string   `json:"user_id"`                      // 用户ID
	ThreadID          string   `json:"thread_id"`                    // 线程ID
	CheckpointID      string   `json:"checkpoint_id"`                // 检查点ID
	Source            string   `json:"source"`                       // 消息来源
	Step              int      `json:"step"`                         // 步骤序号
	Role              string   `json:"role"`                         // 角色
	Content           []string `json:"content"`                      // 内容
	CreatedAt         int64    `json:"created_at"`                   // 创建时间
	UseTime           int64    `json:"use_time"`                     // 使用时间
	TokenCount        int      `json:"token_count"`                  // Token计数
	Metadata          string   `json:"metadata"`                     // 元数据
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	Total    int64      `json:"total"`    // 总记录数
	Messages []*Message `json:"messages"` // 消息列表
}

type MessageTemplate struct {
	ID          int64              `gorm:"column:id;primaryKey;autoIncrement"`
	TemplateID  string             `gorm:"column:template_id;type:varchar(64);not null;uniqueIndex"`
	Type        string             `gorm:"column:type;type:varchar(32);not null"`
	Channel     string             `gorm:"column:channel;type:varchar(32);not null"`
	DeviceType  string             `gorm:"column:device_type;type:varchar(32);not null"`
	Platform    string             `gorm:"column:platform;type:varchar(32);not null"`
	Language    string             `gorm:"column:language;type:varchar(32);not null"`
	Region      string             `gorm:"column:region;type:varchar(32);not null"`
	Tags        []string           `gorm:"column:tags;type:varchar(255);serializer:json"`
	Title       *template.Template `gorm:"column:title;type:varchar(255)"`
	Status      string             `gorm:"column:status;type:varchar(32);default:''"`
	Description string             `gorm:"column:description;type:text"`
	Content     *template.Template `gorm:"-"`
}
