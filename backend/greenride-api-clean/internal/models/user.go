package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID     int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(64);uniqueIndex"`
	Salt   string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*UserValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type UserValues struct {
	UserType    *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index;default:'user'"` // user, driver
	Password    *string `json:"-" gorm:"column:password;type:varchar(255)"`                              // 加密后的密码，不在JSON中显示
	Email       *string `json:"email" gorm:"column:email;type:varchar(255);index"`
	Phone       *string `json:"phone" gorm:"column:phone;type:varchar(64);index"`
	CountryCode *string `json:"country_code" gorm:"column:country_code;type:varchar(10)"`
	Username    *string `json:"username" gorm:"column:username;type:varchar(100);index"`
	DisplayName *string `json:"display_name" gorm:"column:display_name;type:varchar(100)"`
	FirstName   *string `json:"first_name" gorm:"column:first_name;type:varchar(100)"`
	LastName    *string `json:"last_name" gorm:"column:last_name;type:varchar(100)"`
	Avatar      *string `json:"avatar" gorm:"column:avatar;type:varchar(500)"`
	Gender      *string `json:"gender" gorm:"column:gender;type:varchar(10)"` // male, female, other
	Birthday    *int64  `json:"birthday" gorm:"column:birthday"`
	Language    *string `json:"language" gorm:"column:language;type:varchar(10);default:'en'"`
	Timezone    *string `json:"timezone" gorm:"column:timezone;type:varchar(50);default:'UTC'"`

	// 联系信息
	Address    *string `json:"address" gorm:"column:address;type:text"`
	City       *string `json:"city" gorm:"column:city;type:varchar(100)"`
	State      *string `json:"state" gorm:"column:state;type:varchar(100)"`
	Country    *string `json:"country" gorm:"column:country;type:varchar(100)"`
	PostalCode *string `json:"postal_code" gorm:"column:postal_code;type:varchar(20)"`

	// 位置信息
	Latitude          *float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8)"`
	Longitude         *float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8)"`
	LocationUpdatedAt *int64   `json:"location_updated_at" gorm:"column:location_updated_at"`

	// 状态信息
	Status          *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"` // active, inactive, suspended, banned
	IsEmailVerified *bool   `json:"is_email_verified" gorm:"column:is_email_verified;default:false"`
	IsPhoneVerified *bool   `json:"is_phone_verified" gorm:"column:is_phone_verified;default:false"`
	OnlineStatus    *string `json:"online_status" gorm:"column:online_status;type:varchar(20);default:'offline'"` // online, offline, busy

	// 司机相关信息
	LicenseNumber *string `json:"license_number" gorm:"column:license_number;type:varchar(50)"`
	LicenseExpiry *int64  `json:"license_expiry" gorm:"column:license_expiry"`

	// 司机队列管理字段（仅司机类型用户使用）
	QueuedOrderIds   []string `json:"queued_order_ids" gorm:"column:queued_order_ids;type:json;serializer:json"` // 队列订单ID列表
	CurrentOrderId   *string  `json:"current_order_id" gorm:"column:current_order_id;type:varchar(64)"`          // 当前执行订单ID
	NextAvailableAt  *int64   `json:"next_available_at" gorm:"column:next_available_at"`                         // 下次可用时间戳
	MaxQueueCapacity *int     `json:"max_queue_capacity" gorm:"column:max_queue_capacity;default:3"`             // 最大队列容量
	QueueUpdatedAt   *int64   `json:"queue_updated_at" gorm:"column:queue_updated_at"`                           // 队列最后更新时间

	// 通用评分和行程信息
	Score      *float64 `json:"score" gorm:"column:score;type:decimal(3,2);default:5.0"`
	TotalRides *int     `json:"total_rides" gorm:"column:total_rides;default:0"`

	// 推荐系统
	InviteCode  *string `json:"invite_code" gorm:"column:invite_code;type:varchar(20);uniqueIndex"`
	InvitedBy   *string `json:"invited_by" gorm:"column:invited_by;type:varchar(64);index"`
	InviteCount *int    `json:"invite_count" gorm:"column:invite_count;default:0"`

	// 设备和认证信息
	FCMToken   *string `json:"fcm_token" gorm:"column:fcm_token;type:varchar(500)"`
	DeviceID   *string `json:"device_id" gorm:"column:device_id;type:varchar(255)"`
	DeviceType *string `json:"device_type" gorm:"column:device_type;type:varchar(20)"` // ios, android, web
	AppVersion *string `json:"app_version" gorm:"column:app_version;type:varchar(20)"`

	// 时间戳
	LastLoginAt     *int64 `json:"last_login_at" gorm:"column:last_login_at"`
	EmailVerifiedAt *int64 `json:"email_verified_at" gorm:"column:email_verified_at"`
	PhoneVerifiedAt *int64 `json:"phone_verified_at" gorm:"column:phone_verified_at"`

	// 管理员操作记录
	Notes         *string `json:"notes" gorm:"column:notes;type:text"`                            // 管理员备注
	LastUpdatedBy *string `json:"last_updated_by" gorm:"column:last_updated_by;type:varchar(64)"` // 最后更新者ID

	// 扩展信息
	Metadata  *string `json:"metadata" gorm:"column:metadata;type:json"`               // JSON格式的额外信息
	Sandbox   *int    `json:"sandbox" gorm:"column:sandbox;type:tinyint(1);default:0"` // 是否为沙盒用户(0/1)
	UpdatedAt int64   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt *int64  `json:"deleted_at" gorm:"column:deleted_at"` // 删除时间戳，用于软删除
}

