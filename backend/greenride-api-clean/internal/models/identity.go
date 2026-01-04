package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// Identity 身份认证表 - 基于最新设计文档
type Identity struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	IdentityID string `json:"identity_id" gorm:"column:identity_id;type:varchar(64);uniqueIndex"`
	Salt       string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*IdentityValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type IdentityValues struct {
	UserID   *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index;default:'user'"` // user, driver

	// 身份证件信息
	IDType    *string `json:"id_type" gorm:"column:id_type;type:varchar(32)"`       // passport, driver_license, national_id
	IDNumber  *string `json:"id_number" gorm:"column:id_number;type:varchar(100)"`  // 证件号码
	IDCountry *string `json:"id_country" gorm:"column:id_country;type:varchar(10)"` // 证件签发国家

	// 个人信息
	FirstName   *string `json:"first_name" gorm:"column:first_name;type:varchar(100)"`
	LastName    *string `json:"last_name" gorm:"column:last_name;type:varchar(100)"`
	MiddleName  *string `json:"middle_name" gorm:"column:middle_name;type:varchar(100)"`
	FullName    *string `json:"full_name" gorm:"column:full_name;type:varchar(255)"`
	DateOfBirth *string `json:"date_of_birth" gorm:"column:date_of_birth;type:varchar(32)"` // YYYY-MM-DD格式
	Gender      *string `json:"gender" gorm:"column:gender;type:varchar(16)"`               // male, female, other
	Nationality *string `json:"nationality" gorm:"column:nationality;type:varchar(64)"`

	// 地址信息
	Address    *string `json:"address" gorm:"column:address;type:varchar(500)"`
	City       *string `json:"city" gorm:"column:city;type:varchar(100)"`
	State      *string `json:"state" gorm:"column:state;type:varchar(100)"`
	Country    *string `json:"country" gorm:"column:country;type:varchar(100)"`
	PostalCode *string `json:"postal_code" gorm:"column:postal_code;type:varchar(32)"`

	// 证件图片
	FrontImageURL *string `json:"front_image_url" gorm:"column:front_image_url;type:varchar(512)"` // 证件正面
	BackImageURL  *string `json:"back_image_url" gorm:"column:back_image_url;type:varchar(512)"`   // 证件背面
	SelfieURL     *string `json:"selfie_url" gorm:"column:selfie_url;type:varchar(512)"`           // 自拍照

	// 验证状态
	Status           *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'pending'"` // pending, approved, rejected, expired
	VerificationStep *int    `json:"verification_step" gorm:"column:verification_step;default:1"`          // 验证步骤 1-5
	IsVerified       *bool   `json:"is_verified" gorm:"column:is_verified;default:false"`

	// 审核信息
	ReviewedBy    *string `json:"reviewed_by" gorm:"column:reviewed_by;type:varchar(64)"`      // 审核员ID
	ReviewedAt    *int64  `json:"reviewed_at" gorm:"column:reviewed_at"`                       // 审核时间
	ReviewComment *string `json:"review_comment" gorm:"column:review_comment;type:text"`       // 审核备注
	RejectReason  *string `json:"reject_reason" gorm:"column:reject_reason;type:varchar(255)"` // 拒绝原因

	// 有效期
	ExpiresAt *int64 `json:"expires_at" gorm:"column:expires_at"` // 证件过期时间

	// OCR识别结果
	OCRData       *string  `json:"ocr_data" gorm:"column:ocr_data;type:json"`                     // OCR识别的原始数据
	OCRConfidence *float64 `json:"ocr_confidence" gorm:"column:ocr_confidence;type:decimal(5,4)"` // OCR识别置信度

	// 扩展信息
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // JSON格式的额外信息

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Identity) TableName() string {
	return "t_identities"
}

// 身份验证状态常量
const (
	IdentityStatusPending  = "pending"
	IdentityStatusApproved = "approved"
	IdentityStatusRejected = "rejected"
	IdentityStatusExpired  = "expired"
)

// 证件类型常量
const (
	IDTypePassport      = "passport"
	IDTypeDriverLicense = "driver_license"
	IDTypeNationalID    = "national_id"
)

