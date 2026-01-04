package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	SALT_KEY = "jK9$mP2#nL5@qR8*"

	// ID前缀常量 - 参考stay-hub-api设计
	ID_PREFIX_USER                = "U"
	ID_PREFIX_DRIVER              = "D"
	ID_PREFIX_ORDER               = "O"
	ID_PREFIX_RIDE_ORDER          = "R"
	ID_PREFIX_VEHICLE             = "V"
	ID_PREFIX_PAYMENT             = "P"
	ID_PREFIX_WALLET              = "W"
	ID_PREFIX_WALLET_TX           = "WT"
	ID_PREFIX_IDENTITY            = "ID"
	ID_PREFIX_NOTIFICATION        = "N"
	ID_PREFIX_REVIEW              = "RV"
	ID_PREFIX_CARD                = "C"
	ID_PREFIX_PROMOTION           = "PM"
	ID_PREFIX_USER_PROMOTION      = "UPM"
	ID_PREFIX_ANALYTICS           = "A"
	ID_PREFIX_SERVICE_AREA        = "SA"
	ID_PREFIX_USER_ACCOUNT        = "UA"
	ID_PREFIX_WITHDRAWAL          = "WD"
	ID_PREFIX_ADMIN               = "AD"
	ID_PREFIX_PAYMENT_METHOD      = "PM"
	ID_PREFIX_VEHICLE_TYPE        = "VT"
	ID_PREFIX_USER_ADDRESS        = "ADDR"
	ID_PREFIX_PRICE_RULE          = "PR"
	ID_PREFIX_USER_PAYMENT_METHOD = "UPM"
	ID_PREFIX_FEEDBACK            = "FB"
	ID_PREFIX_FCM_TOKEN           = "FT"
	ID_PREFIX_ANNOUNCEMENT        = "ANN"
	ID_PREFIX_PRICE_SNAPSHOT      = "PS"
	ID_PREFIX_RULE_VERSION        = "RV"
	ID_PREFIX_CHECKOUT            = "CO"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return ulid.Make().String()
}

func GenerateID() string {
	return ulid.Make().String()
}

// GenerateShortID 生成短ID
func GenerateShortID() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixMilli(), 36)
	}
	return hex.EncodeToString(b)
}

func GenerateCheckoutID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_CHECKOUT, GenerateID())
}

// 用户相关ID生成
func GenerateUserID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_USER, GenerateID())
}

func GenerateDriverID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_DRIVER, GenerateID())
}

// 订单相关ID生成
func GenerateOrderID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_ORDER, GenerateID())
}

func GenerateRideOrderID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_RIDE_ORDER, GenerateID())
}

// GenerateOrderIDByType 根据订单类型生成对应前缀的订单ID
func GenerateOrderIDByType(orderType string) string {
	switch orderType {
	case "ride":
		return GenerateRideOrderID()
	case "delivery":
		// 如果将来有外卖订单，可以添加专门的前缀
		return fmt.Sprintf("D%v", GenerateID())
	case "shopping":
		// 如果将来有购物订单，可以添加专门的前缀
		return fmt.Sprintf("S%v", GenerateID())
	default:
		// 默认使用通用订单前缀
		return GenerateOrderID()
	}
}

// 车辆相关ID生成
func GenerateVehicleID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_VEHICLE, GenerateID())
}

// 支付相关ID生成
func GeneratePaymentID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_PAYMENT, GenerateID())
}

// GenerateSandboxChannelPaymentID 生成沙盒渠道支付ID
func GenerateSandboxChannelPaymentID() string {
	return fmt.Sprintf("sandbox_%v", GenerateID())
}

// GenerateCashChannelPaymentID 生成现金渠道支付ID
func GenerateCashChannelPaymentID() string {
	return fmt.Sprintf("cash_%v", GenerateID())
}

func GenerateWalletID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_WALLET, GenerateID())
}

func GenerateWalletTransactionID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_WALLET_TX, GenerateID())
}