func (u *UserValues) SetUsername(username string) *UserValues {
	u.Username = &username
	return u
}
func (u *UserValues) GetUsername() string {
	if u.Username != nil {
		return *u.Username
	}
	return ""
}

func (u *UserValues) GetPassword() string {
	if u.Password != nil {
		return *u.Password
	}
	return ""
}

func (User) TableName() string {
	return "t_users"
}

// 创建新的用户对象
func NewUser() *User {
	userID := utils.GenerateUserID()
	defaultAvatar := protocol.GenerateDefaultAvatar(userID) // 生成默认头像

	return &User{
		UserID: userID,
		Salt:   utils.GenerateSalt(),
		UserValues: &UserValues{
			UserType:        utils.StringPtr(protocol.UserTypePassenger),
			Status:          utils.StringPtr(protocol.StatusActive),
			IsEmailVerified: utils.BoolPtr(false),
			IsPhoneVerified: utils.BoolPtr(false),
			OnlineStatus:    utils.StringPtr(protocol.StatusOffline),
			Language:        utils.StringPtr("en"),
			Timezone:        utils.StringPtr("UTC"),
			Score:           utils.Float64Ptr(5.0),
			TotalRides:      utils.IntPtr(0),
			InviteCount:     utils.IntPtr(0),
			InviteCode:      utils.StringPtr(utils.GenerateInviteCode()),
			Avatar:          utils.StringPtr(defaultAvatar), // 设置默认头像
			Sandbox:         utils.IntPtr(0),                // 默认设置为非沙盒用户
		},
	}
}