// 性别常量
const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

// 创建新的身份认证对象
func NewIdentityV2() *Identity {
	return &Identity{
		IdentityID: utils.GenerateIdentityID(),
		Salt:       utils.GenerateSalt(),
		IdentityValues: &IdentityValues{
			UserType:         utils.StringPtr(protocol.UserTypePassenger),
			Status:           utils.StringPtr(IdentityStatusPending),
			VerificationStep: utils.IntPtr(1),
			IsVerified:       utils.BoolPtr(false),
		},
	}
}

// SetValues 更新IdentityV2Values中的非nil值
func (i *IdentityValues) SetValues(values *IdentityValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		i.UserID = values.UserID
	}
	if values.UserType != nil {
		i.UserType = values.UserType
	}
	if values.IDType != nil {
		i.IDType = values.IDType
	}
	if values.IDNumber != nil {
		i.IDNumber = values.IDNumber
	}
	if values.FirstName != nil {
		i.FirstName = values.FirstName
	}
	if values.LastName != nil {
		i.LastName = values.LastName
	}
	if values.FullName != nil {
		i.FullName = values.FullName
	}
	if values.DateOfBirth != nil {
		i.DateOfBirth = values.DateOfBirth
	}
	if values.Gender != nil {
		i.Gender = values.Gender
	}
	if values.Status != nil {
		i.Status = values.Status
	}
	if values.IsVerified != nil {
		i.IsVerified = values.IsVerified
	}
	if values.FrontImageURL != nil {
		i.FrontImageURL = values.FrontImageURL
	}
	if values.BackImageURL != nil {
		i.BackImageURL = values.BackImageURL
	}
	if values.SelfieURL != nil {
		i.SelfieURL = values.SelfieURL
	}
	if values.Metadata != nil {
		i.Metadata = values.Metadata
	}
	if values.UpdatedAt > 0 {
		i.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (i *IdentityValues) GetUserID() string {
	if i.UserID == nil {
		return ""
	}
	return *i.UserID
}

func (i *IdentityValues) GetIDType() string {
	if i.IDType == nil {
		return ""
	}
	return *i.IDType
}

func (i *IdentityValues) GetIDNumber() string {
	if i.IDNumber == nil {
		return ""
	}
	return *i.IDNumber
}

func (i *IdentityValues) GetFirstName() string {
	if i.FirstName == nil {
		return ""
	}
	return *i.FirstName
}

func (i *IdentityValues) GetLastName() string {
	if i.LastName == nil {
		return ""
	}
	return *i.LastName
}

func (i *IdentityValues) GetFullName() string {
	if i.FullName == nil {
		return ""
	}
	return *i.FullName
}

func (i *IdentityValues) GetStatus() string {
	if i.Status == nil {
		return IdentityStatusPending
	}
	return *i.Status
}

func (i *IdentityValues) GetIsVerified() bool {
	if i.IsVerified == nil {
		return false
	}
	return *i.IsVerified
}

func (i *IdentityValues) GetVerificationStep() int {
	if i.VerificationStep == nil {
		return 1
	}
	return *i.VerificationStep
}

func (i *IdentityValues) GetGender() string {
	if i.Gender == nil {
		return ""
	}
	return *i.Gender
}

func (i *IdentityValues) GetDateOfBirth() string {
	if i.DateOfBirth == nil {
		return ""
	}
	return *i.DateOfBirth
}

func (i *IdentityValues) GetOCRConfidence() float64 {
	if i.OCRConfidence == nil {
		return 0.0
	}
	return *i.OCRConfidence
}

// Setter 方法
func (i *IdentityValues) SetUserID(userID string) *IdentityValues {
	i.UserID = &userID
	return i
}

func (i *IdentityValues) SetIDInfo(idType, idNumber, country string) *IdentityValues {
	i.IDType = &idType
	i.IDNumber = &idNumber
	i.IDCountry = &country
	return i
}

func (i *IdentityValues) SetPersonalInfo(firstName, lastName, fullName string) *IdentityValues {
	i.FirstName = &firstName
	i.LastName = &lastName
	i.FullName = &fullName
	return i
}

func (i *IdentityValues) SetStatus(status string) *IdentityValues {
	i.Status = &status
	return i
}

func (i *IdentityValues) SetVerificationStep(step int) *IdentityValues {
	i.VerificationStep = &step
	return i
}

func (i *IdentityValues) SetImageURLs(frontURL, backURL, selfieURL string) *IdentityValues {
	if frontURL != "" {
		i.FrontImageURL = &frontURL
	}
	if backURL != "" {
		i.BackImageURL = &backURL
	}
	if selfieURL != "" {
		i.SelfieURL = &selfieURL
	}
	return i
}

// 业务方法
func (i *Identity) IsPending() bool {
	return i.GetStatus() == IdentityStatusPending
}

func (i *Identity) IsApproved() bool {
	return i.GetStatus() == IdentityStatusApproved
}

func (i *Identity) IsRejected() bool {
	return i.GetStatus() == IdentityStatusRejected
}

func (i *Identity) IsExpired() bool {
	return i.GetStatus() == IdentityStatusExpired
}

func (i *Identity) IsVerified() bool {
	return i.GetIsVerified() && i.IsApproved()
}

func (i *Identity) CanSubmitForReview() bool {
	return i.IsPending() && i.HasRequiredImages()
}

func (i *Identity) HasRequiredImages() bool {
	return i.FrontImageURL != nil && *i.FrontImageURL != "" &&
		i.SelfieURL != nil && *i.SelfieURL != ""
}

func (i *Identity) IsDriverLicense() bool {
	return i.GetIDType() == IDTypeDriverLicense
}

func (i *Identity) IsPassport() bool {
	return i.GetIDType() == IDTypePassport
}

func (i *Identity) GetProgressPercentage() int {
	step := i.GetVerificationStep()
	if step >= 5 {
		return 100
	}
	return step * 20 // 每步20%
}

// 审核相关方法
func (i *IdentityValues) ApproveIdentity(reviewerID string) {
	i.SetStatus(IdentityStatusApproved)
	i.IsVerified = utils.BoolPtr(true)
	i.ReviewedBy = &reviewerID
	now := utils.TimeNowMilli()
	i.ReviewedAt = &now
	i.VerificationStep = utils.IntPtr(5) // 最终步骤
}

func (i *IdentityValues) RejectIdentity(reviewerID, reason string) {
	i.SetStatus(IdentityStatusRejected)
	i.IsVerified = utils.BoolPtr(false)
	i.ReviewedBy = &reviewerID
	i.RejectReason = &reason
	now := utils.TimeNowMilli()
	i.ReviewedAt = &now
}

func (i *IdentityValues) SetOCRResult(data string, confidence float64) *IdentityValues {
	i.OCRData = &data
	i.OCRConfidence = &confidence
	return i
}

func (i *IdentityValues) SetExpiryDate(expiresAt int64) *IdentityValues {
	i.ExpiresAt = &expiresAt
	return i
}

func (i *IdentityValues) SetAddress(address, city, state, country, postalCode string) *IdentityValues {
	if address != "" {
		i.Address = &address
	}
	if city != "" {
		i.City = &city
	}
	if state != "" {
		i.State = &state
	}
	if country != "" {
		i.Country = &country
	}
	if postalCode != "" {
		i.PostalCode = &postalCode
	}
	return i
}

func (i *IdentityValues) NextStep() *IdentityValues {
	currentStep := i.GetVerificationStep()
	if currentStep < 5 {
		i.SetVerificationStep(currentStep + 1)
	}
	return i
}

// 检查是否需要额外验证
func (i *Identity) RequiresAdditionalVerification() bool {
	// OCR置信度低于阈值
	if i.GetOCRConfidence() < 0.8 {
		return true
	}

	// 高风险国家或地区
	if i.IDCountry != nil {
		riskCountries := []string{"XX", "YY"} // 配置高风险国家列表
		for _, country := range riskCountries {
			if *i.IDCountry == country {
				return true
			}
		}
	}

	return false
}
