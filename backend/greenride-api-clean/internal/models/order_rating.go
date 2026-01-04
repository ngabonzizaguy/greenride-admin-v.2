package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"slices"

	"gorm.io/gorm"
)

// OrderRating 订单评价模型
type OrderRating struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID     string         `gorm:"not null;index" json:"order_id"`                  // 订单ID
	RaterID     string         `gorm:"not null;index" json:"rater_id"`                  // 评价者ID
	RateeID     string         `gorm:"not null;index" json:"ratee_id"`                  // 被评价者ID
	RaterType   string         `gorm:"not null" json:"rater_type"`                      // 评价者类型: user, provider
	RateeType   string         `gorm:"not null" json:"ratee_type"`                      // 被评价者类型: user, provider
	Rating      float64        `gorm:"not null" json:"rating"`                          // 评分 1-5
	Comment     *string        `gorm:"type:text" json:"comment,omitempty"`              // 评价内容
	Tags        []string       `gorm:"type:json;serializer:json" json:"tags,omitempty"` // 评价标签数组
	Reply       *string        `gorm:"type:text" json:"reply,omitempty"`                // 回复内容
	IsAnonymous bool           `gorm:"default:false" json:"is_anonymous"`               // 是否匿名评价
	RepliedAt   *int64         `json:"replied_at,omitempty"`                            // 回复时间
	CreatedAt   int64          `gorm:"not null" json:"created_at"`                      // 创建时间
	UpdatedAt   int64          `gorm:"not null" json:"updated_at"`                      // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`               // 删除时间
}

// TableName 返回表名
func (OrderRating) TableName() string {
	return "t_order_ratings"
}

// BeforeCreate GORM钩子：创建前
func (r *OrderRating) BeforeCreate(tx *gorm.DB) error {
	if r.CreatedAt == 0 {
		r.CreatedAt = utils.TimeNowMilli()
	}
	if r.UpdatedAt == 0 {
		r.UpdatedAt = r.CreatedAt
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (r *OrderRating) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = utils.TimeNowMilli()
	return nil
}

// GetComment 安全获取评价内容
func (r *OrderRating) GetComment() string {
	if r.Comment == nil {
		return ""
	}
	return *r.Comment
}

// GetTags 安全获取评价标签
func (r *OrderRating) GetTags() []string {
	if r.Tags == nil {
		return []string{}
	}
	return r.Tags
}

// GetReply 安全获取回复内容
func (r *OrderRating) GetReply() string {
	if r.Reply == nil {
		return ""
	}
	return *r.Reply
}

// SetComment 设置评价内容
func (r *OrderRating) SetComment(comment string) {
	if comment == "" {
		r.Comment = nil
	} else {
		r.Comment = &comment
	}
}

// SetTags 设置评价标签
func (r *OrderRating) SetTags(tags []string) {
	if len(tags) == 0 {
		r.Tags = nil
	} else {
		r.Tags = tags
	}
}

// AddTag 添加单个标签
func (r *OrderRating) AddTag(tag string) {
	if tag == "" {
		return
	}

	// 检查标签是否已存在
	for _, existingTag := range r.Tags {
		if existingTag == tag {
			return
		}
	}

	r.Tags = append(r.Tags, tag)
}

// RemoveTag 移除单个标签
func (r *OrderRating) RemoveTag(tag string) {
	if tag == "" {
		return
	}

	for i, existingTag := range r.Tags {
		if existingTag == tag {
			r.Tags = append(r.Tags[:i], r.Tags[i+1:]...)
			break
		}
	}
}

// HasTag 检查是否包含特定标签
func (r *OrderRating) HasTag(tag string) bool {
	return slices.Contains(r.Tags, tag)
}

// GetTagsCount 获取标签数量
func (r *OrderRating) GetTagsCount() int {
	return len(r.Tags)
}

// SetReply 设置回复内容
func (r *OrderRating) SetReply(reply string) {
	if reply == "" {
		r.Reply = nil
	} else {
		r.Reply = &reply
		timestamp := utils.TimeNowMilli()
		r.RepliedAt = &timestamp
	}
}

// IsReplied 检查是否已回复
func (r *OrderRating) IsReplied() bool {
	return r.Reply != nil && *r.Reply != ""
}

// CanBeRepliedBy 检查用户是否可以回复此评价
func (r *OrderRating) CanBeRepliedBy(userID, userType string) bool {
	// 被评价者可以回复
	if r.RateeID == userID && r.RateeType == userType {
		return true
	}

	// 管理员可以回复
	if userType == "admin" {
		return true
	}

	return false
}

// CanBeEditedBy 检查用户是否可以编辑此评价
func (r *OrderRating) CanBeEditedBy(userID, userType string) bool {
	// 评价者可以编辑自己的评价
	return r.RaterID == userID && r.RaterType == userType
}

// CanBeDeletedBy 检查用户是否可以删除此评价
func (r *OrderRating) CanBeDeletedBy(userID, userType string) bool {
	// 评价者可以删除自己的评价
	if r.RaterID == userID && r.RaterType == userType {
		return true
	}

	// 管理员可以删除任何评价
	if userType == "admin" {
		return true
	}

	return false
}

func GetRatingsByOrderID(orderID string) []*OrderRating {
	var ratings []*OrderRating
	if err := GetDB().Where("order_id = ?", orderID).Find(&ratings).Error; err != nil {
		return nil
	}
	return ratings
}

type OrderRatings []*OrderRating

func (r *OrderRatings) Protocol() []*protocol.Rating {
	var result []*protocol.Rating
	for _, rating := range *r {
		result = append(result, rating.Protocol())
	}
	return result
}

// GetRatingLevel 获取评分等级
func (r *OrderRating) GetRatingLevel() string {
	switch {
	case r.Rating >= 4.5:
		return "excellent"
	case r.Rating >= 3.5:
		return "good"
	case r.Rating >= 2.5:
		return "average"
	case r.Rating >= 1.5:
		return "poor"
	default:
		return "terrible"
	}
}

// IsPositive 是否为正面评价
func (r *OrderRating) IsPositive() bool {
	return r.Rating >= 3.5
}

// IsNegative 是否为负面评价
func (r *OrderRating) IsNegative() bool {
	return r.Rating < 2.5
}

func GetOrderRatingsByID(ratingID string) *OrderRating {
	var rating OrderRating
	if err := GetDB().Where("id = ?", ratingID).First(&rating).Error; err != nil {
		return nil
	}
	return &rating
}

func GetOrderRatingsByOrderID(orderID string) []*OrderRating {
	var ratings []*OrderRating
	if err := GetDB().Where("order_id = ?", orderID).Find(&ratings).Error; err != nil {
		return nil
	}
	return ratings
}

// ToProtocol 转换为协议结构
func (r *OrderRating) Protocol() *protocol.Rating {
	return &protocol.Rating{
		ID:          r.ID,
		OrderID:     r.OrderID,
		RaterID:     r.RaterID,
		RateeID:     r.RateeID,
		RaterType:   r.RaterType,
		RateeType:   r.RateeType,
		Rating:      r.Rating,
		Comment:     r.GetComment(),
		Tags:        r.GetTags(),
		Reply:       r.GetReply(),
		IsAnonymous: r.IsAnonymous,
		RepliedAt:   r.RepliedAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