// SetValues 更新UserV2Values中的非nil值
func (u *UserValues) SetValues(values *UserValues) {
	if values == nil {
		return
	}

	if values.UserType != nil {
		u.UserType = values.UserType
	}
	if values.Password != nil {
		u.Password = values.Password
	}
	if values.Email != nil {
		u.Email = values.Email
	}
	if values.Phone != nil {
		u.Phone = values.Phone
	}
	if values.CountryCode != nil {
		u.CountryCode = values.CountryCode
	}
	if values.Username != nil {
		u.Username = values.Username
	}
	if values.DisplayName != nil {
		u.DisplayName = values.DisplayName
	}
	if values.FirstName != nil {
		u.FirstName = values.FirstName
	}
	if values.LastName != nil {
		u.LastName = values.LastName
	}
	if values.Avatar != nil {
		u.Avatar = values.Avatar
	}
	if values.Gender != nil {
		u.Gender = values.Gender
	}
	if values.Birthday != nil {
		u.Birthday = values.Birthday
	}
	if values.Language != nil {
		u.Language = values.Language
	}
	if values.Timezone != nil {
		u.Timezone = values.Timezone
	}
	if values.Address != nil {
		u.Address = values.Address
	}
	if values.City != nil {
		u.City = values.City
	}
	if values.State != nil {
		u.State = values.State
	}
	if values.Country != nil {
		u.Country = values.Country
	}
	if values.PostalCode != nil {
		u.PostalCode = values.PostalCode
	}
	if values.Latitude != nil {
		u.Latitude = values.Latitude
	}
	if values.Longitude != nil {
		u.Longitude = values.Longitude
	}
	if values.LocationUpdatedAt != nil {
		u.LocationUpdatedAt = values.LocationUpdatedAt
	}
	if values.Status != nil {
		u.Status = values.Status
	}
	if values.IsEmailVerified != nil {
		u.IsEmailVerified = values.IsEmailVerified
	}
	if values.IsPhoneVerified != nil {
		u.IsPhoneVerified = values.IsPhoneVerified
	}
	if values.OnlineStatus != nil {
		u.OnlineStatus = values.OnlineStatus
	}
	if values.LicenseNumber != nil {
		u.LicenseNumber = values.LicenseNumber
	}
	if values.LicenseExpiry != nil {
		u.LicenseExpiry = values.LicenseExpiry
	}
	// 队列管理相关字段
	if values.QueuedOrderIds != nil {
		u.QueuedOrderIds = values.QueuedOrderIds
	}
	if values.CurrentOrderId != nil {
		u.CurrentOrderId = values.CurrentOrderId
	}
	if values.NextAvailableAt != nil {
		u.NextAvailableAt = values.NextAvailableAt
	}
	if values.MaxQueueCapacity != nil {
		u.MaxQueueCapacity = values.MaxQueueCapacity
	}
	if values.QueueUpdatedAt != nil {
		u.QueueUpdatedAt = values.QueueUpdatedAt
	}
	if values.Score != nil {
		u.Score = values.Score
	}
	if values.TotalRides != nil {
		u.TotalRides = values.TotalRides
	}
	if values.InviteCode != nil {
		u.InviteCode = values.InviteCode
	}
	if values.InvitedBy != nil {
		u.InvitedBy = values.InvitedBy
	}
	if values.InviteCount != nil {
		u.InviteCount = values.InviteCount
	}
	if values.FCMToken != nil {
		u.FCMToken = values.FCMToken
	}
	if values.DeviceID != nil {
		u.DeviceID = values.DeviceID
	}
	if values.DeviceType != nil {
		u.DeviceType = values.DeviceType
	}
	if values.AppVersion != nil {
		u.AppVersion = values.AppVersion
	}
	// 时间戳相关字段
	if values.LastLoginAt != nil {
		u.LastLoginAt = values.LastLoginAt
	}
	if values.EmailVerifiedAt != nil {
		u.EmailVerifiedAt = values.EmailVerifiedAt
	}
	if values.PhoneVerifiedAt != nil {
		u.PhoneVerifiedAt = values.PhoneVerifiedAt
	}
	// 管理员相关字段
	if values.Notes != nil {
		u.Notes = values.Notes
	}
	if values.LastUpdatedBy != nil {
		u.LastUpdatedBy = values.LastUpdatedBy
	}
	if values.Metadata != nil {
		u.Metadata = values.Metadata
	}
	if values.Sandbox != nil {
		u.Sandbox = values.Sandbox
	}
	if values.UpdatedAt > 0 {
		u.UpdatedAt = values.UpdatedAt
	}
	if values.DeletedAt != nil {
		u.DeletedAt = values.DeletedAt
	}
}

