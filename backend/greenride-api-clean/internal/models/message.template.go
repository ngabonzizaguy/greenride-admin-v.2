package models

type MessageTemplate struct {
	ID          int64    `gorm:"column:id;primaryKey;autoIncrement"`
	TemplateID  string   `gorm:"column:template_id;type:varchar(64);uniqueIndex"`
	Type        string   `gorm:"column:type;type:varchar(32);"`
	Channel     string   `gorm:"column:channel;type:varchar(32);"`
	DeviceType  string   `gorm:"column:device_type;type:varchar(32);"`
	Platform    string   `gorm:"column:platform;type:varchar(32);"`
	Language    string   `gorm:"column:language;type:varchar(32);"`
	Region      string   `gorm:"column:region;type:varchar(32);"`
	Tags        []string `gorm:"column:tags;type:varchar(255);serializer:json"`
	Title       string   `gorm:"column:title;type:varchar(255)"`
	Status      string   `gorm:"column:status;type:varchar(32);default:''"`
	Description string   `gorm:"column:description;type:text"`
	Content     string   `gorm:"column:content;type:text"`
	Url         string   `gorm:"column:url;type:varchar(255)"`
	CreatedAt   int64    `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64    `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

func (m *MessageTemplate) TableName() string {
	return "t_message_template"
}
