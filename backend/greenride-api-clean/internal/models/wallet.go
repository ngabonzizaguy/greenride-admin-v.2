package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// Wallet 钱包表 - 基于最新设计文档
type Wallet struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	WalletID string `json:"wallet_id" gorm:"column:wallet_id;type:varchar(64);uniqueIndex"`
	Salt     string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*WalletValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type WalletValues struct {
	UserID   *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index;default:'user'"` // user, driver

	// 余额信息
	Balance  *float64 `json:"balance" gorm:"column:balance;type:decimal(12,2);default:0"`
	Currency *string  `json:"currency" gorm:"column:currency;type:varchar(3);default:'USD'"`

	// 冻结金额 (用于订单预扣等)
	FrozenAmount *float64 `json:"frozen_amount" gorm:"column:frozen_amount;type:decimal(12,2);default:0"`

	// 状态信息
	Status   *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"` // active, inactive, suspended, frozen
	IsActive *bool   `json:"is_active" gorm:"column:is_active;default:true"`

	// 限额设置
	DailyLimit   *float64 `json:"daily_limit" gorm:"column:daily_limit;type:decimal(10,2)"`
	MonthlyLimit *float64 `json:"monthly_limit" gorm:"column:monthly_limit;type:decimal(10,2)"`

	// 统计信息
	TotalEarnings  *float64 `json:"total_earnings" gorm:"column:total_earnings;type:decimal(12,2);default:0"`   // 总收入
	TotalSpending  *float64 `json:"total_spending" gorm:"column:total_spending;type:decimal(12,2);default:0"`   // 总支出
	TotalWithdrawn *float64 `json:"total_withdrawn" gorm:"column:total_withdrawn;type:decimal(12,2);default:0"` // 总提现
	TotalDeposited *float64 `json:"total_deposited" gorm:"column:total_deposited;type:decimal(12,2);default:0"` // 总充值

	// 时间戳
	LastTransactionAt *int64 `json:"last_transaction_at" gorm:"column:last_transaction_at"`
	LastDepositAt     *int64 `json:"last_deposit_at" gorm:"column:last_deposit_at"`
	LastWithdrawAt    *int64 `json:"last_withdraw_at" gorm:"column:last_withdraw_at"`

	// 验证状态
	IsVerified *bool  `json:"is_verified" gorm:"column:is_verified;default:false"`
	VerifiedAt *int64 `json:"verified_at" gorm:"column:verified_at"`

	// 扩展信息
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // JSON格式的额外信息

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Wallet) TableName() string {
	return "t_wallets"
}

// 钱包状态常量
const (
	WalletStatusActive    = "active"
	WalletStatusInactive  = "inactive"
	WalletStatusSuspended = "suspended"
	WalletStatusFrozen    = "frozen"
)

// 创建新的钱包对象
func NewWalletV2() *Wallet {
	return &Wallet{
		WalletID: utils.GenerateWalletID(),
		Salt:     utils.GenerateSalt(),
		WalletValues: &WalletValues{
			UserType:       utils.StringPtr(protocol.UserTypePassenger),
			Balance:        utils.Float64Ptr(0),
			Currency:       utils.StringPtr("USD"),
			FrozenAmount:   utils.Float64Ptr(0),
			Status:         utils.StringPtr(WalletStatusActive),
			IsActive:       utils.BoolPtr(true),
			IsVerified:     utils.BoolPtr(false),
			TotalEarnings:  utils.Float64Ptr(0),
			TotalSpending:  utils.Float64Ptr(0),
			TotalWithdrawn: utils.Float64Ptr(0),
			TotalDeposited: utils.Float64Ptr(0),
		},
	}
}