func (u *UserValues) IsPassenger() bool {
	return u.GetUserType() == protocol.UserTypePassenger
}

func (u *UserValues) IsDriver() bool {
	return u.GetUserType() == protocol.UserTypeDriver
}

// Getter 方法
func (u *UserValues) GetUserType() string {
	if u.UserType == nil {
		return ""
	}
	return *u.UserType
}

func (u *UserValues) GetEmail() string {
	if u.Email == nil {
		return ""
	}
	return *u.Email
}

func (u *UserValues) GetPhone() string {
	if u.Phone == nil {
		return ""
	}
	return *u.Phone
}

func (u *UserValues) GetDisplayName() string {
	if u.DisplayName == nil {
		return ""
	}
	return *u.DisplayName
}

func (u *UserValues) GetFirstName() string {
	if u.FirstName == nil {
		return ""
	}
	return *u.FirstName
}

func (u *UserValues) GetLastName() string {
	if u.LastName == nil {
		return ""
	}
	return *u.LastName
}

func (u *UserValues) GetFullName() string {
	firstName := ""
	lastName := ""
	if u.FirstName != nil {
		firstName = *u.FirstName
	}
	if u.LastName != nil {
		lastName = *u.LastName
	}
	if firstName == "" && lastName == "" {
		return u.GetDisplayName()
	}
	return firstName + " " + lastName
}

func (u *UserValues) GetStatus() string {
	if u.Status == nil {
		return protocol.StatusActive
	}
	return *u.Status
}

func (u *UserValues) IsActive() bool {
	return u.GetStatus() == protocol.StatusActive
}

func (u *UserValues) GetIsEmailVerified() bool {
	if u.IsEmailVerified == nil {
		return false
	}
	return *u.IsEmailVerified
}

func (u *UserValues) GetIsPhoneVerified() bool {
	if u.IsPhoneVerified == nil {
		return false
	}
	return *u.IsPhoneVerified
}

func (u *UserValues) GetOnlineStatus() string {
	if u.OnlineStatus == nil {
		return protocol.StatusOffline
	}
	return *u.OnlineStatus
}

func (u *UserValues) GetScore() float64 {
	if u.Score == nil {
		return 5.0
	}
	return *u.Score
}

func (u *UserValues) GetRating() float64 {
	// Score 即为评分
	return u.GetScore()
}

func (u *UserValues) GetMaxQueueCapacity() int {
	if u.MaxQueueCapacity == nil {
		return 3 // 默认3
	}
	return *u.MaxQueueCapacity
}

func (u *UserValues) GetVehicleID() string {
	// 这里需要查询关联的车辆表，暂时返回空字符串
	// 在实际实现中应该查询vehicle表获取绑定的车辆ID
	return ""
}

func (u *UserValues) GetCurrentOrderID() string {
	if u.CurrentOrderId == nil {
		return ""
	}
	return *u.CurrentOrderId
}

func (u *UserValues) GetInviteCode() string {
	if u.InviteCode == nil {
		return ""
	}
	return *u.InviteCode
}

// Setter 方法
func (u *UserValues) SetUserType(userType string) *UserValues {
	u.UserType = &userType
	return u
}

func (u *UserValues) SetEmail(email string) *UserValues {
	u.Email = &email
	return u
}

func (u *UserValues) SetPhone(phone string) *UserValues {
	u.Phone = &phone
	return u
}

func (u *UserValues) SetFirstName(firstName string) *UserValues {
	u.FirstName = &firstName
	return u
}

func (u *UserValues) SetLastName(lastName string) *UserValues {
	u.LastName = &lastName
	return u
}

func (u *UserValues) SetDisplayName(name string) *UserValues {
	u.DisplayName = &name
	return u
}

