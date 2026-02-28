package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// VerifyCodeService handles verification code operations
type VerifyCodeService struct {
	config     *config.VerifyCodeConfig
	msgService *MessageService
}

var (
	verifyCodeService *VerifyCodeService
	serviceLock       sync.Once
)

// SetupVerifyCodeService initializes verify code service with proper email service
func SetupVerifyCodeService() {
	cfg := config.Get()

	// Setup email service
	SetupEmailService()

	// Ensure VerifyCode config exists with defaults
	verifyCodeConfig := cfg.VerifyCode
	if verifyCodeConfig == nil {
		log.Get().Warn("[OTP] VerifyCode config is nil — using safe defaults (BypassOTP=false, real SMS will be sent)")
		verifyCodeConfig = &config.VerifyCodeConfig{
			Length:       4,
			Expiration:   5,
			SendInterval: 60,
			MaxSendTimes: 10,
			BypassOTP:    false,
		}
	}

	verifyCodeService = &VerifyCodeService{
		config:     verifyCodeConfig,
		msgService: GetMessageService(),
	}
}

func GetVerifyCodeService() *VerifyCodeService {
	serviceLock.Do(func() {
		SetupVerifyCodeService()
	})
	// If setup failed, try again (sync.Once won't retry)
	if verifyCodeService == nil {
		SetupVerifyCodeService()
	}
	return verifyCodeService
}

