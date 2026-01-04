package protocol

// 优惠券折扣类型常量
const (
	PromoDiscountTypePercentage  = "percentage"   // 百分比折扣
	PromoDiscountTypeFixedAmount = "fixed_amount" // 固定金额折扣
)

// 优惠券来源类型常量
const (
	UserPromotionSourceSystem   = "system"   // 系统发放
	UserPromotionSourceAdmin    = "admin"    // 管理员发放
	UserPromotionSourceEvent    = "event"    // 事件触发
	UserPromotionSourceReferral = "referral" // 推荐获得
	UserPromotionSourceWelcome  = "welcome"  // 欢迎奖励
	UserPromotionSourceReward   = "reward"   // 奖励发放
)