func (u *UserValues) SetStatus(status string) *UserValues {
	u.Status = &status
	return u
}

// SetActiveStatus 设置用户在线状态
func (u *UserValues) SetActiveStatus(status string) *UserValues {
	u.OnlineStatus = &status
	return u
}

// SetOnlineStatus 设置用户在线状态（语义化方法）
func (u *UserValues) SetOnlineStatus(status string) *UserValues {
	return u.SetActiveStatus(status)
}

func (u *UserValues) SetLocation(lat, lng float64) *UserValues {
	u.Latitude = &lat
	u.Longitude = &lng
	now := utils.TimeNowMilli()
	u.LocationUpdatedAt = &now
	return u
}

func (u *UserValues) SetPassword(password string) *UserValues {
	u.Password = &password
	return u
}

func (u *UserValues) SetCountryCode(countryCode string) *UserValues {
	u.CountryCode = &countryCode
	return u
}

func (u *UserValues) SetAvatar(avatar string) *UserValues {
	u.Avatar = &avatar
	return u
}

func (u *UserValues) SetGender(gender string) *UserValues {
	u.Gender = &gender
	return u
}

func (u *UserValues) SetBirthday(birthday int64) *UserValues {
	u.Birthday = &birthday
	return u
}

func (u *UserValues) SetLanguage(language string) *UserValues {
	u.Language = &language
	return u
}

func (u *UserValues) SetTimezone(timezone string) *UserValues {
	u.Timezone = &timezone
	return u
}

func (u *UserValues) SetAddress(address string) *UserValues {
	u.Address = &address
	return u
}

func (u *UserValues) SetCity(city string) *UserValues {
	u.City = &city
	return u
}

func (u *UserValues) SetState(state string) *UserValues {
	u.State = &state
	return u
}

func (u *UserValues) SetCountry(country string) *UserValues {
	u.Country = &country
	return u
}

func (u *UserValues) SetPostalCode(postalCode string) *UserValues {
	u.PostalCode = &postalCode
	return u
}

func (u *UserValues) SetLatitude(latitude float64) *UserValues {
	u.Latitude = &latitude
	return u
}

func (u *UserValues) SetLongitude(longitude float64) *UserValues {
	u.Longitude = &longitude
	return u
}

func (u *UserValues) SetLocationUpdatedAt(timestamp int64) *UserValues {
	u.LocationUpdatedAt = &timestamp
	return u
}

func (u *UserValues) SetIsEmailVerified(isVerified bool) *UserValues {
	u.IsEmailVerified = &isVerified
	return u
}

func (u *UserValues) SetIsPhoneVerified(isVerified bool) *UserValues {
	u.IsPhoneVerified = &isVerified
	return u
}

func (u *UserValues) SetLicenseNumber(licenseNumber string) *UserValues {
	u.LicenseNumber = &licenseNumber
	return u
}

func (u *UserValues) SetLicenseExpiry(expiry int64) *UserValues {
	u.LicenseExpiry = &expiry
	return u
}

func (u *UserValues) SetScore(score float64) *UserValues {
	u.Score = &score
	return u
}

func (u *UserValues) SetTotalRides(count int) *UserValues {
	u.TotalRides = &count
	return u
}

func (u *UserValues) SetInviteCode(inviteCode string) *UserValues {
	u.InviteCode = &inviteCode
	return u
}

func (u *UserValues) SetInvitedBy(invitedBy string) *UserValues {
	u.InvitedBy = &invitedBy
	return u
}

func (u *UserValues) SetInviteCount(count int) *UserValues {
	u.InviteCount = &count
	return u
}

func (u *UserValues) SetFCMToken(token string) *UserValues {
	u.FCMToken = &token
	return u
}

func (u *UserValues) SetDeviceID(deviceID string) *UserValues {
	u.DeviceID = &deviceID
	return u
}

