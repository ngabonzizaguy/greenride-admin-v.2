package protocol

// 价格规则类别常量
const (
	PriceRuleCategoryBasePricing   = "base_pricing"
	PriceRuleCategorySurgePricing  = "surge_pricing"
	PriceRuleCategoryDiscount      = "discount"
	PriceRuleCategoryPromotion     = "promotion"
	PriceRuleCategoryUserPromotion = "user_promotion"
	PriceRuleCategorySpecialOffer  = "special_offer"
	PriceRuleCategoryDistanceFare  = "distance_fare"
	PriceRuleCategoryTimeFare      = "time_fare"
	PriceRuleCategoryServiceFee    = "service_fee"
)

// 价格规则类型常量
const (
	PriceRuleTypePercentage  = "percentage"
	PriceRuleTypeFixedAmount = "fixed_amount"
	PriceRuleTypeMultiplier  = "multiplier"
	PriceRuleTypeTiered      = "tiered"
	PriceRuleTypeCustom      = "custom"
)

// 定价模型常量
const (
	PricingModelDistanceBased = "distance_based"
	PricingModelTimeBased     = "time_based"
	PricingModelFixedRate     = "fixed_rate"
	PricingModelDynamic       = "dynamic"
)

// 折扣类型常量
const (
	DiscountTypePercentage   = "percentage"
	DiscountTypeFixed        = "fixed"
	DiscountTypeBuyXGetY     = "buy_x_get_y"
	DiscountTypeFreeDelivery = "free_delivery"
)

// 变更类型常量
const (
	ChangeTypeMajor  = "major"
	ChangeTypeMinor  = "minor"
	ChangeTypePatch  = "patch"
	ChangeTypeHotfix = "hotfix"
)
