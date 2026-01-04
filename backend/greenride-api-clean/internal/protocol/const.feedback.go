package protocol

// 反馈类型常量
const (
	FeedbackTypeComplaint      = "complaint"
	FeedbackTypeSuggestion     = "suggestion"
	FeedbackTypeCompliment     = "compliment"
	FeedbackTypeBugReport      = "bug_report"
	FeedbackTypeFeatureRequest = "feature_request"
	FeedbackTypeSafetyIssue    = "safety_issue"
)

// 反馈分类常量
const (
	FeedbackCategoryService = "service"
	FeedbackCategoryDriver  = "driver"
	FeedbackCategoryVehicle = "vehicle"
	FeedbackCategoryApp     = "app"
	FeedbackCategoryPayment = "payment"
	FeedbackCategorySafety  = "safety"
	FeedbackCategoryOther   = "other"
)

// 严重程度常量
const (
	FeedbackSeverityLow      = "low"
	FeedbackSeverityMedium   = "medium"
	FeedbackSeverityHigh     = "high"
	FeedbackSeverityCritical = "critical"
)

// 优先级常量
const (
	FeedbackPriorityLow    = "low"
	FeedbackPriorityMedium = "medium"
	FeedbackPriorityHigh   = "high"
	FeedbackPriorityUrgent = "urgent"
)

// SLA级别常量
const (
	FeedbackSLAStandard  = "standard"
	FeedbackSLAPriority  = "priority"
	FeedbackSLAEmergency = "emergency"
)

// 情感类型常量
const (
	EmotionTypeAngry      = "angry"
	EmotionTypeFrustrated = "frustrated"
	EmotionTypeSatisfied  = "satisfied"
	EmotionTypeHappy      = "happy"
	EmotionTypeNeutral    = "neutral"
)

// 特定影响级别常量
const (
	FeedbackImpactMinor    = "minor"
	FeedbackImpactModerate = "moderate"
	FeedbackImpactMajor    = "major"
	FeedbackImpactSevere   = "severe"
)

// 反馈相关的评分等级常量
const (
	RatingExcellent = 5
	RatingGood      = 4
	RatingAverage   = 3
	RatingPoor      = 2
	RatingTerrible  = 1
)

// 反馈分析语气常量
const (
	ToneAnalysisFormal     = "formal"
	ToneAnalysisCasual     = "casual"
	ToneAnalysisAggressive = "aggressive"
	ToneAnalysisPolite     = "polite"
)