func (u *UserValues) SetDeviceType(deviceType string) *UserValues {
	u.DeviceType = &deviceType
	return u
}

func (u *UserValues) SetAppVersion(version string) *UserValues {
	u.AppVersion = &version
	return u
}

func (u *UserValues) SetLastLoginAt(timestamp int64) *UserValues {
	u.LastLoginAt = &timestamp
	return u
}

func (u *UserValues) SetEmailVerifiedAt(timestamp int64) *UserValues {
	u.EmailVerifiedAt = &timestamp
	return u
}

func (u *UserValues) SetPhoneVerifiedAt(timestamp int64) *UserValues {
	u.PhoneVerifiedAt = &timestamp
	return u
}

func (u *UserValues) SetNotes(notes string) *UserValues {
	u.Notes = &notes
	return u
}

func (u *UserValues) SetLastUpdatedBy(userID string) *UserValues {
	u.LastUpdatedBy = &userID
	return u
}

func (u *UserValues) SetMetadata(metadata string) *UserValues {
	u.Metadata = &metadata
	return u
}

// 业务方法
func (u *User) IsDriver() bool {
	return u.GetUserType() == protocol.UserTypeDriver
}

func (u *User) IsUser() bool {
	return u.GetUserType() == protocol.UserTypePassenger
}

func (u *User) CanTakeRides() bool {
	return u.IsActive()
}

func (u *User) CanDriveRides() bool {
	return u.IsDriver() && u.CanTakeRides() && u.GetIsPhoneVerified()
}

func (u *UserValues) MarkEmailAsVerified() {
	verified := true
	u.IsEmailVerified = &verified
	now := utils.TimeNowMilli()
	u.EmailVerifiedAt = &now
}

func (u *UserValues) MarkPhoneAsVerified() {
	verified := true
	u.IsPhoneVerified = &verified
	now := utils.TimeNowMilli()
	u.PhoneVerifiedAt = &now
}

func (u *UserValues) UpdateLastLogin() {
	now := utils.TimeNowMilli()
	u.LastLoginAt = &now
}

func (u *UserValues) IncrementRideCount() {
	count := u.GetTotalRides() + 1
	u.TotalRides = &count
}

func (u *UserValues) GetTotalRides() int {
	if u.TotalRides == nil {
		return 0
	}
	return *u.TotalRides
}

func (u *UserValues) GetLanguage() string {
	if u.Language == nil {
		return "en"
	}
	return *u.Language
}

func (u *UserValues) GetTimezone() string {
	if u.Timezone == nil {
		return "UTC"
	}
	return *u.Timezone
}

func (u *UserValues) GetCountryCode() string {
	if u.CountryCode == nil {
		return ""
	}
	return *u.CountryCode
}

func (u *UserValues) GetAvatar() string {
	if u.Avatar == nil {
		return ""
	}
	return *u.Avatar
}

func (u *UserValues) GetGender() string {
	if u.Gender == nil {
		return ""
	}
	return *u.Gender
}

func (u *UserValues) GetBirthday() int64 {
	if u.Birthday == nil {
		return 0
	}
	return *u.Birthday
}

func (u *UserValues) GetAddress() string {
	if u.Address == nil {
		return ""
	}
	return *u.Address
}

func (u *UserValues) GetCity() string {
	if u.City == nil {
		return ""
	}
	return *u.City
}

func (u *UserValues) GetState() string {
	if u.State == nil {
		return ""
	}
	return *u.State
}

func (u *UserValues) GetCountry() string {
	if u.Country == nil {
		return ""
	}
	return *u.Country
}

func (u *UserValues) GetPostalCode() string {
	if u.PostalCode == nil {
		return ""
	}
	return *u.PostalCode
}

func (u *UserValues) GetLatitude() float64 {
	if u.Latitude == nil {
		return 0.0
	}
	return *u.Latitude
}

func (u *UserValues) GetLongitude() float64 {
	if u.Longitude == nil {
		return 0.0
	}
	return *u.Longitude
}

