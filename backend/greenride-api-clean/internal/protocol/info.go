package protocol

// =============================================================================
// 通用信息实体 (Info Entities)
// =============================================================================

// Location 当前位置响应
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
	UpdatedAt int64   `json:"updated_at"`
}

// Rating 评价信息
type Rating struct {
	ID          int64    `json:"id"`
	OrderID     string   `json:"order_id"`
	RaterID     string   `json:"rater_id"`
	RateeID     string   `json:"ratee_id"`
	RaterType   string   `json:"rater_type"` // user, provider
	RateeType   string   `json:"ratee_type"` // user, provider
	Rating      float64  `json:"rating"`     // 评分 1-5
	Comment     string   `json:"comment,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Reply       string   `json:"reply,omitempty"`
	IsAnonymous bool     `json:"is_anonymous"`
	RepliedAt   *int64   `json:"replied_at,omitempty"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}