// GenerateCode generates a random verification code
func (s *VerifyCodeService) GenerateCode() string {
	// Defensive nil check
	length := 4
	if s.config != nil {
		length = s.config.Length
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result string
	for i := 0; i < length; i++ {
		result += fmt.Sprintf("%d", r.Intn(10))
	}
	return result
}

// getInt64FromCache gets int64 value from cache
func (s *VerifyCodeService) getInt64FromCache(key string) (int64, error) {
	value, err := models.GetCache(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

// setInt64ToCache sets int64 value to cache
func (s *VerifyCodeService) setInt64ToCache(key string, value int64, expiration time.Duration) error {
	return models.SetCache(key, strconv.FormatInt(value, 10), expiration)
}

// SendVerifyCode sends verification code via email or SMS
func (s *VerifyCodeService) SendVerifyCode(contactType, contact, user_type, purpose, language string) (protocol.ErrorCode, int) {
	// CRITICAL: Check if service itself is nil (happens if setup failed)
	if s == nil {
		log.Get().Error("VerifyCodeService is nil - service not initialized")
		return protocol.SystemError, 0
	}
	// Defensive nil check - ensure config is always initialized
	if s.config == nil {
		log.Get().Warn("[OTP] Config nil at send time — using safe defaults (BypassOTP=false)")
		s.config = &config.VerifyCodeConfig{
			Length:       4,
			Expiration:   5,
			SendInterval: 60,
			MaxSendTimes: 10,
			BypassOTP:    false,
		}
	}

	// Validate contact type
	if contactType != protocol.MsgChannelEmail && contactType != protocol.MsgChannelSms {
		return protocol.InvalidVerificationMethod, 0
	}
	if contactType == protocol.MsgChannelSms {
		normalized, ok := normalizeSMSPhone(contact)
		if !ok {
			log.Get().Warnf("[OTP] Invalid phone format rejected: raw=%s purpose=%s user_type=%s", contact, purpose, user_type)
			return protocol.InvalidParams, 0
		}
		contact = normalized
	}

	// Check send frequency
	lastTimeKey := fmt.Sprintf("%s_verify_code_%v_%v_%s_time", contactType, purpose, user_type, contact)
	lastTime, err := s.getInt64FromCache(lastTimeKey)
	if err == nil && lastTime > 0 {
		if time.Now().Unix()-lastTime < int64(s.config.SendInterval) {
			remainingSeconds := s.config.SendInterval - int(time.Now().Unix()-lastTime)
			return protocol.VerificationCooldown, remainingSeconds
		}
	}
	isSandbox := config.Get().IsSandbox()
	sandboxReason := ""
	if isSandbox {
		sandboxReason = "env=sandbox"
	}
	if contactType == protocol.MsgChannelSms && strings.HasPrefix(contact, "+86") {
		isSandbox = true
		sandboxReason = "china_number"
	}
	// Universal Bypass: If enabled in config, force sandbox mode for ALL SMS
	if s.config.BypassOTP && contactType == protocol.MsgChannelSms {
		isSandbox = true
		sandboxReason = "bypass_otp=true"
	}

	log.Get().Infof("[OTP] Preparing code for %s via %s (purpose=%s, user_type=%s, sandbox=%v reason=%s, env=%s)",
		contact, contactType, purpose, user_type, isSandbox, sandboxReason, config.Get().Env)

	// Generate verification code
	var code string
	if isSandbox {
		// For sandbox users, always use "1234"
		code = "1234"
	} else {
		// For real users, generate a random code
		code = s.GenerateCode()
	}

	// Store verification code in cache (convert minutes to seconds)
	codeKey := fmt.Sprintf("%s_verify_code_%v_%v_%s", contactType, purpose, user_type, contact)
	if err := models.SetCache(codeKey, code, time.Duration(s.config.Expiration)*time.Minute); err != nil {
		log.Get().Errorf("[OTP] Failed to store code in cache for %s: %v", contact, err)
		return protocol.CacheError, 0
	}
	if !isSandbox {
		msg := &Message{
			Type:     protocol.MsgTypeVerifyCode,
			Channels: []string{contactType},
			Language: language,
			To:       contact,
			Params: map[string]any{
				"to":         contact,
				"code":       code,
				"expiration": s.config.Expiration,
			},
		}
		if err := s.msgService.SendMessage(msg); err != nil {
			log.Get().Errorf("[OTP] SEND FAILED to %s via %s: %v", contact, contactType, err)
			return protocol.VerificationCodeSendFailed, 0
		}
		log.Get().Infof("[OTP] Code sent successfully to %s via %s", contact, contactType)
	} else {
		log.Get().Infof("[OTP] Sandbox mode — code %s stored for %s but NOT sent via SMS", code, contact)
	}
	log.Get().Infof("[OTP] Verification code %s prepared for %s (purpose=%v, user_type=%v)", code, contact, purpose, user_type)
	// Update send records
	s.setInt64ToCache(lastTimeKey, time.Now().Unix(), 24*time.Hour)

	return protocol.Success, 0
}

// SendEmailCode sends verification code via email
func (s *VerifyCodeService) SendEmailCode(email, user_type, purpose, language string) (protocol.ErrorCode, int) {
	return s.SendVerifyCode(protocol.MsgChannelEmail, email, user_type, purpose, language)
}

// SendSMSCode sends verification code via SMS
func (s *VerifyCodeService) SendSMSCode(phone, user_type, purpose, language string) (protocol.ErrorCode, int) {
	return s.SendVerifyCode(protocol.MsgChannelSms, phone, user_type, purpose, language)
}

// VerifyCode verifies a verification code for either email or SMS
func (s *VerifyCodeService) VerifyCode(contactType, purpose, user_type, contact, code string) bool {
	// Validate contact type
	if contactType != protocol.MsgChannelEmail && contactType != protocol.MsgChannelSms {
		return false
	}
	if contactType == protocol.MsgChannelSms {
		normalized, ok := normalizeSMSPhone(contact)
		if !ok {
			return false
		}
		contact = normalized
	}

	codeKey := fmt.Sprintf("%s_verify_code_%v_%v_%s", contactType, purpose, user_type, contact)
	var storedCode string
	storedCode, err := models.GetCache(codeKey)
	if err != nil {
		return false
	}

	if storedCode != code {
		return false
	}

	// Delete the code after successful verification
	models.DelCache(codeKey)
	return true
}

// VerifyEmailCode verifies email verification code
func (s *VerifyCodeService) VerifyEmailCode(purpose, user_type, email, code string) bool {
	return s.VerifyCode(protocol.MsgChannelEmail, purpose, user_type, email, code)
}

// VerifySMSCode verifies SMS verification code
func (s *VerifyCodeService) VerifySMSCode(purpose, user_type, phone, code string) bool {
	return s.VerifyCode(protocol.MsgChannelSms, purpose, user_type, phone, code)
}

// normalizeSMSPhone converts various user inputs into a strict E.164-like value.
// Returns false when the phone is not valid enough for OTP delivery.
func normalizeSMSPhone(raw string) (string, bool) {
	phone := strings.TrimSpace(raw)
	if phone == "" {
		return "", false
	}
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Convert international prefix 00 -> +
	if strings.HasPrefix(phone, "00") {
		phone = "+" + phone[2:]
	}
	// Rwanda local format 07XXXXXXXX -> +2507XXXXXXXX
	if strings.HasPrefix(phone, "0") && len(phone) == 10 {
		phone = "+250" + phone[1:]
	}
	// Rwanda with country digits but no plus
	if strings.HasPrefix(phone, "250") && len(phone) == 12 {
		phone = "+" + phone
	}
	// Reject accidental trunk zero after +250 (e.g. +25007...).
	// This format frequently causes carrier-side delivery failures.
	if strings.HasPrefix(phone, "+2500") {
		return "", false
	}
	if !strings.HasPrefix(phone, "+") {
		return "", false
	}
	digits := phone[1:]
	for _, ch := range digits {
		if ch < '0' || ch > '9' {
			return "", false
		}
	}
	// E.164 max digits after "+" is 15
	if len(digits) < 8 || len(digits) > 15 {
		return "", false
	}
	// Rwanda-specific hard checks to prevent malformed +2500... numbers
	if strings.HasPrefix(digits, "250") {
		if len(digits) != 12 {
			return "", false
		}
		// Mobile ranges in Rwanda are currently under 7XXXXXXXX.
		if digits[3] != '7' {
			return "", false
		}
	}
	return "+" + digits, true
}