func (u *UserValues) GetLocationUpdatedAt() int64 {
	if u.LocationUpdatedAt == nil {
		return 0
	}
	return *u.LocationUpdatedAt
}

func (u *UserValues) GetLicenseNumber() string {
	if u.LicenseNumber == nil {
		return ""
	}
	return *u.LicenseNumber
}

func (u *UserValues) GetLicenseExpiry() int64 {
	if u.LicenseExpiry == nil {
		return 0
	}
	return *u.LicenseExpiry
}

func (u *UserValues) GetInvitedBy() string {
	if u.InvitedBy == nil {
		return ""
	}
	return *u.InvitedBy
}

func (u *UserValues) GetInviteCount() int {
	if u.InviteCount == nil {
		return 0
	}
	return *u.InviteCount
}

func (u *UserValues) GetFCMToken() string {
	if u.FCMToken == nil {
		return ""
	}
	return *u.FCMToken
}

func (u *UserValues) GetDeviceID() string {
	if u.DeviceID == nil {
		return ""
	}
	return *u.DeviceID
}

func (u *UserValues) GetDeviceType() string {
	if u.DeviceType == nil {
		return ""
	}
	return *u.DeviceType
}

func (u *UserValues) GetAppVersion() string {
	if u.AppVersion == nil {
		return ""
	}
	return *u.AppVersion
}

func (u *UserValues) GetLastLoginAt() int64 {
	if u.LastLoginAt == nil {
		return 0
	}
	return *u.LastLoginAt
}

func (u *UserValues) GetEmailVerifiedAt() int64 {
	if u.EmailVerifiedAt == nil {
		return 0
	}
	return *u.EmailVerifiedAt
}

func (u *UserValues) GetPhoneVerifiedAt() int64 {
	if u.PhoneVerifiedAt == nil {
		return 0
	}
	return *u.PhoneVerifiedAt
}

func (u *UserValues) GetNotes() string {
	if u.Notes == nil {
		return ""
	}
	return *u.Notes
}

func (u *UserValues) GetLastUpdatedBy() string {
	if u.LastUpdatedBy == nil {
		return ""
	}
	return *u.LastUpdatedBy
}

func (u *UserValues) GetMetadata() string {
	if u.Metadata == nil {
		return ""
	}
	return *u.Metadata
}

func (u *UserValues) GetSandbox() int {
	if u.Sandbox == nil {
		return 0
	}
	return *u.Sandbox
}

func (u *UserValues) IsSandbox() bool {
	return u.GetSandbox() > 0
}

func (u *UserValues) SetSandbox(sandbox int) *UserValues {
	u.Sandbox = &sandbox
	return u
}

// DeletedAt getter和setter
func (u *UserValues) GetDeletedAt() int64 {
	if u.DeletedAt == nil {
		return 0
	}
	return *u.DeletedAt
}

func (u *UserValues) SetDeletedAt(timestamp int64) *UserValues {
	u.DeletedAt = &timestamp
	return u
}

// IsDeleted 检查用户是否已被删除
func (u *UserValues) IsDeleted() bool {
	return u.DeletedAt != nil && *u.DeletedAt > 0
}

// MarkAsDeleted 标记用户为已删除
func (u *UserValues) MarkAsDeleted() *UserValues {
	now := utils.TimeNowMilli()
	u.DeletedAt = &now
	return u
}

// RestoreDeleted 恢复已删除的用户
func (u *UserValues) RestoreDeleted() *UserValues {
	u.DeletedAt = nil
	return u
}

// ================== 缺少的 Getter/Setter 方法 ==================

// CurrentOrderId setter
func (u *UserValues) SetCurrentOrderId(orderID string) *UserValues {
	u.CurrentOrderId = &orderID
	return u
}

// NextAvailableAt getter和setter
func (u *UserValues) GetNextAvailableAt() int64 {
	if u.NextAvailableAt == nil {
		return 0
	}
	return *u.NextAvailableAt
}

