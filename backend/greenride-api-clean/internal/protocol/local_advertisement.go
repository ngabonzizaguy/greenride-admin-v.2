package protocol

// LocalAdvertisementListRequest 本地广告列表请求
type LocalAdvertisementListRequest struct {
	City     string `json:"city,omitempty" form:"city"`         // 城市过滤
	Region   string `json:"region,omitempty" form:"region"`     // 地区过滤
	Category string `json:"category,omitempty" form:"category"` // 类别过滤
}

// LocalAdvertisementDetailRequest 本地广告详情请求
type LocalAdvertisementDetailRequest struct {
	AdID string `json:"ad_id" binding:"required"` // 广告ID
}

// LocalAdvertisement 本地广告响应
type LocalAdvertisement struct {
	AdID              string  `json:"ad_id"`
	Name              string  `json:"name"`
	Description       string  `json:"description,omitempty"`
	Category          string  `json:"category"`
	Address           string  `json:"address"`
	City              string  `json:"city,omitempty"`
	Region            string  `json:"region,omitempty"`
	Phone             string  `json:"phone,omitempty"`
	Website           string  `json:"website,omitempty"`
	ImageURL          string  `json:"image_url,omitempty"`
	LogoURL           string  `json:"logo_url,omitempty"`
	GoogleRating      float64 `json:"google_rating,omitempty"`
	GoogleReviewCount int     `json:"google_review_count,omitempty"`
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	IsPromoted        bool    `json:"is_promoted"`
	IsFeatured        bool    `json:"is_featured"`
	OpeningHours      string  `json:"opening_hours,omitempty"` // JSON字符串
	PriceLevel        int     `json:"price_level,omitempty"`
	Priority          int     `json:"priority"`
	DisplayOrder      int     `json:"display_order"`
}

// LocalAdvertisementStatsRequest 广告统计请求
type LocalAdvertisementStatsRequest struct {
	AdID      string `json:"ad_id" binding:"required"`
	StatsType string `json:"stats_type" binding:"required"` // view, click, call
	UserID    string `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}