// SetValues 更新WalletV2Values中的非nil值
func (w *WalletValues) SetValues(values *WalletValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		w.UserID = values.UserID
	}
	if values.UserType != nil {
		w.UserType = values.UserType
	}
	if values.Balance != nil {
		w.Balance = values.Balance
	}
	if values.Currency != nil {
		w.Currency = values.Currency
	}
	if values.FrozenAmount != nil {
		w.FrozenAmount = values.FrozenAmount
	}
	if values.Status != nil {
		w.Status = values.Status
	}
	if values.IsActive != nil {
		w.IsActive = values.IsActive
	}
	if values.DailyLimit != nil {
		w.DailyLimit = values.DailyLimit
	}
	if values.MonthlyLimit != nil {
		w.MonthlyLimit = values.MonthlyLimit
	}
	if values.IsVerified != nil {
		w.IsVerified = values.IsVerified
	}
	if values.Metadata != nil {
		w.Metadata = values.Metadata
	}
	if values.UpdatedAt > 0 {
		w.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (w *WalletValues) GetBalance() float64 {
	if w.Balance == nil {
		return 0
	}
	return *w.Balance
}

func (w *WalletValues) GetFrozenAmount() float64 {
	if w.FrozenAmount == nil {
		return 0
	}
	return *w.FrozenAmount
}

func (w *WalletValues) GetAvailableBalance() float64 {
	return w.GetBalance() - w.GetFrozenAmount()
}

func (w *WalletValues) GetStatus() string {
	if w.Status == nil {
		return WalletStatusActive
	}
	return *w.Status
}

func (w *WalletValues) GetIsActive() bool {
	if w.IsActive == nil {
		return true
	}
	return *w.IsActive
}

func (w *WalletValues) GetIsVerified() bool {
	if w.IsVerified == nil {
		return false
	}
	return *w.IsVerified
}

func (w *WalletValues) GetUserID() string {
	if w.UserID == nil {
		return ""
	}
	return *w.UserID
}

func (w *WalletValues) GetCurrency() string {
	if w.Currency == nil {
		return "USD"
	}
	return *w.Currency
}

func (w *WalletValues) GetTotalEarnings() float64 {
	if w.TotalEarnings == nil {
		return 0
	}
	return *w.TotalEarnings
}

func (w *WalletValues) GetTotalSpending() float64 {
	if w.TotalSpending == nil {
		return 0
	}
	return *w.TotalSpending
}

// Setter 方法
func (w *WalletValues) SetBalance(balance float64) *WalletValues {
	w.Balance = &balance
	return w
}

func (w *WalletValues) SetStatus(status string) *WalletValues {
	w.Status = &status
	return w
}

func (w *WalletValues) SetUserID(userID string) *WalletValues {
	w.UserID = &userID
	return w
}

func (w *WalletValues) SetLastTransactionAt(timestamp int64) *WalletValues {
	w.LastTransactionAt = &timestamp
	return w
}

// 业务方法
func (w *Wallet) CanTransact() bool {
	return w.GetIsActive() && w.GetStatus() == WalletStatusActive
}

func (w *Wallet) HasSufficientBalance(amount float64) bool {
	return w.GetAvailableBalance() >= amount
}

func (w *Wallet) IsActive() bool {
	return w.GetIsActive() && w.GetStatus() == WalletStatusActive
}

func (w *Wallet) IsFrozen() bool {
	return w.GetStatus() == WalletStatusFrozen
}

// 余额操作方法
func (w *WalletValues) AddBalance(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	newBalance := w.GetBalance() + amount
	w.SetBalance(newBalance)
	now := utils.TimeNowMilli()
	w.SetLastTransactionAt(now)

	// 更新总充值金额
	totalDeposited := w.GetTotalDeposited() + amount
	w.TotalDeposited = &totalDeposited
	w.LastDepositAt = &now

	return nil
}

func (w *WalletValues) SubtractBalance(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if !w.HasSufficientBalance(amount) {
		return fmt.Errorf("insufficient balance")
	}

	newBalance := w.GetBalance() - amount
	w.SetBalance(newBalance)
	now := utils.TimeNowMilli()
	w.SetLastTransactionAt(now)

	// 更新总支出金额
	totalSpending := w.GetTotalSpending() + amount
	w.TotalSpending = &totalSpending

	return nil
}

func (w *WalletValues) FreezeAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if w.GetAvailableBalance() < amount {
		return fmt.Errorf("insufficient available balance")
	}

	newFrozenAmount := w.GetFrozenAmount() + amount
	w.FrozenAmount = &newFrozenAmount

	return nil
}

func (w *WalletValues) UnfreezeAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if w.GetFrozenAmount() < amount {
		return fmt.Errorf("insufficient frozen amount")
	}

	newFrozenAmount := w.GetFrozenAmount() - amount
	w.FrozenAmount = &newFrozenAmount

	return nil
}

func (w *WalletValues) TransferFrozenToBalance(amount float64) error {
	if err := w.UnfreezeAmount(amount); err != nil {
		return err
	}

	// 这里不需要再调用AddBalance，因为冻结的钱已经在余额中
	// 只需要从冻结金额中减去即可
	now := utils.TimeNowMilli()
	w.SetLastTransactionAt(now)

	return nil
}

func (w *WalletValues) Withdraw(amount float64) error {
	if err := w.SubtractBalance(amount); err != nil {
		return err
	}

	// 更新总提现金额
	totalWithdrawn := w.GetTotalWithdrawn() + amount
	w.TotalWithdrawn = &totalWithdrawn
	now := utils.TimeNowMilli()
	w.LastWithdrawAt = &now

	return nil
}

func (w *WalletValues) AddEarnings(amount float64) error {
	if err := w.AddBalance(amount); err != nil {
		return err
	}

	// 更新总收入
	totalEarnings := w.GetTotalEarnings() + amount
	w.TotalEarnings = &totalEarnings

	return nil
}

func (w *WalletValues) MarkAsVerified() {
	verified := true
	w.IsVerified = &verified
	now := utils.TimeNowMilli()
	w.VerifiedAt = &now
}

func (w *WalletValues) GetTotalDeposited() float64 {
	if w.TotalDeposited == nil {
		return 0
	}
	return *w.TotalDeposited
}

func (w *WalletValues) GetTotalWithdrawn() float64 {
	if w.TotalWithdrawn == nil {
		return 0
	}
	return *w.TotalWithdrawn
}

func (w *WalletValues) HasSufficientBalance(amount float64) bool {
	return w.GetAvailableBalance() >= amount
}