func GenerateUserAccountID() string {
	return fmt.Sprintf("UA%v", GenerateID())
}

// GenerateWithdrawalID 生成提现记录ID
func GenerateWithdrawalID() string {
	return fmt.Sprintf("WD%v", GenerateID())
}

// GenerateAdminID 生成管理员ID
func GenerateAdminID() string {
	return fmt.Sprintf("ADM%v", GenerateID())
}

func GenerateDispatchID() string {
	return fmt.Sprintf("DIS%v", GenerateID())
}

// GeneratePaymentMethodID 生成支付方式ID
func GeneratePaymentMethodID() string {
	return fmt.Sprintf("PM%v", GenerateID())
}

// GenerateVehicleTypeID 生成车型类型ID
func GenerateVehicleTypeID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_VEHICLE_TYPE, GenerateID())
}

func GenerateCardID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_CARD, GenerateID())
}

// 身份认证相关ID生成
func GenerateIdentityID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_IDENTITY, GenerateID())
}

// 通知相关ID生成
func GenerateNotificationID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_NOTIFICATION, GenerateID())
}

// 评价相关ID生成
func GenerateReviewID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_REVIEW, GenerateID())
}

// 促销代码相关ID生成
func GeneratePromotionID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_PROMOTION, GenerateID())
}

// 用户优惠券相关ID生成
func GenerateUserPromotionID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_USER_PROMOTION, GenerateID())
}

// 分析统计相关ID生成
func GenerateAnalyticsID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_ANALYTICS, GenerateID())
}

// 服务区域相关ID生成
func GenerateServiceAreaID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_SERVICE_AREA, GenerateID())
}

// 用户地址相关ID生成
func GenerateUserAddressID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_USER_ADDRESS, GenerateID())
}

// 价格规则相关ID生成
func GeneratePriceRuleID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_PRICE_RULE, GenerateID())
}

// FCM消息相关ID生成
func GenerateFCMMessageID() string {
	return fmt.Sprintf("FCM%v", GenerateID())
}

// 用户支付方式关联相关ID生成
func GenerateUserPaymentMethodID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_USER_PAYMENT_METHOD, GenerateID())
}

// 反馈相关ID生成
// GenerateFeedbackID 生成反馈ID
func GenerateFeedbackID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_FEEDBACK, GenerateID())
}

// GenerateFCMTokenID 生成FCM令牌ID
func GenerateFCMTokenID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_FCM_TOKEN, GenerateID())
}

// GenerateAnnouncementID 生成公告ID
func GenerateAnnouncementID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_ANNOUNCEMENT, GenerateID())
}

// GeneratePriceSnapshotID 生成价格快照ID
func GenerateSnapshotID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_PRICE_SNAPSHOT, GenerateID())
}

// GenerateRuleVersionID 生成规则版本ID
func GenerateRuleVersionID() string {
	return fmt.Sprintf("%v%v", ID_PREFIX_RULE_VERSION, GenerateID())
}

// GenerateSessionID 生成会话ID
func GenerateSessionID() string {
	return fmt.Sprintf("session%v%v", time.Now().UnixMilli(), GenerateShortID())
}

// GenerateSalt 生成加密盐值
func GenerateSalt() string {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(salt)
}

// GenerateInviteCode 生成10位随机字母数字组合的邀请码
func GenerateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 10

	b := make([]byte, codeLength)
	for i := range b {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果获取随机数失败，使用时间作为备用
			return fmt.Sprintf("%010d", time.Now().UnixNano())[:10]
		}
		b[i] = charset[randNum.Int64()]
	}
	return string(b)
}

// GenerateVerifyCode 生成6位数字验证码
func GenerateVerifyCode() string {
	return generateRandomNumber(6)
}

// generateRandomNumber 生成指定长度的随机数字字符串
func generateRandomNumber(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果获取随机数失败，使用简单的循环
			result[i] = charset[i%len(charset)]
		} else {
			result[i] = charset[randNum.Int64()]
		}
	}
	return string(result)
}
