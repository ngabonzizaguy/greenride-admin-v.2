package protocol

// Currency constants
const (
	// CurrencyRWF 卢旺达法郎
	CurrencyRWF = "RWF"

	// DefaultCurrency 默认货币
	DefaultCurrency = CurrencyRWF
)

// GetDefaultCurrency 获取默认货币
func GetDefaultCurrency() string {
	return DefaultCurrency
}
