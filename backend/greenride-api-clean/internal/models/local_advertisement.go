package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"

	"gorm.io/gorm"
)

// LocalAdvertisement 本地广告表 - 管理首页显示的本地商家广告
type LocalAdvertisement struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AdID string `json:"ad_id" gorm:"column:ad_id;type:varchar(64);uniqueIndex"`
	Salt string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*LocalAdvertisementValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type LocalAdvertisementValues struct {
	// 基本信息
	Name        *string `json:"name" gorm:"column:name;type:varchar(255)"`               // 商家名称
	Description *string `json:"description" gorm:"column:description;type:text"`         // 商家描述
	Category    *string `json:"category" gorm:"column:category;type:varchar(100);index"` // 商家类别

	// 位置信息
	Address   *string  `json:"address" gorm:"column:address;type:varchar(500)"`            // 详细地址
	City      *string  `json:"city" gorm:"column:city;type:varchar(100);index"`            // 城市
	Region    *string  `json:"region" gorm:"column:region;type:varchar(100);index"`        // 地区
	Country   *string  `json:"country" gorm:"column:country;type:varchar(100);index"`      // 国家
	Latitude  *float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8);index"`   // 纬度
	Longitude *float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8);index"` // 经度

	// 联系信息
	Phone   *string `json:"phone" gorm:"column:phone;type:varchar(20)"`      // 电话号码
	Email   *string `json:"email" gorm:"column:email;type:varchar(255)"`     // 邮箱
	Website *string `json:"website" gorm:"column:website;type:varchar(500)"` // 网站URL

	// 媒体信息
	ImageURL  *string `json:"image_url" gorm:"column:image_url;type:varchar(1024)"` // 主图片URL (支持长URL)
	ImageURLs *string `json:"image_urls" gorm:"column:image_urls;type:json"`        // 多张图片URLs (JSON数组)
	LogoURL   *string `json:"logo_url" gorm:"column:logo_url;type:varchar(1024)"`   // Logo URL

	// Google Maps 相关信息
	GooglePlaceID     *string  `json:"google_place_id" gorm:"column:google_place_id;type:varchar(255);index"` // Google Place ID
	GoogleRating      *float64 `json:"google_rating" gorm:"column:google_rating;type:decimal(3,2)"`           // Google评分
	GoogleReviewCount *int     `json:"google_review_count" gorm:"column:google_review_count"`                 // Google评论数量

	// 营业信息
	OpeningHours  *string `json:"opening_hours" gorm:"column:opening_hours;type:json"`           // 营业时间 (JSON对象)
	IsOpen24Hours *bool   `json:"is_open_24_hours" gorm:"column:is_open_24_hours;default:false"` // 是否24小时营业
	PriceLevel    *int    `json:"price_level" gorm:"column:price_level"`                         // 价格等级 (0-4)

	// 显示设置
	Status       *string `json:"status" gorm:"column:status;type:varchar(20);index;default:'active'"` // active, inactive, pending
	Priority     *int    `json:"priority" gorm:"column:priority;default:0"`                           // 显示优先级，数字越大优先级越高
	DisplayOrder *int    `json:"display_order" gorm:"column:display_order;default:0"`                 // 显示顺序
	IsPromoted   *bool   `json:"is_promoted" gorm:"column:is_promoted;default:false"`                 // 是否为推广广告
	IsFeatured   *bool   `json:"is_featured" gorm:"column:is_featured;default:false"`                 // 是否为精选广告

	// 地理限制
	ServiceRadius *float64 `json:"service_radius" gorm:"column:service_radius"`                   // 服务半径(公里)
	TargetCities  *string  `json:"target_cities" gorm:"column:target_cities;type:varchar(500)"`   // 目标城市 (逗号分隔)
	TargetRegions *string  `json:"target_regions" gorm:"column:target_regions;type:varchar(500)"` // 目标地区 (逗号分隔)

	// 时间控制
	StartAt *int64 `json:"start_at" gorm:"column:start_at"` // 开始展示时间 (毫秒时间戳)
	EndAt   *int64 `json:"end_at" gorm:"column:end_at"`     // 结束展示时间 (毫秒时间戳)

	// 统计信息
	ViewCount  *int `json:"view_count" gorm:"column:view_count;default:0"`   // 浏览次数
	ClickCount *int `json:"click_count" gorm:"column:click_count;default:0"` // 点击次数
	CallCount  *int `json:"call_count" gorm:"column:call_count;default:0"`   // 电话点击次数

	// 管理信息
	CreatedBy *string `json:"created_by" gorm:"column:created_by;type:varchar(64)"` // 创建者ID
	UpdatedBy *string `json:"updated_by" gorm:"column:updated_by;type:varchar(64)"` // 更新者ID

	// 扩展信息
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"`         // 标签 (逗号分隔)
	Keywords *string `json:"keywords" gorm:"column:keywords;type:varchar(500)"` // 关键词 (逗号分隔)
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`               // 备注

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (LocalAdvertisement) TableName() string {
	return "t_local_advertisements"
}

// 创建新的本地广告对象
func NewLocalAdvertisement() *LocalAdvertisement {
	return &LocalAdvertisement{
		AdID: utils.GenerateID(),
		Salt: utils.GenerateSalt(),
		LocalAdvertisementValues: &LocalAdvertisementValues{
			Status:       utils.StringPtr(protocol.StatusPending),
			Priority:     utils.IntPtr(0),
			DisplayOrder: utils.IntPtr(0),
			IsPromoted:   utils.BoolPtr(false),
			IsFeatured:   utils.BoolPtr(false),
			ViewCount:    utils.IntPtr(0),
			ClickCount:   utils.IntPtr(0),
			CallCount:    utils.IntPtr(0),
		},
	}
}

// Getter 方法
func (la *LocalAdvertisementValues) GetName() string {
	if la.Name == nil {
		return ""
	}
	return *la.Name
}

func (la *LocalAdvertisementValues) GetDescription() string {
	if la.Description == nil {
		return ""
	}
	return *la.Description
}

func (la *LocalAdvertisementValues) GetAddress() string {
	if la.Address == nil {
		return ""
	}
	return *la.Address
}

func (la *LocalAdvertisementValues) GetPhone() string {
	if la.Phone == nil {
		return ""
	}
	return *la.Phone
}

func (la *LocalAdvertisementValues) GetWebsite() string {
	if la.Website == nil {
		return ""
	}
	return *la.Website
}

func (la *LocalAdvertisementValues) GetImageURL() string {
	if la.ImageURL == nil {
		return ""
	}
	return *la.ImageURL
}

func (la *LocalAdvertisementValues) GetLogoURL() string {
	if la.LogoURL == nil {
		return ""
	}
	return *la.LogoURL
}

func (la *LocalAdvertisementValues) GetCategory() string {
	if la.Category == nil {
		return protocol.LocalAdCategoryOther
	}
	return *la.Category
}

func (la *LocalAdvertisementValues) GetStatus() string {
	if la.Status == nil {
		return protocol.StatusPending
	}
	return *la.Status
}

func (la *LocalAdvertisementValues) GetPriority() int {
	if la.Priority == nil {
		return 0
	}
	return *la.Priority
}

func (la *LocalAdvertisementValues) GetDisplayOrder() int {
	if la.DisplayOrder == nil {
		return 0
	}
	return *la.DisplayOrder
}

func (la *LocalAdvertisementValues) GetGoogleRating() float64 {
	if la.GoogleRating == nil {
		return 0.0
	}
	return *la.GoogleRating
}

func (la *LocalAdvertisementValues) GetGoogleReviewCount() int {
	if la.GoogleReviewCount == nil {
		return 0
	}
	return *la.GoogleReviewCount
}

func (la *LocalAdvertisementValues) GetLatitude() float64 {
	if la.Latitude == nil {
		return 0.0
	}
	return *la.Latitude
}

func (la *LocalAdvertisementValues) GetLongitude() float64 {
	if la.Longitude == nil {
		return 0.0
	}
	return *la.Longitude
}

func (la *LocalAdvertisementValues) GetCity() string {
	if la.City == nil {
		return ""
	}
	return *la.City
}

func (la *LocalAdvertisementValues) GetRegion() string {
	if la.Region == nil {
		return ""
	}
	return *la.Region
}

func (la *LocalAdvertisementValues) GetOpeningHours() string {
	if la.OpeningHours == nil {
		return ""
	}
	return *la.OpeningHours
}

func (la *LocalAdvertisementValues) GetPriceLevel() int {
	if la.PriceLevel == nil {
		return 0
	}
	return *la.PriceLevel
}

// Setter 方法
func (la *LocalAdvertisementValues) SetName(name string) *LocalAdvertisementValues {
	la.Name = &name
	return la
}

func (la *LocalAdvertisementValues) SetDescription(description string) *LocalAdvertisementValues {
	la.Description = &description
	return la
}

func (la *LocalAdvertisementValues) SetCategory(category string) *LocalAdvertisementValues {
	la.Category = &category
	return la
}

func (la *LocalAdvertisementValues) SetAddress(address string) *LocalAdvertisementValues {
	la.Address = &address
	return la
}

func (la *LocalAdvertisementValues) SetLocation(address, city, region, country string, lat, lng float64) *LocalAdvertisementValues {
	la.Address = &address
	la.City = &city
	la.Region = &region
	la.Country = &country
	la.Latitude = &lat
	la.Longitude = &lng
	return la
}

func (la *LocalAdvertisementValues) SetContact(phone, email, website string) *LocalAdvertisementValues {
	if phone != "" {
		la.Phone = &phone
	}
	if email != "" {
		la.Email = &email
	}
	if website != "" {
		la.Website = &website
	}
	return la
}

func (la *LocalAdvertisementValues) SetImages(imageURL, logoURL string) *LocalAdvertisementValues {
	if imageURL != "" {
		la.ImageURL = &imageURL
	}
	if logoURL != "" {
		la.LogoURL = &logoURL
	}
	return la
}

func (la *LocalAdvertisementValues) SetGoogleInfo(placeID string, rating float64, reviewCount int) *LocalAdvertisementValues {
	if placeID != "" {
		la.GooglePlaceID = &placeID
	}
	if rating > 0 {
		la.GoogleRating = &rating
	}
	if reviewCount > 0 {
		la.GoogleReviewCount = &reviewCount
	}
	return la
}

func (la *LocalAdvertisementValues) SetStatus(status string) *LocalAdvertisementValues {
	la.Status = &status
	return la
}

func (la *LocalAdvertisementValues) SetPriority(priority int) *LocalAdvertisementValues {
	la.Priority = &priority
	return la
}

func (la *LocalAdvertisementValues) SetDisplayOrder(order int) *LocalAdvertisementValues {
	la.DisplayOrder = &order
	return la
}

func (la *LocalAdvertisementValues) SetPromoted(promoted bool) *LocalAdvertisementValues {
	la.IsPromoted = &promoted
	return la
}

func (la *LocalAdvertisementValues) SetFeatured(featured bool) *LocalAdvertisementValues {
	la.IsFeatured = &featured
	return la
}

func (la *LocalAdvertisementValues) SetCreator(createdBy string) *LocalAdvertisementValues {
	la.CreatedBy = &createdBy
	return la
}

// 业务方法
func (la *LocalAdvertisement) IsActive() bool {
	return la.GetStatus() == protocol.StatusActive
}

func (la *LocalAdvertisement) IsTimeValid() bool {
	now := utils.TimeNowMilli()

	// 检查开始时间
	if la.StartAt != nil && now < *la.StartAt {
		return false
	}

	// 检查结束时间
	if la.EndAt != nil && now > *la.EndAt {
		return false
	}

	return true
}

func (la *LocalAdvertisement) IsVisible() bool {
	return la.IsActive() && la.IsTimeValid()
}

func (la *LocalAdvertisement) IsPromotedAd() bool {
	return la.IsPromoted != nil && *la.IsPromoted
}

func (la *LocalAdvertisement) IsFeaturedAd() bool {
	return la.IsFeatured != nil && *la.IsFeatured
}

// Protocol 将模型转换为协议响应格式
func (la *LocalAdvertisement) Protocol() protocol.LocalAdvertisement {
	return protocol.LocalAdvertisement{
		AdID:              la.AdID,
		Name:              la.GetName(),
		Description:       la.GetDescription(),
		Category:          la.GetCategory(),
		Address:           la.GetAddress(),
		City:              la.GetCity(),
		Region:            la.GetRegion(),
		Phone:             la.GetPhone(),
		Website:           la.GetWebsite(),
		ImageURL:          la.GetImageURL(),
		LogoURL:           la.GetLogoURL(),
		GoogleRating:      la.GetGoogleRating(),
		GoogleReviewCount: la.GetGoogleReviewCount(),
		Latitude:          la.GetLatitude(),
		Longitude:         la.GetLongitude(),
		IsPromoted:        la.IsPromotedAd(),
		IsFeatured:        la.IsFeaturedAd(),
		OpeningHours:      la.GetOpeningHours(),
		PriceLevel:        la.GetPriceLevel(),
		Priority:          la.GetPriority(),
		DisplayOrder:      la.GetDisplayOrder(),
	}
}

// 数据库查询方法
func GetActiveLocalAdvertisements(city, region string, limit int) ([]*LocalAdvertisement, error) {
	var ads []*LocalAdvertisement
	db := DB

	query := db.Where("status = ?", protocol.StatusActive)

	if city != "" {
		query = query.Where("city = ?", city)
	}

	if region != "" {
		query = query.Where("region = ?", region)
	}

	// 添加时间有效性检查
	now := utils.TimeNowMilli()
	query = query.Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", now, now)

	err := query.Order("priority DESC, display_order ASC, created_at DESC").
		Limit(limit).
		Find(&ads).Error

	return ads, err
}

func GetLocalAdvertisementByID(adID string) (*LocalAdvertisement, error) {
	var ad LocalAdvertisement
	err := DB.Where("ad_id = ?", adID).First(&ad).Error
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

// 获取所有激活的本地广告（不限制数量，用于API）
func GetAllActiveLocalAdvertisements(city, region, category string) ([]*LocalAdvertisement, error) {
	var ads []*LocalAdvertisement
	db := DB

	query := db.Where("status = ?", protocol.StatusActive)

	if city != "" {
		query = query.Where("city = ?", city)
	}

	if region != "" {
		query = query.Where("region = ?", region)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 添加时间有效性检查
	now := utils.TimeNowMilli()
	query = query.Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", now, now)

	err := query.Order("priority DESC, display_order ASC, created_at DESC").
		Find(&ads).Error

	return ads, err
}

// 更新广告统计信息
func UpdateLocalAdvertisementStats(db *gorm.DB, adID string, statsType string) error {
	var updateField string
	switch statsType {
	case "view":
		updateField = "view_count = view_count + 1"
	case "click":
		updateField = "click_count = click_count + 1"
	case "call":
		updateField = "call_count = call_count + 1"
	default:
		return nil
	}

	return db.Model(&LocalAdvertisement{}).
		Where("ad_id = ?", adID).
		UpdateColumn(updateField, gorm.Expr(updateField)).Error
}

// UpdateLocalAdvertisementGoogleInfo 更新广告的Google信息
func UpdateLocalAdvertisementGoogleInfo(adID, googlePlaceID string, rating float64, reviewCount int, lat, lng float64, phone, website string) error {
	updates := map[string]interface{}{}

	if googlePlaceID != "" {
		updates["google_place_id"] = googlePlaceID
	}
	if rating > 0 {
		updates["google_rating"] = rating
	}
	if reviewCount > 0 {
		updates["google_review_count"] = reviewCount
	}
	if lat != 0 {
		updates["latitude"] = lat
	}
	if lng != 0 {
		updates["longitude"] = lng
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if website != "" {
		updates["website"] = website
	}

	if len(updates) == 0 {
		return nil // 没有要更新的字段
	}

	return DB.Model(&LocalAdvertisement{}).
		Where("ad_id = ?", adID).
		Updates(updates).Error
}

// UpdateLocalAdvertisementWithCompleteGoogleInfo 使用完整的Google信息更新广告
func UpdateLocalAdvertisementWithCompleteGoogleInfo(adID string, googleData map[string]interface{}) error {
	updates := map[string]interface{}{}

	// Google提供的真实数据
	if placeID, ok := googleData["place_id"].(string); ok && placeID != "" {
		updates["google_place_id"] = placeID
	}
	if name, ok := googleData["name"].(string); ok && name != "" {
		updates["name"] = name
	}
	if address, ok := googleData["formatted_address"].(string); ok && address != "" {
		updates["address"] = address
	}
	if phone, ok := googleData["international_phone_number"].(string); ok && phone != "" {
		updates["phone"] = phone
	}
	if website, ok := googleData["website"].(string); ok && website != "" {
		updates["website"] = website
	}
	if rating, ok := googleData["rating"].(float64); ok && rating > 0 {
		updates["google_rating"] = rating
	}
	if reviewCount, ok := googleData["user_ratings_total"].(int); ok && reviewCount > 0 {
		updates["google_review_count"] = reviewCount
	}
	if lat, ok := googleData["latitude"].(float64); ok && lat != 0 {
		updates["latitude"] = lat
	}
	if lng, ok := googleData["longitude"].(float64); ok && lng != 0 {
		updates["longitude"] = lng
	}
	if priceLevel, ok := googleData["price_level"].(int); ok && priceLevel >= 0 {
		updates["price_level"] = priceLevel
	}
	if types, ok := googleData["types"].([]string); ok && len(types) > 0 {
		// 将Google的types转换为tags字符串
		updates["tags"] = strings.Join(types, ",")
	}
	if openingHours, ok := googleData["opening_hours"].(string); ok && openingHours != "" {
		updates["opening_hours"] = openingHours
	}
	if imageURL, ok := googleData["image_url"].(string); ok && imageURL != "" {
		updates["image_url"] = imageURL
	}

	// 清空Google API不提供的假数据字段
	updates["email"] = nil       // Google Places API不提供邮箱
	updates["description"] = nil // Google Places API不提供具体描述
	updates["logo_url"] = nil    // Google Places API不提供Logo
	updates["image_urls"] = nil  // 暂时清空多图片，只保留第一张

	return DB.Model(&LocalAdvertisement{}).
		Where("ad_id = ?", adID).
		Updates(updates).Error
}

// GetLocalAdvertisementsByName 根据名称获取本地广告（用于Google信息更新）
func GetLocalAdvertisementsByName(name string) ([]*LocalAdvertisement, error) {
	var ads []*LocalAdvertisement
	err := DB.Where("name LIKE ?", "%"+name+"%").Find(&ads).Error
	return ads, err
}