func (u *UserValues) SetNextAvailableAt(timestamp int64) *UserValues {
	u.NextAvailableAt = &timestamp
	return u
}

// MaxQueueCapacity setter
func (u *UserValues) SetMaxQueueCapacity(capacity int) *UserValues {
	u.MaxQueueCapacity = &capacity
	return u
}

// QueueUpdatedAt getter和setter
func (u *UserValues) GetQueueUpdatedAt() int64 {
	if u.QueueUpdatedAt == nil {
		return 0
	}
	return *u.QueueUpdatedAt
}

func (u *UserValues) SetQueueUpdatedAt(timestamp int64) *UserValues {
	u.QueueUpdatedAt = &timestamp
	return u
}

// QueuedOrderIds getter和setter
func (u *UserValues) GetQueuedOrderIds() []string {
	if u.QueuedOrderIds == nil {
		return []string{}
	}
	return u.QueuedOrderIds
}

func (u *UserValues) SetQueuedOrderIds(orderIds []string) *UserValues {
	u.QueuedOrderIds = orderIds
	return u
}

func GetUserByID(userID string) *User {
	var user User
	err := DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

// GetUserByInviteCode 根据邀请码获取用户
func GetUserByInviteCode(inviteCode string) *User {
	var user User
	err := DB.Where("invite_code = ?", inviteCode).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

// IncrementUserInviteCount 增加用户邀请计数
func IncrementUserInviteCount(userID string) error {
	return DB.Model(&User{}).
		Where("user_id = ?", userID).
		UpdateColumn("invite_count", gorm.Expr("invite_count + ?", 1)).
		Error
}
func GetUserByEmailAndType(email, userType string) *User {
	var user User
	err := DB.Where("email = ? AND user_type = ?", email, userType).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

// GetUserByPhoneAndType 根据手机号和用户类型获取用户
func GetUserByPhoneAndType(phone, userType string) *User {
	var user User
	err := DB.Where("phone = ? AND user_type = ?", phone, userType).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

// NewUser 将 UserV2 模型转换为 User 协议结构
func (user *User) Protocol() *protocol.User {
	info := &protocol.User{
		UserID:            user.UserID,
		UserType:          user.GetUserType(),
		Email:             user.GetEmail(),
		Phone:             user.GetPhone(),
		FullName:          user.GetFullName(),
		Language:          user.GetLanguage(),
		Timezone:          user.GetTimezone(),
		Status:            user.GetStatus(),
		IsEmailVerified:   user.GetIsEmailVerified(),
		IsPhoneVerified:   user.GetIsPhoneVerified(),
		IsActive:          user.IsActive(),
		OnlineStatus:      user.GetOnlineStatus(),
		CountryCode:       user.GetCountryCode(),
		Username:          user.GetUsername(),
		DisplayName:       user.GetDisplayName(),
		FirstName:         user.GetFirstName(),
		LastName:          user.GetLastName(),
		Avatar:            user.GetAvatar(),
		Gender:            user.GetGender(),
		Birthday:          user.GetBirthday(),
		City:              user.GetCity(),
		State:             user.GetState(),
		Country:           user.GetCountry(),
		PostalCode:        user.GetPostalCode(),
		Latitude:          user.GetLatitude(),
		Longitude:         user.GetLongitude(),
		LocationUpdatedAt: user.GetLocationUpdatedAt(),
		Score:             user.GetScore(),
		TotalRides:        user.GetTotalRides(),
		LicenseNumber:     user.GetLicenseNumber(),
		LicenseExpiry:     user.GetLicenseExpiry(),
		InviteCode:        user.GetInviteCode(),
		InviteCount:       user.GetInviteCount(),
		Sandbox:           user.GetSandbox(),
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		DeletedAt:         user.GetDeletedAt(),
	}
	return info
}
