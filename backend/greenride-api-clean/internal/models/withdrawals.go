package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// Withdrawal 提现记录表 - 管理用户提现申请和处理
type Withdrawal struct {
	ID           int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	WithdrawalID string `json:"withdrawal_id" gorm:"column:withdrawal_id;type:varchar(64);uniqueIndex"`
	Salt         string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*WithdrawalValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type WithdrawalValues struct {
	// 账户信息
	AccountID *string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	UserID    *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType  *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index"`

	// 提现信息
	Amount    *float64 `json:"amount" gorm:"column:amount;type:decimal(15,2)"`                      // 提现金额
	FeeAmount *float64 `json:"fee_amount" gorm:"column:fee_amount;type:decimal(10,2);default:0.00"` // 手续费
	NetAmount *float64 `json:"net_amount" gorm:"column:net_amount;type:decimal(15,2)"`              // 实际到账金额
	Currency  *string  `json:"currency" gorm:"column:currency;type:varchar(10);default:'RWF'"`

	// 提现方式
	WithdrawalMethod *string `json:"withdrawal_method" gorm:"column:withdrawal_method;type:varchar(50)"` // bank_transfer, paypal, stripe, alipay, wechat
	DestinationType  *string `json:"destination_type" gorm:"column:destination_type;type:varchar(32)"`   // bank_account, paypal, stripe_account, alipay, wechat

	// 银行信息
	BankAccountName   *string `json:"bank_account_name" gorm:"column:bank_account_name;type:varchar(255)"`
	BankAccountNumber *string `json:"bank_account_number" gorm:"column:bank_account_number;type:varchar(100)"` // 加密存储
	BankName          *string `json:"bank_name" gorm:"column:bank_name;type:varchar(100)"`
	BankBranch        *string `json:"bank_branch" gorm:"column:bank_branch;type:varchar(100)"`
	BankSwiftCode     *string `json:"bank_swift_code" gorm:"column:bank_swift_code;type:varchar(20)"`
	BankRoutingNumber *string `json:"bank_routing_number" gorm:"column:bank_routing_number;type:varchar(20)"`

	// 第三方账户信息
	PaypalEmail     *string `json:"paypal_email" gorm:"column:paypal_email;type:varchar(255)"`
	StripeAccountID *string `json:"stripe_account_id" gorm:"column:stripe_account_id;type:varchar(255)"`
	AlipayAccount   *string `json:"alipay_account" gorm:"column:alipay_account;type:varchar(100)"`
	WechatAccount   *string `json:"wechat_account" gorm:"column:wechat_account;type:varchar(100)"`

	// 处理状态
	Status      *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'pending'"` // pending, processing, completed, failed, cancelled, rejected
	RequestedAt *int64  `json:"requested_at" gorm:"column:requested_at"`
	ApprovedAt  *int64  `json:"approved_at" gorm:"column:approved_at"`
	ProcessedAt *int64  `json:"processed_at" gorm:"column:processed_at"`
	CompletedAt *int64  `json:"completed_at" gorm:"column:completed_at"`
	FailedAt    *int64  `json:"failed_at" gorm:"column:failed_at"`

	// 审批信息
	ApprovalStatus  *string `json:"approval_status" gorm:"column:approval_status;type:varchar(32);default:'pending'"` // pending, approved, rejected, auto_approved
	ApprovedBy      *string `json:"approved_by" gorm:"column:approved_by;type:varchar(64)"`
	RejectionReason *string `json:"rejection_reason" gorm:"column:rejection_reason;type:varchar(500)"`
	ApprovalNotes   *string `json:"approval_notes" gorm:"column:approval_notes;type:text"`

	// 处理信息
	ProcessedBy           *string `json:"processed_by" gorm:"column:processed_by;type:varchar(64)"`
	ProcessingReference   *string `json:"processing_reference" gorm:"column:processing_reference;type:varchar(255)"`
	ExternalTransactionID *string `json:"external_transaction_id" gorm:"column:external_transaction_id;type:varchar(255)"`
	FailureReason         *string `json:"failure_reason" gorm:"column:failure_reason;type:text"`

	// 风控信息
	RiskScore            *float64 `json:"risk_score" gorm:"column:risk_score;type:decimal(5,2);default:0.00"`
	RiskFlags            *string  `json:"risk_flags" gorm:"column:risk_flags;type:json"`
	RequiresManualReview *bool    `json:"requires_manual_review" gorm:"column:requires_manual_review;default:false"`
	ReviewNotes          *string  `json:"review_notes" gorm:"column:review_notes;type:text"`

	// 余额验证
	AccountBalanceBefore *float64 `json:"account_balance_before" gorm:"column:account_balance_before;type:decimal(15,2)"`
	AccountBalanceAfter  *float64 `json:"account_balance_after" gorm:"column:account_balance_after;type:decimal(15,2)"`

	// 限额检查
	DailyWithdrawalCount    *int     `json:"daily_withdrawal_count" gorm:"column:daily_withdrawal_count;default:1"`
	DailyWithdrawalAmount   *float64 `json:"daily_withdrawal_amount" gorm:"column:daily_withdrawal_amount;type:decimal(15,2)"`
	MonthlyWithdrawalAmount *float64 `json:"monthly_withdrawal_amount" gorm:"column:monthly_withdrawal_amount;type:decimal(15,2)"`

	// 元数据
	Metadata   *string `json:"metadata" gorm:"column:metadata;type:json"`
	UserNotes  *string `json:"user_notes" gorm:"column:user_notes;type:varchar(500)"`
	AdminNotes *string `json:"admin_notes" gorm:"column:admin_notes;type:text"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Withdrawal) TableName() string {
	return "t_withdrawals"
}

// 提现状态常量
const (
	WithdrawalStatusPending    = "pending"
	WithdrawalStatusProcessing = "processing"
	WithdrawalStatusCompleted  = "completed"
	WithdrawalStatusFailed     = "failed"
	WithdrawalStatusCancelled  = "cancelled"
	WithdrawalStatusRejected   = "rejected"
)

// 审批状态常量
const (
	ApprovalStatusPending      = "pending"
	ApprovalStatusApproved     = "approved"
	ApprovalStatusRejected     = "rejected"
	ApprovalStatusAutoApproved = "auto_approved"
)

// 提现方式常量
const (
	WithdrawalMethodBankTransfer = "bank_transfer"
	WithdrawalMethodPaypal       = "paypal"
	WithdrawalMethodStripe       = "stripe"
	WithdrawalMethodAlipay       = "alipay"
	WithdrawalMethodWechat       = "wechat"
)

// 目标类型常量
const (
	DestinationTypeBankAccount   = "bank_account"
	DestinationTypePaypal        = "paypal"
	DestinationTypeStripeAccount = "stripe_account"
	DestinationTypeAlipay        = "alipay"
	DestinationTypeWechat        = "wechat"
)

// 创建新的提现记录对象
func NewWithdrawalV2() *Withdrawal {
	now := utils.TimeNowMilli()
	return &Withdrawal{
		WithdrawalID: utils.GenerateWithdrawalID(),
		Salt:         utils.GenerateSalt(),
		WithdrawalValues: &WithdrawalValues{
			UserType:             utils.StringPtr(protocol.UserTypePassenger),
			Currency:             utils.StringPtr("RWF"),
			Status:               utils.StringPtr(WithdrawalStatusPending),
			ApprovalStatus:       utils.StringPtr(ApprovalStatusPending),
			FeeAmount:            utils.Float64Ptr(0.00),
			RiskScore:            utils.Float64Ptr(0.00),
			RequiresManualReview: utils.BoolPtr(false),
			DailyWithdrawalCount: utils.IntPtr(1),
			RequestedAt:          &now,
		},
	}
}

// SetValues 更新WithdrawalV2Values中的非nil值
func (w *WithdrawalValues) SetValues(values *WithdrawalValues) {
	if values == nil {
		return
	}

	if values.AccountID != nil {
		w.AccountID = values.AccountID
	}
	if values.UserID != nil {
		w.UserID = values.UserID
	}
	if values.UserType != nil {
		w.UserType = values.UserType
	}
	if values.Amount != nil {
		w.Amount = values.Amount
	}
	if values.NetAmount != nil {
		w.NetAmount = values.NetAmount
	}
	if values.WithdrawalMethod != nil {
		w.WithdrawalMethod = values.WithdrawalMethod
	}
	if values.DestinationType != nil {
		w.DestinationType = values.DestinationType
	}
	if values.Status != nil {
		w.Status = values.Status
	}
	if values.ApprovalStatus != nil {
		w.ApprovalStatus = values.ApprovalStatus
	}
	if values.BankAccountName != nil {
		w.BankAccountName = values.BankAccountName
	}
	if values.BankAccountNumber != nil {
		w.BankAccountNumber = values.BankAccountNumber
	}
	if values.BankName != nil {
		w.BankName = values.BankName
	}
	if values.UserNotes != nil {
		w.UserNotes = values.UserNotes
	}
	if values.Metadata != nil {
		w.Metadata = values.Metadata
	}
	if values.UpdatedAt > 0 {
		w.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (w *WithdrawalValues) GetAccountID() string {
	if w.AccountID == nil {
		return ""
	}
	return *w.AccountID
}

func (w *WithdrawalValues) GetUserID() string {
	if w.UserID == nil {
		return ""
	}
	return *w.UserID
}

func (w *WithdrawalValues) GetAmount() float64 {
	if w.Amount == nil {
		return 0.00
	}
	return *w.Amount
}

func (w *WithdrawalValues) GetFeeAmount() float64 {
	if w.FeeAmount == nil {
		return 0.00
	}
	return *w.FeeAmount
}

func (w *WithdrawalValues) GetNetAmount() float64 {
	if w.NetAmount == nil {
		return 0.00
	}
	return *w.NetAmount
}

func (w *WithdrawalValues) GetCurrency() string {
	if w.Currency == nil {
		return "RWF"
	}
	return *w.Currency
}

func (w *WithdrawalValues) GetWithdrawalMethod() string {
	if w.WithdrawalMethod == nil {
		return ""
	}
	return *w.WithdrawalMethod
}

func (w *WithdrawalValues) GetDestinationType() string {
	if w.DestinationType == nil {
		return ""
	}
	return *w.DestinationType
}

func (w *WithdrawalValues) GetStatus() string {
	if w.Status == nil {
		return WithdrawalStatusPending
	}
	return *w.Status
}

func (w *WithdrawalValues) GetApprovalStatus() string {
	if w.ApprovalStatus == nil {
		return ApprovalStatusPending
	}
	return *w.ApprovalStatus
}

func (w *WithdrawalValues) GetRiskScore() float64 {
	if w.RiskScore == nil {
		return 0.00
	}
	return *w.RiskScore
}

func (w *WithdrawalValues) GetRequiresManualReview() bool {
	if w.RequiresManualReview == nil {
		return false
	}
	return *w.RequiresManualReview
}

func (w *WithdrawalValues) GetDailyWithdrawalCount() int {
	if w.DailyWithdrawalCount == nil {
		return 1
	}
	return *w.DailyWithdrawalCount
}

// Setter 方法
func (w *WithdrawalValues) SetAccountID(accountID string) *WithdrawalValues {
	w.AccountID = &accountID
	return w
}

func (w *WithdrawalValues) SetUserID(userID string) *WithdrawalValues {
	w.UserID = &userID
	return w
}

func (w *WithdrawalValues) SetAmount(amount float64) *WithdrawalValues {
	w.Amount = &amount
	return w
}

func (w *WithdrawalValues) SetFeeAmount(fee float64) *WithdrawalValues {
	w.FeeAmount = &fee
	return w
}

func (w *WithdrawalValues) SetNetAmount(netAmount float64) *WithdrawalValues {
	w.NetAmount = &netAmount
	return w
}

func (w *WithdrawalValues) SetWithdrawalMethod(method string) *WithdrawalValues {
	w.WithdrawalMethod = &method
	return w
}

func (w *WithdrawalValues) SetDestinationType(destType string) *WithdrawalValues {
	w.DestinationType = &destType
	return w
}

func (w *WithdrawalValues) SetStatus(status string) *WithdrawalValues {
	w.Status = &status
	return w
}

func (w *WithdrawalValues) SetApprovalStatus(status string) *WithdrawalValues {
	w.ApprovalStatus = &status
	return w
}

func (w *WithdrawalValues) SetBankInfo(accountName, accountNumber, bankName, branch string) *WithdrawalValues {
	w.BankAccountName = &accountName
	w.BankAccountNumber = &accountNumber
	w.BankName = &bankName
	if branch != "" {
		w.BankBranch = &branch
	}
	return w
}

func (w *WithdrawalValues) SetPaypalEmail(email string) *WithdrawalValues {
	w.PaypalEmail = &email
	return w
}

func (w *WithdrawalValues) SetRiskScore(score float64) *WithdrawalValues {
	w.RiskScore = &score
	return w
}

func (w *WithdrawalValues) SetRequiresManualReview(required bool) *WithdrawalValues {
	w.RequiresManualReview = &required
	return w
}

func (w *WithdrawalValues) SetBalanceSnapshot(before, after float64) *WithdrawalValues {
	w.AccountBalanceBefore = &before
	w.AccountBalanceAfter = &after
	return w
}

// 业务方法
func (w *Withdrawal) IsPending() bool {
	return w.GetStatus() == WithdrawalStatusPending
}

func (w *Withdrawal) IsProcessing() bool {
	return w.GetStatus() == WithdrawalStatusProcessing
}

func (w *Withdrawal) IsCompleted() bool {
	return w.GetStatus() == WithdrawalStatusCompleted
}

func (w *Withdrawal) IsFailed() bool {
	return w.GetStatus() == WithdrawalStatusFailed
}

func (w *Withdrawal) IsCancelled() bool {
	return w.GetStatus() == WithdrawalStatusCancelled
}

func (w *Withdrawal) IsRejected() bool {
	return w.GetStatus() == WithdrawalStatusRejected
}

func (w *Withdrawal) IsApproved() bool {
	status := w.GetApprovalStatus()
	return status == ApprovalStatusApproved || status == ApprovalStatusAutoApproved
}

func (w *Withdrawal) IsAwaitingApproval() bool {
	return w.GetApprovalStatus() == ApprovalStatusPending
}

func (w *Withdrawal) IsHighRisk() bool {
	return w.GetRiskScore() >= 80.0
}

func (w *Withdrawal) CanCancel() bool {
	status := w.GetStatus()
	return status == WithdrawalStatusPending || (status == WithdrawalStatusProcessing && w.IsAwaitingApproval())
}

func (w *Withdrawal) CanProcess() bool {
	return w.IsApproved() && w.GetStatus() == WithdrawalStatusPending
}

// 状态更新方法
func (w *WithdrawalValues) Approve(adminID string) *WithdrawalValues {
	w.SetApprovalStatus(ApprovalStatusApproved)
	w.ApprovedBy = &adminID
	now := utils.TimeNowMilli()
	w.ApprovedAt = &now
	return w
}

func (w *WithdrawalValues) Reject(adminID, reason string) *WithdrawalValues {
	w.SetApprovalStatus(ApprovalStatusRejected)
	w.SetStatus(WithdrawalStatusRejected)
	w.ApprovedBy = &adminID
	w.RejectionReason = &reason
	now := utils.TimeNowMilli()
	w.ApprovedAt = &now
	return w
}

func (w *WithdrawalValues) StartProcessing(adminID string) *WithdrawalValues {
	w.SetStatus(WithdrawalStatusProcessing)
	w.ProcessedBy = &adminID
	now := utils.TimeNowMilli()
	w.ProcessedAt = &now
	return w
}

func (w *WithdrawalValues) Complete(externalTxID string) *WithdrawalValues {
	w.SetStatus(WithdrawalStatusCompleted)
	if externalTxID != "" {
		w.ExternalTransactionID = &externalTxID
	}
	now := utils.TimeNowMilli()
	w.CompletedAt = &now
	return w
}

func (w *WithdrawalValues) Fail(reason string) *WithdrawalValues {
	w.SetStatus(WithdrawalStatusFailed)
	w.FailureReason = &reason
	now := utils.TimeNowMilli()
	w.FailedAt = &now
	return w
}

func (w *WithdrawalValues) Cancel(reason string) *WithdrawalValues {
	w.SetStatus(WithdrawalStatusCancelled)
	if reason != "" {
		w.AdminNotes = &reason
	}
	return w
}

func (w *WithdrawalValues) SetProcessingReference(reference string) *WithdrawalValues {
	w.ProcessingReference = &reference
	return w
}

// 风控相关方法
func (w *WithdrawalValues) FlagForManualReview(reason string) *WithdrawalValues {
	w.SetRequiresManualReview(true)
	if reason != "" {
		w.ReviewNotes = &reason
	}
	return w
}

func (w *WithdrawalValues) AddRiskFlag(flag string) error {
	var flags []string
	if w.RiskFlags != nil {
		if err := utils.FromJSON(*w.RiskFlags, &flags); err != nil {
			return fmt.Errorf("failed to parse existing risk flags: %v", err)
		}
	}

	// 避免重复添加
	for _, existingFlag := range flags {
		if existingFlag == flag {
			return nil
		}
	}

	flags = append(flags, flag)
	flagsJSON, err := utils.ToJSON(flags)
	if err != nil {
		return fmt.Errorf("failed to marshal risk flags: %v", err)
	}

	w.RiskFlags = &flagsJSON
	return nil
}

// 计算手续费
func (w *WithdrawalValues) CalculateFee(feeRate float64, fixedFee float64) *WithdrawalValues {
	amount := w.GetAmount()

	// 按比例手续费
	percentageFee := amount * feeRate

	// 总手续费 = 固定费用 + 按比例费用
	totalFee := fixedFee + percentageFee

	// 实际到账金额
	netAmount := amount - totalFee

	w.SetFeeAmount(totalFee)
	w.SetNetAmount(netAmount)

	return w
}

// 设置限额统计
func (w *WithdrawalValues) SetDailyStats(count int, amount float64) *WithdrawalValues {
	w.DailyWithdrawalCount = &count
	w.DailyWithdrawalAmount = &amount
	return w
}

func (w *WithdrawalValues) SetMonthlyAmount(amount float64) *WithdrawalValues {
	w.MonthlyWithdrawalAmount = &amount
	return w
}

// 检查是否超出限额
func (w *WithdrawalValues) ExceedsDailyLimit(dailyLimit float64) bool {
	if w.DailyWithdrawalAmount == nil {
		return false
	}
	return *w.DailyWithdrawalAmount > dailyLimit
}

func (w *WithdrawalValues) ExceedsMonthlyLimit(monthlyLimit float64) bool {
	if w.MonthlyWithdrawalAmount == nil {
		return false
	}
	return *w.MonthlyWithdrawalAmount > monthlyLimit
}

// 创建银行转账提现
func NewBankTransferWithdrawal(accountID, userID string, amount float64, bankInfo map[string]string) *Withdrawal {
	withdrawal := NewWithdrawalV2()
	withdrawal.SetAccountID(accountID).
		SetUserID(userID).
		SetAmount(amount).
		SetWithdrawalMethod(WithdrawalMethodBankTransfer).
		SetDestinationType(DestinationTypeBankAccount).
		SetBankInfo(
			bankInfo["account_name"],
			bankInfo["account_number"],
			bankInfo["bank_name"],
			bankInfo["branch"],
		)

	if swiftCode, exists := bankInfo["swift_code"]; exists {
		withdrawal.BankSwiftCode = &swiftCode
	}

	return withdrawal
}

// 创建PayPal提现
func NewPaypalWithdrawal(accountID, userID string, amount float64, email string) *Withdrawal {
	withdrawal := NewWithdrawalV2()
	withdrawal.SetAccountID(accountID).
		SetUserID(userID).
		SetAmount(amount).
		SetWithdrawalMethod(WithdrawalMethodPaypal).
		SetDestinationType(DestinationTypePaypal).
		SetPaypalEmail(email)

	return withdrawal
}
