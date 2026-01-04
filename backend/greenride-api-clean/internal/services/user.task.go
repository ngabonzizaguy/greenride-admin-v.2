package services

import (
	"context"
	"fmt"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/task"
	"greenride/internal/utils"

	"github.com/spf13/cast"
)

const (
	TaskUserRegisteredHandler = "user_registered_handler"
)

// InitUserTaskHandlers 初始化用户任务处理器
func InitUserTaskHandlers() {
	task.RegisterHandler(protocol.UserRegisteredHandler, HandleUserRegistered)
}

// HandleUserRegistered 处理用户注册信号
func HandleUserRegistered(ctx context.Context, params protocol.MapData) error {
	// 尝试从不同的参数名获取用户ID
	userID := cast.ToString(params["biz_id"])
	if userID == "" {
		userID = cast.ToString(params["user_id"])
	}
	if userID == "" {
		return fmt.Errorf("invalid user_id parameter")
	}

	logger := log.Get().WithField("user_id", userID)

	// 1. 给新注册用户发放欢迎优惠券
	if err := IssueWelcomeCouponForUser(userID); err != nil {
		logger.Errorf("Failed to issue welcome coupon: %v", err)
		// 错误不中断流程，继续处理邀请关系
	}

	// 2. 处理邀请关系 - 给邀请人发放推荐优惠券
	if err := IssueReferralCouponForInviter(userID); err != nil {
		logger.Errorf("Failed to issue referral coupon to inviter: %v", err)
		// 错误记录日志但不影响整体流程
	}

	return nil
}

// IssueWelcomeCouponForUser 为新用户发放欢迎优惠券
func IssueWelcomeCouponForUser(userID string) error {
	logger := log.Get().WithField("user_id", userID)

	// 1. 获取配置（带默认值）
	promotionConfig := config.GetPromotionConfig()

	// 2. 检查是否启用欢迎优惠券
	if promotionConfig.EnableWelcomeCoupon != config.StatusOn {
		logger.Debug("Welcome coupon issuance is disabled")
		return nil
	}

	// 3. 获取用户信息并检查是否为乘客类型
	user := models.GetUserByID(userID)
	if user == nil {
		return fmt.Errorf("user not found: %s", userID)
	}

	if !user.IsPassenger() {
		logger.Debug("Skip coupon issuance for non-passenger user")
		return nil
	}

	// 4. 防重复检查
	if models.CheckUserHasWelcomeCoupon(userID) {
		logger.Info("User already has welcome coupon")
		return nil
	}

	// 5. 获取或创建欢迎优惠券模板
	var promotion *models.Promotion
	if promotionConfig.WelcomeCouponCode != "" {
		promotion = models.GetWelcomePromotionTemplate(promotionConfig.WelcomeCouponCode)
	}

	// 如果指定的优惠券模板不存在，使用默认模板
	if promotion == nil {
		promotion = models.GetOrCreateDefaultWelcomePromotion()
		if promotion == nil {
			return fmt.Errorf("failed to get or create welcome promotion template")
		}
		logger.Infof("Using default welcome promotion template: %s", promotion.GetCode())
	}

	// 6. 创建用户优惠券实例
	userPromotion := models.CreateWelcomePromoForUser(userID, promotion)

	// 7. 设置过期时间（根据配置）
	if promotionConfig.WelcomeCouponValidDays > 0 {
		expiredAt := utils.TimeNowMilli() + int64(promotionConfig.WelcomeCouponValidDays*24*3600*1000)
		userPromotion.SetExpiredAt(expiredAt)
	}

	// 8. 保存到数据库
	if err := models.CreateUserPromotionInDB(userPromotion); err != nil {
		return fmt.Errorf("failed to create user promotion: %v", err)
	}

	logger.Infof("Welcome coupon issued successfully: %s", userPromotion.GetCode())
	return nil
}

// IssueReferralCouponForInviter 为邀请人发放推荐优惠券
func IssueReferralCouponForInviter(invitedUserID string) error {
	logger := log.Get().WithField("invited_user_id", invitedUserID)

	// 1. 获取配置
	promotionConfig := config.GetPromotionConfig()

	// 2. 检查是否启用推荐优惠券
	if promotionConfig.EnableReferralCoupon != config.StatusOn {
		logger.Debug("Referral coupon issuance is disabled")
		return nil
	}

	// 3. 获取新用户信息
	invitedUser := models.GetUserByID(invitedUserID)
	if invitedUser == nil {
		return fmt.Errorf("invited user not found: %s", invitedUserID)
	}

	// 5. 获取邀请人信息
	inviter := models.GetUserByID(invitedUser.GetInvitedBy())
	if inviter == nil {
		logger.Warnf("Inviter not found for user: %s", invitedUser.GetInvitedBy())
		return nil
	}

	// 只有乘客才能获得推荐优惠券
	if !inviter.IsPassenger() {
		logger.Debug("Skip referral coupon issuance for non-passenger inviter")
		return nil
	}

	logger = logger.WithField("inviter_id", inviter.UserID)

	// 6. 获取或创建推荐优惠券模板
	var promotion *models.Promotion
	if promotionConfig.ReferralCouponCode != "" {
		promotion = models.GetReferralPromotionTemplate(promotionConfig.ReferralCouponCode)
	}

	// 如果指定的优惠券模板不存在，使用默认模板
	if promotion == nil {
		promotion = models.GetOrCreateDefaultReferralPromotion()
		if promotion == nil {
			return fmt.Errorf("failed to get or create referral promotion template")
		}
		logger.Infof("Using default referral promotion template: %s", promotion.GetCode())
	}

	// 7. 创建用户优惠券实例
	userPromotion := models.CreateReferralPromoForUser(inviter.UserID, invitedUserID, promotion)

	// 8. 设置过期时间
	if promotionConfig.ReferralCouponValidDays > 0 {
		expiredAt := utils.TimeNowMilli() + int64(promotionConfig.ReferralCouponValidDays*24*3600*1000)
		userPromotion.ExpiredAt = &expiredAt
	}

	// 9. 保存到数据库
	if err := models.CreateUserPromotionInDB(userPromotion); err != nil {
		return fmt.Errorf("failed to create referral promotion: %v", err)
	}

	// 10. 更新邀请统计
	if err := models.IncrementUserInviteCount(inviter.UserID); err != nil {
		logger.Warnf("Failed to increment invite count: %v", err)
		// 不影响优惠券发放，继续执行
	}

	logger.Infof("Referral coupon issued successfully to inviter %s: %s",
		inviter.UserID, userPromotion.GetCode())
	return nil
}
