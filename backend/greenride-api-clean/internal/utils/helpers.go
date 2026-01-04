package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

var NUM_PATTERN = regexp.MustCompile(`^\d+.?(\d+)$`)

func SignMD5(sign_str, secret string, lower bool) string {
	sign_str = fmt.Sprintf("%v&key=%v", sign_str, secret)
	sign := Md5(sign_str)
	if lower {
		return strings.ToLower(sign)
	}
	return strings.ToUpper(sign)
}

func ToInt(in interface{}) int {
	return cast.ToInt(in)
}

func ToString(in interface{}) string {
	return cast.ToString(in)
}

func ToMap(in interface{}) map[string]interface{} {
	data, err := cast.ToStringMapE(in)
	if err == nil {
		return data
	}
	return map[string]interface{}{}
}

// 判断字符串是否为 纯数字 包涵浮点型
func IsNumber(s string) bool {
	return NUM_PATTERN.MatchString(s)
}

func ToJsonByte(obj interface{}) []byte {
	_data, _ := json.Marshal(obj)
	return _data
}

func ToJsonString(obj interface{}) string {
	_data, _ := json.Marshal(obj)
	return string(_data)
}

func ToQueryUrl(data map[string]interface{}) string {
	if len(data) == 0 {
		return ""
	}
	params := []string{}
	for k, v := range data {
		params = append(params, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(params, "&")
}

func GetBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Md5(signStr string) string {
	h := md5.New()
	h.Write([]byte(signStr))
	md5Str := hex.EncodeToString(h.Sum(nil))
	sign := strings.ToLower(md5Str)
	return sign
}

func GetHmacSha1(keyStr, value string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	// 进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return res
}

func GetSha256String(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

func GetSha256Base64(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes)
}

func GetSha512String(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

func GetHmacSha256Base64(str, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes)
}

func GetHmacSha256Hex(str, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}

func NewPwd(cert_no string, pwd string) string {
	md5str := fmt.Sprintf("%s|%s", cert_no, pwd)
	pwdStr := fmt.Sprintf("%x", md5.Sum([]byte(md5str)))
	return pwdStr
}

// GetEnv get key environment variable if exist otherwise return defalutValue
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func NewULID() string {
	now := time.Now()
	entropy := rand.New(rand.NewSource(now.UnixNano()))
	ms := ulid.Timestamp(now)
	id := ulid.MustNew(ms, entropy)
	return id.String()
}

// 计算以给定经纬度为中心，半径为radius公里的坐标范围
func CalculateBoundsForRadius(lat, lon, radius float64) (minLat, maxLat, minLon, maxLon float64) {
	// 纬度1度约等于111公里
	latDelta := radius / 111.0
	// 经度1度的距离随着纬度增加而减少
	lonDelta := radius / (111.0 * math.Cos(lat*(math.Pi/180.0)))

	minLat = lat - latDelta
	maxLat = lat + latDelta
	minLon = lon - lonDelta
	maxLon = lon + lonDelta

	return
}

// 工具函数
func BoolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func Normalize(value, min, max float64) float64 {
	if max <= min {
		return 0.5
	}
	return math.Max(0, math.Min(1, (value-min)/(max-min)))
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// EscapeCSV 转义CSV字段，处理包含逗号、引号和换行符的内容
func EscapeCSV(field string) string {
	// 如果字段包含逗号、引号或换行符，需要用引号包围并转义内部引号
	if strings.Contains(field, ",") || strings.Contains(field, "\"") || strings.Contains(field, "\n") || strings.Contains(field, "\r") {
		// 将字段中的引号转义为两个引号
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		// 用引号包围整个字段
		return "\"" + escaped + "\""
	}
	return field
}

// 指针辅助函数

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// IntPtr 返回int指针
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr 返回int64指针
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr 返回float64指针
func Float64Ptr(f float64) *float64 {
	return &f
}

// DecimalPtr 返回decimal.Decimal指针
func DecimalPtr(d decimal.Decimal) *decimal.Decimal {
	return &d
}

// NewDecimalFromFloat 创建一个decimal.Decimal从float64
func NewDecimalFromFloat(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// DecimalToFloat64 将decimal.Decimal转换为float64
func DecimalToFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

// IsDecimalZero 检查decimal.Decimal是否为零
func IsDecimalZero(d decimal.Decimal) bool {
	return d.Equal(decimal.Zero)
}

// IsDecimalGreaterThanZero 检查decimal.Decimal是否大于零
func IsDecimalGreaterThanZero(d decimal.Decimal) bool {
	return d.GreaterThan(decimal.Zero)
}

// RoundDecimal 将decimal.Decimal四舍五入到指定小数位
func RoundDecimal(d decimal.Decimal, places int32) decimal.Decimal {
	return d.Round(places)
}

// BoolPtr 返回bool指针
func BoolPtr(b bool) *bool {
	return &b
}

// JSON序列化辅助函数
func ToJSON(v any) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJSON(jsonStr string, v any) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// RoundFloat 浮点数四舍五入到指定小数位
func RoundFloat(f float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(f*multiplier) / multiplier
}

// SafeStringDeref 安全地解引用字符串指针
func SafeStringDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SafeFloat64Deref 安全地解引用float64指针
func SafeFloat64Deref(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// SafeIntDeref 安全地解引用int指针
func SafeIntDeref(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// SafeDecimalDeref 安全地解引用decimal.Decimal指针
func SafeDecimalDeref(d *decimal.Decimal) decimal.Decimal {
	if d == nil {
		return decimal.Zero
	}
	return *d
}

// RoundToTwoDecimal 将浮点数四舍五入到两位小数
func RoundToTwoDecimal(num float64) float64 {
	return math.Round(num*100) / 100
}
