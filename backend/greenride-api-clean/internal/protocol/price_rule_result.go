package protocol

// PriceRuleResult 价格规则计算结果，包含基本规则信息和计算结果
type PriceRuleResult struct {
	// 基本规则信息
	RuleID      string `json:"rule_id,omitempty"`
	RuleName    string `json:"rule_name,omitempty"`
	Category    string `json:"category,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`

	// 计算结果字段
	Amount  float64 `json:"amount,omitempty"`  // 计算金额
	Applied bool    `json:"applied,omitempty"` // 是否应用
	Ccy     string  `json:"ccy,omitempty"`     // 货币
	Reason  string  `json:"reason,omitempty"`  // 原因
}
