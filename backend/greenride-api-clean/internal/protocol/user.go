package protocol

import (
	"fmt"
	"net/url"
)

// GenerateDefaultAvatar 基于用户ID生成默认头像URL
func GenerateDefaultAvatar(userID string) string {
	// 基础URL
	baseURL := "https://api.dicebear.com/7.x/avataaars/svg"

	// 参数配置
	params := url.Values{}
	params.Set("seed", userID)              // 使用用户ID作为种子
	params.Set("size", "200")               // 头像大小200px
	params.Set("backgroundColor", "ffffff") // 白色背景

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// User 用户信息响应结构
type User struct {
	UserID      string `json:"user_id"`
	UserType    string `json:"user_type"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	FullName    string `json:"full_name"`
	Avatar      string `json:"avatar,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Birthday    int64  `json:"birthday,omitempty"`
	Language    string `json:"language"`
	Timezone    string `json:"timezone"`

	// 联系信息
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`

	// 位置信息
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	LocationUpdatedAt int64   `json:"location_updated_at,omitempty"`

	// 状态信息
	Status          string `json:"status"`
	IsEmailVerified bool   `json:"is_email_verified"`
	IsPhoneVerified bool   `json:"is_phone_verified"`
	IsActive        bool   `json:"is_active"`
	OnlineStatus    string `json:"online_status"`

	// 司机相关信息
	LicenseNumber string  `json:"license_number,omitempty"`
	LicenseExpiry int64   `json:"license_expiry,omitempty"`
	Score         float64 `json:"score,omitempty"`
	TotalRides    int     `json:"total_rides,omitempty"`

	// 推荐系统
	InviteCode  string `json:"invite_code,omitempty"`
	InviteCount int    `json:"invite_count,omitempty"`

	// 沙盒信息
	Sandbox int `json:"sandbox"`

	// 时间戳
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
	DeletedAt       int64    `json:"deleted_at,omitempty"` // 删除时间戳
	LastLoginAt     int64    `json:"last_login_at,omitempty"`
	EmailVerifiedAt int64    `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt int64    `json:"phone_verified_at,omitempty"`
	Vehicle         *Vehicle `json:"vehicle,omitempty"`
}
