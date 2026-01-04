package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// WalletTransaction 钱包交易表 - 记录所有钱包相关交易
type WalletTransaction struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TransactionID string `json:"transaction_id" gorm:"column:transaction_id;type:varchar(64);uniqueIndex"`
	Salt          string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*WalletTransactionValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type WalletTransactionValues struct {
	// 账户信息
	AccountID *string `json:"account_id" gorm:"column:account_id;type:varchar(64);index"`
	UserID    *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType  *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index"`

	// 交易信息
	Type      *string  `json:"type" gorm:"column:type;type:varchar(32);index"`         // income, expense, transfer, withdrawal, refund, bonus, penalty
	Category  *string  `json:"category" gorm:"column:category;type:varchar(32);index"` // ride_payment, ride_earning, tip, withdrawal_fee, platform_fee, bonus, refund
	Amount    *float64 `json:"amount" gorm:"column:amount;type:decimal(15,2)"`
	Currency  *string  `json:"currency" gorm:"column:currency;type:varchar(10);default:'RWF'"`
	FeeAmount *float64 `json:"fee_amount" gorm:"column:fee_amount;type:decimal(10,2);default:0.00"`

	// 关联信息
	RelatedID    *string `json:"related_id" gorm:"column:related_id;type:varchar(64);index"`       // 关联对象ID
	RelatedType  *string `json:"related_type" gorm:"column:related_type;type:varchar(64);index"`   // ride_order, payment, withdrawal, transfer
	PaymentID    *string `json:"payment_id" gorm:"column:payment_id;type:varchar(64);index"`       // 与t_payments.payment_id关联
	WithdrawalID *string `json:"withdrawal_id" gorm:"column:withdrawal_id;type:varchar(64);index"` // 提现记录ID
	TransferID   *string `json:"transfer_id" gorm:"column:transfer_id;type:varchar(64);index"`     // 转账记录ID

	// 对方信息(用于转账)
	CounterpartAccountID *string `json:"counterpart_account_id" gorm:"column:counterpart_account_id;type:varchar(64);index"`
	CounterpartUserID    *string `json:"counterpart_user_id" gorm:"column:counterpart_user_id;type:varchar(64);index"`
	CounterpartUserType  *string `json:"counterpart_user_type" gorm:"column:counterpart_user_type;type:varchar(32)"`

	// 交易状态
	Status      *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'completed'"` // pending, processing, completed, failed, cancelled, reversed
	ProcessedAt *int64  `json:"processed_at" gorm:"column:processed_at"`
	CompletedAt *int64  `json:"completed_at" gorm:"column:completed_at"`
	FailedAt    *int64  `json:"failed_at" gorm:"column:failed_at"`

	// 描述和备注
	Title       *string `json:"title" gorm:"column:title;type:varchar(255)"`
	Description *string `json:"description" gorm:"column:description;type:text"`
	AdminNotes  *string `json:"admin_notes" gorm:"column:admin_notes;type:text"`
	UserNotes   *string `json:"user_notes" gorm:"column:user_notes;type:varchar(500)"`

	// 风控信息
	RiskScore    *float64 `json:"risk_score" gorm:"column:risk_score;type:decimal(5,2);default:0.00"`
	IsSuspicious *bool    `json:"is_suspicious" gorm:"column:is_suspicious;default:false"`
	ReviewStatus *string  `json:"review_status" gorm:"column:review_status;type:varchar(32);default:'auto_approved'"` // auto_approved, pending_review, manual_approved, rejected
	ReviewedBy   *string  `json:"reviewed_by" gorm:"column:reviewed_by;type:varchar(64)"`
	ReviewedAt   *int64   `json:"reviewed_at" gorm:"column:reviewed_at"`

	// 元数据
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"`
	Tags     *string `json:"tags" gorm:"column:tags;type:json"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (WalletTransaction) TableName() string {
	return "t_wallet_transactions"
}

// 交易类型常量
const (
	TransactionTypeIncome     = "income"
	TransactionTypeExpense    = "expense"
	TransactionTypeTransfer   = "transfer"
	TransactionTypeWithdrawal = "withdrawal"
	TransactionTypeRefund     = "refund"
	TransactionTypeBonus      = "bonus"
	TransactionTypePenalty    = "penalty"
)

// 交易分类常量
const (
	TransactionCategoryRidePayment   = "ride_payment"
	TransactionCategoryRideEarning   = "ride_earning"
	TransactionCategoryTip           = "tip"
	TransactionCategoryWithdrawalFee = "withdrawal_fee"
	TransactionCategoryPlatformFee   = "platform_fee"
	TransactionCategoryBonus         = "bonus"
	TransactionCategoryRefund        = "refund"
	TransactionCategoryTopup         = "topup"
	TransactionCategoryTransfer      = "transfer"
)

// 交易状态常量
const (
	TransactionStatusPending    = "pending"
	TransactionStatusProcessing = "processing"
	TransactionStatusCompleted  = "completed"
	TransactionStatusFailed     = "failed"
	TransactionStatusCancelled  = "cancelled"
	TransactionStatusReversed   = "reversed"
)

// 审核状态常量
const (
	ReviewStatusAutoApproved   = "auto_approved"
	ReviewStatusPendingReview  = "pending_review"
	ReviewStatusManualApproved = "manual_approved"
	ReviewStatusRejected       = "rejected"
)

// 创建新的钱包交易对象
func NewWalletTransactionV2() *WalletTransaction {
	return &WalletTransaction{
		TransactionID: utils.GenerateWalletTransactionID(),
		Salt:          utils.GenerateSalt(),
		WalletTransactionValues: &WalletTransactionValues{
			UserType:     utils.StringPtr(protocol.UserTypePassenger),
			Currency:     utils.StringPtr("RWF"),
			Status:       utils.StringPtr(TransactionStatusCompleted),
			FeeAmount:    utils.Float64Ptr(0.00),
			RiskScore:    utils.Float64Ptr(0.00),
			IsSuspicious: utils.BoolPtr(false),
			ReviewStatus: utils.StringPtr(ReviewStatusAutoApproved),
		},
	}
}

// SetValues 更新WalletTransactionV2Values中的非nil值
func (t *WalletTransactionValues) SetValues(values *WalletTransactionValues) {
	if values == nil {
		return
	}

	if values.AccountID != nil {
		t.AccountID = values.AccountID
	}
	if values.UserID != nil {
		t.UserID = values.UserID
	}
	if values.UserType != nil {
		t.UserType = values.UserType
	}
	if values.Type != nil {
		t.Type = values.Type
	}
	if values.Category != nil {
		t.Category = values.Category
	}
	if values.Amount != nil {
		t.Amount = values.Amount
	}
	if values.Currency != nil {
		t.Currency = values.Currency
	}
	if values.Status != nil {
		t.Status = values.Status
	}
	if values.Title != nil {
		t.Title = values.Title
	}
	if values.Description != nil {
		t.Description = values.Description
	}
	if values.RelatedID != nil {
		t.RelatedID = values.RelatedID
	}
	if values.RelatedType != nil {
		t.RelatedType = values.RelatedType
	}
	if values.PaymentID != nil {
		t.PaymentID = values.PaymentID
	}
	if values.Metadata != nil {
		t.Metadata = values.Metadata
	}
	if values.UpdatedAt > 0 {
		t.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (t *WalletTransactionValues) GetAccountID() string {
	if t.AccountID == nil {
		return ""
	}
	return *t.AccountID
}

func (t *WalletTransactionValues) GetUserID() string {
	if t.UserID == nil {
		return ""
	}
	return *t.UserID
}

func (t *WalletTransactionValues) GetType() string {
	if t.Type == nil {
		return ""
	}
	return *t.Type
}

func (t *WalletTransactionValues) GetCategory() string {
	if t.Category == nil {
		return ""
	}
	return *t.Category
}

func (t *WalletTransactionValues) GetAmount() float64 {
	if t.Amount == nil {
		return 0.00
	}
	return *t.Amount
}

func (t *WalletTransactionValues) GetCurrency() string {
	if t.Currency == nil {
		return "RWF"
	}
	return *t.Currency
}

func (t *WalletTransactionValues) GetFeeAmount() float64 {
	if t.FeeAmount == nil {
		return 0.00
	}
	return *t.FeeAmount
}

func (t *WalletTransactionValues) GetStatus() string {
	if t.Status == nil {
		return TransactionStatusCompleted
	}
	return *t.Status
}

func (t *WalletTransactionValues) GetTitle() string {
	if t.Title == nil {
		return ""
	}
	return *t.Title
}

func (t *WalletTransactionValues) GetDescription() string {
	if t.Description == nil {
		return ""
	}
	return *t.Description
}

func (t *WalletTransactionValues) GetRelatedID() string {
	if t.RelatedID == nil {
		return ""
	}
	return *t.RelatedID
}

func (t *WalletTransactionValues) GetRelatedType() string {
	if t.RelatedType == nil {
		return ""
	}
	return *t.RelatedType
}

func (t *WalletTransactionValues) GetRiskScore() float64 {
	if t.RiskScore == nil {
		return 0.00
	}
	return *t.RiskScore
}

func (t *WalletTransactionValues) GetIsSuspicious() bool {
	if t.IsSuspicious == nil {
		return false
	}
	return *t.IsSuspicious
}

func (t *WalletTransactionValues) GetReviewStatus() string {
	if t.ReviewStatus == nil {
		return ReviewStatusAutoApproved
	}
	return *t.ReviewStatus
}

// Setter 方法
func (t *WalletTransactionValues) SetAccountID(accountID string) *WalletTransactionValues {
	t.AccountID = &accountID
	return t
}

func (t *WalletTransactionValues) SetUserID(userID string) *WalletTransactionValues {
	t.UserID = &userID
	return t
}

func (t *WalletTransactionValues) SetType(transactionType string) *WalletTransactionValues {
	t.Type = &transactionType
	return t
}

func (t *WalletTransactionValues) SetCategory(category string) *WalletTransactionValues {
	t.Category = &category
	return t
}

func (t *WalletTransactionValues) SetAmount(amount float64) *WalletTransactionValues {
	t.Amount = &amount
	return t
}

func (t *WalletTransactionValues) SetStatus(status string) *WalletTransactionValues {
	t.Status = &status
	return t
}

func (t *WalletTransactionValues) SetTitle(title string) *WalletTransactionValues {
	t.Title = &title
	return t
}

func (t *WalletTransactionValues) SetDescription(description string) *WalletTransactionValues {
	t.Description = &description
	return t
}

func (t *WalletTransactionValues) SetRelated(relatedType, relatedID string) *WalletTransactionValues {
	t.RelatedType = &relatedType
	t.RelatedID = &relatedID
	return t
}

func (t *WalletTransactionValues) SetPaymentID(paymentID string) *WalletTransactionValues {
	t.PaymentID = &paymentID
	return t
}

func (t *WalletTransactionValues) SetWithdrawalID(withdrawalID string) *WalletTransactionValues {
	t.WithdrawalID = &withdrawalID
	return t
}

func (t *WalletTransactionValues) SetCounterpart(accountID, userID, userType string) *WalletTransactionValues {
	t.CounterpartAccountID = &accountID
	t.CounterpartUserID = &userID
	t.CounterpartUserType = &userType
	return t
}

func (t *WalletTransactionValues) SetRiskScore(score float64) *WalletTransactionValues {
	t.RiskScore = &score
	return t
}

func (t *WalletTransactionValues) MarkAsSuspicious() *WalletTransactionValues {
	t.IsSuspicious = utils.BoolPtr(true)
	return t
}

// 业务方法
func (t *WalletTransaction) IsIncome() bool {
	return t.GetType() == TransactionTypeIncome
}

func (t *WalletTransaction) IsExpense() bool {
	return t.GetType() == TransactionTypeExpense
}

func (t *WalletTransaction) IsTransfer() bool {
	return t.GetType() == TransactionTypeTransfer
}

func (t *WalletTransaction) IsWithdrawal() bool {
	return t.GetType() == TransactionTypeWithdrawal
}

func (t *WalletTransaction) IsCompleted() bool {
	return t.GetStatus() == TransactionStatusCompleted
}

func (t *WalletTransaction) IsPending() bool {
	return t.GetStatus() == TransactionStatusPending
}

func (t *WalletTransaction) IsFailed() bool {
	return t.GetStatus() == TransactionStatusFailed
}

func (t *WalletTransaction) IsHighRisk() bool {
	return t.GetRiskScore() >= 80.0
}

func (t *WalletTransaction) RequiresReview() bool {
	return t.GetReviewStatus() == ReviewStatusPendingReview
}

// 状态更新方法
func (t *WalletTransactionValues) MarkAsProcessing() *WalletTransactionValues {
	t.SetStatus(TransactionStatusProcessing)
	now := utils.TimeNowMilli()
	t.ProcessedAt = &now
	return t
}

func (t *WalletTransactionValues) MarkAsCompleted() *WalletTransactionValues {
	t.SetStatus(TransactionStatusCompleted)
	now := utils.TimeNowMilli()
	t.CompletedAt = &now
	return t
}

func (t *WalletTransactionValues) MarkAsFailed(reason string) *WalletTransactionValues {
	t.SetStatus(TransactionStatusFailed)
	now := utils.TimeNowMilli()
	t.FailedAt = &now
	if reason != "" {
		t.AdminNotes = &reason
	}
	return t
}

func (t *WalletTransactionValues) Cancel(reason string) *WalletTransactionValues {
	t.SetStatus(TransactionStatusCancelled)
	if reason != "" {
		t.AdminNotes = &reason
	}
	return t
}

func (t *WalletTransactionValues) Reverse(reason string) *WalletTransactionValues {
	t.SetStatus(TransactionStatusReversed)
	if reason != "" {
		t.AdminNotes = &reason
	}
	return t
}

// 审核相关方法
func (t *WalletTransactionValues) SubmitForReview(reason string) *WalletTransactionValues {
	t.ReviewStatus = utils.StringPtr(ReviewStatusPendingReview)
	if reason != "" {
		t.AdminNotes = &reason
	}
	return t
}

func (t *WalletTransactionValues) ApproveManually(reviewerID string) *WalletTransactionValues {
	t.ReviewStatus = utils.StringPtr(ReviewStatusManualApproved)
	t.ReviewedBy = &reviewerID
	now := utils.TimeNowMilli()
	t.ReviewedAt = &now
	return t
}

func (t *WalletTransactionValues) RejectReview(reviewerID, reason string) *WalletTransactionValues {
	t.ReviewStatus = utils.StringPtr(ReviewStatusRejected)
	t.ReviewedBy = &reviewerID
	now := utils.TimeNowMilli()
	t.ReviewedAt = &now
	if reason != "" {
		t.AdminNotes = &reason
	}
	return t
}

// 创建特定类型的交易
func NewRideEarningTransaction(accountID, userID, rideOrderID string, amount float64) *WalletTransaction {
	tx := NewWalletTransactionV2()
	tx.SetAccountID(accountID).
		SetUserID(userID).
		SetType(TransactionTypeIncome).
		SetCategory(TransactionCategoryRideEarning).
		SetAmount(amount).
		SetTitle("行程收入").
		SetDescription(fmt.Sprintf("完成行程获得收入: %.2f RWF", amount)).
		SetRelated("ride_order", rideOrderID)

	return tx
}

func NewRidePaymentTransaction(accountID, userID, rideOrderID string, amount float64) *WalletTransaction {
	tx := NewWalletTransactionV2()
	tx.SetAccountID(accountID).
		SetUserID(userID).
		SetType(TransactionTypeExpense).
		SetCategory(TransactionCategoryRidePayment).
		SetAmount(amount).
		SetTitle("行程支付").
		SetDescription(fmt.Sprintf("支付行程费用: %.2f RWF", amount)).
		SetRelated("ride_order", rideOrderID)

	return tx
}

func NewTopupTransaction(accountID, userID, paymentID string, amount float64) *WalletTransaction {
	tx := NewWalletTransactionV2()
	tx.SetAccountID(accountID).
		SetUserID(userID).
		SetType(TransactionTypeIncome).
		SetCategory(TransactionCategoryTopup).
		SetAmount(amount).
		SetTitle("钱包充值").
		SetDescription(fmt.Sprintf("钱包充值: %.2f RWF", amount)).
		SetPaymentID(paymentID)

	return tx
}

func NewWithdrawalTransaction(accountID, userID, withdrawalID string, amount, fee float64) *WalletTransaction {
	tx := NewWalletTransactionV2()
	tx.SetAccountID(accountID).
		SetUserID(userID).
		SetType(TransactionTypeWithdrawal).
		SetCategory(TransactionCategoryWithdrawalFee).
		SetAmount(amount).
		SetWithdrawalID(withdrawalID).
		SetTitle("提现申请").
		SetDescription(fmt.Sprintf("申请提现: %.2f RWF (手续费: %.2f RWF)", amount, fee))

	if fee > 0 {
		tx.FeeAmount = &fee
	}

	return tx
}
