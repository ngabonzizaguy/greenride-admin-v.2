package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// UserAccount 用户账户表 - 管理用户钱包余额
type UserAccount struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:varchar(64);uniqueIndex"`
	Salt      string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*UserAccountValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type UserAccountValues struct {
	// 用户关联
	UserID   *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index"` // passenger, driver

	// 账户基本信息
	Currency *string `json:"currency" gorm:"column:currency;type:varchar(10);default:'RWF'"`
	Status   *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"` // active, suspended, frozen, closed

	// 余额信息
	AvailableBalance *float64 `json:"available_balance" gorm:"column:available_balance;type:decimal(15,2);default:0.00"` // 可用余额
	FrozenBalance    *float64 `json:"frozen_balance" gorm:"column:frozen_balance;type:decimal(15,2);default:0.00"`       // 冻结余额
	PendingBalance   *float64 `json:"pending_balance" gorm:"column:pending_balance;type:decimal(15,2);default:0.00"`     // 待结算余额

	// 版本控制(用于并发控制)
	Version             *int64 `json:"version" gorm:"column:version;default:1"`
	LastBalanceUpdateAt *int64 `json:"last_balance_update_at" gorm:"column:last_balance_update_at"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (UserAccount) TableName() string {
	return "t_user_accounts"
}

// 账户状态常量
const (
	AccountStatusActive    = "active"
	AccountStatusSuspended = "suspended"
	AccountStatusFrozen    = "frozen"
	AccountStatusClosed    = "closed"
)

// 创建新的用户账户对象
func NewUserAccountV2() *UserAccount {
	return &UserAccount{
		AccountID: utils.GenerateUserAccountID(),
		Salt:      utils.GenerateSalt(),
		UserAccountValues: &UserAccountValues{
			UserType:         utils.StringPtr(protocol.UserTypePassenger),
			Currency:         utils.StringPtr("RWF"),
			Status:           utils.StringPtr(AccountStatusActive),
			AvailableBalance: utils.Float64Ptr(0.00),
			FrozenBalance:    utils.Float64Ptr(0.00),
			PendingBalance:   utils.Float64Ptr(0.00),
			Version:          utils.Int64Ptr(1),
		},
	}
}

// SetValues 更新UserAccountV2Values中的非nil值
func (a *UserAccountValues) SetValues(values *UserAccountValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		a.UserID = values.UserID
	}
	if values.UserType != nil {
		a.UserType = values.UserType
	}
	if values.Currency != nil {
		a.Currency = values.Currency
	}
	if values.Status != nil {
		a.Status = values.Status
	}
	if values.AvailableBalance != nil {
		a.AvailableBalance = values.AvailableBalance
	}
	if values.FrozenBalance != nil {
		a.FrozenBalance = values.FrozenBalance
	}
	if values.PendingBalance != nil {
		a.PendingBalance = values.PendingBalance
	}
	if values.UpdatedAt > 0 {
		a.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (a *UserAccountValues) GetUserID() string {
	if a.UserID == nil {
		return ""
	}
	return *a.UserID
}

func (a *UserAccountValues) GetUserType() string {
	if a.UserType == nil {
		return protocol.UserTypePassenger
	}
	return *a.UserType
}

func (a *UserAccountValues) GetCurrency() string {
	if a.Currency == nil {
		return "RWF"
	}
	return *a.Currency
}

func (a *UserAccountValues) GetStatus() string {
	if a.Status == nil {
		return AccountStatusActive
	}
	return *a.Status
}

func (a *UserAccountValues) GetAvailableBalance() float64 {
	if a.AvailableBalance == nil {
		return 0.00
	}
	return *a.AvailableBalance
}

func (a *UserAccountValues) GetFrozenBalance() float64 {
	if a.FrozenBalance == nil {
		return 0.00
	}
	return *a.FrozenBalance
}

func (a *UserAccountValues) GetPendingBalance() float64 {
	if a.PendingBalance == nil {
		return 0.00
	}
	return *a.PendingBalance
}

func (a *UserAccountValues) GetTotalBalance() float64 {
	return a.GetAvailableBalance() + a.GetFrozenBalance() + a.GetPendingBalance()
}

func (a *UserAccountValues) GetVersion() int64 {
	if a.Version == nil {
		return 1
	}
	return *a.Version
}

// Setter 方法
func (a *UserAccountValues) SetUserID(userID string) *UserAccountValues {
	a.UserID = &userID
	return a
}

func (a *UserAccountValues) SetUserType(userType string) *UserAccountValues {
	a.UserType = &userType
	return a
}

func (a *UserAccountValues) SetCurrency(currency string) *UserAccountValues {
	a.Currency = &currency
	return a
}

func (a *UserAccountValues) SetStatus(status string) *UserAccountValues {
	a.Status = &status
	return a
}

func (a *UserAccountValues) SetAvailableBalance(balance float64) *UserAccountValues {
	a.AvailableBalance = &balance
	a.updateLastBalanceTime()
	return a
}

func (a *UserAccountValues) SetFrozenBalance(balance float64) *UserAccountValues {
	a.FrozenBalance = &balance
	a.updateLastBalanceTime()
	return a
}

func (a *UserAccountValues) SetPendingBalance(balance float64) *UserAccountValues {
	a.PendingBalance = &balance
	a.updateLastBalanceTime()
	return a
}

func (a *UserAccountValues) updateLastBalanceTime() {
	now := utils.TimeNowMilli()
	a.LastBalanceUpdateAt = &now
}

func (a *UserAccountValues) IncrementVersion() *UserAccountValues {
	version := a.GetVersion() + 1
	a.Version = &version
	return a
}

// 业务方法
func (a *UserAccount) IsActive() bool {
	return a.GetStatus() == AccountStatusActive
}

func (a *UserAccount) IsFrozen() bool {
	return a.GetStatus() == AccountStatusFrozen
}

func (a *UserAccount) IsSuspended() bool {
	return a.GetStatus() == AccountStatusSuspended
}

func (a *UserAccount) IsClosed() bool {
	return a.GetStatus() == AccountStatusClosed
}

func (a *UserAccount) CanTransact() bool {
	return a.IsActive()
}

func (a *UserAccount) HasSufficientBalance(amount float64) bool {
	return a.GetAvailableBalance() >= amount
}

// 余额操作方法
func (a *UserAccountValues) AddAvailableBalance(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	newBalance := a.GetAvailableBalance() + amount
	a.SetAvailableBalance(newBalance)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) SubtractAvailableBalance(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if !a.HasSufficientBalance(amount) {
		return fmt.Errorf("insufficient available balance")
	}

	newBalance := a.GetAvailableBalance() - amount
	a.SetAvailableBalance(newBalance)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) FreezeAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if !a.HasSufficientBalance(amount) {
		return fmt.Errorf("insufficient available balance to freeze")
	}

	// 从可用余额转移到冻结余额
	newAvailable := a.GetAvailableBalance() - amount
	newFrozen := a.GetFrozenBalance() + amount

	a.SetAvailableBalance(newAvailable)
	a.SetFrozenBalance(newFrozen)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) UnfreezeAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if a.GetFrozenBalance() < amount {
		return fmt.Errorf("insufficient frozen balance")
	}

	// 从冻结余额转移到可用余额
	newFrozen := a.GetFrozenBalance() - amount
	newAvailable := a.GetAvailableBalance() + amount

	a.SetFrozenBalance(newFrozen)
	a.SetAvailableBalance(newAvailable)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) MoveToPending(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if !a.HasSufficientBalance(amount) {
		return fmt.Errorf("insufficient available balance")
	}

	// 从可用余额转移到待结算余额
	newAvailable := a.GetAvailableBalance() - amount
	newPending := a.GetPendingBalance() + amount

	a.SetAvailableBalance(newAvailable)
	a.SetPendingBalance(newPending)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) MoveFromPending(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if a.GetPendingBalance() < amount {
		return fmt.Errorf("insufficient pending balance")
	}

	// 从待结算余额转移到可用余额
	newPending := a.GetPendingBalance() - amount
	newAvailable := a.GetAvailableBalance() + amount

	a.SetPendingBalance(newPending)
	a.SetAvailableBalance(newAvailable)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) DeductFromFrozen(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if a.GetFrozenBalance() < amount {
		return fmt.Errorf("insufficient frozen balance")
	}

	newFrozen := a.GetFrozenBalance() - amount
	a.SetFrozenBalance(newFrozen)
	a.IncrementVersion()

	return nil
}

func (a *UserAccountValues) DeductFromPending(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if a.GetPendingBalance() < amount {
		return fmt.Errorf("insufficient pending balance")
	}

	newPending := a.GetPendingBalance() - amount
	a.SetPendingBalance(newPending)
	a.IncrementVersion()

	return nil
}

// 账户状态管理
func (a *UserAccountValues) Suspend(reason string) *UserAccountValues {
	a.SetStatus(AccountStatusSuspended)
	return a
}

func (a *UserAccountValues) Freeze(reason string) *UserAccountValues {
	a.SetStatus(AccountStatusFrozen)
	return a
}

func (a *UserAccountValues) Activate() *UserAccountValues {
	a.SetStatus(AccountStatusActive)
	return a
}

func (a *UserAccountValues) Close() *UserAccountValues {
	a.SetStatus(AccountStatusClosed)
	return a
}

// 检查是否有足够余额（包含可用余额检查）
func (a *UserAccountValues) HasSufficientBalance(amount float64) bool {
	return a.GetAvailableBalance() >= amount
}

// 获取账户摘要信息
func (a *UserAccountValues) GetAccountSummary() map[string]interface{} {
	return map[string]interface{}{
		"available_balance": a.GetAvailableBalance(),
		"frozen_balance":    a.GetFrozenBalance(),
		"pending_balance":   a.GetPendingBalance(),
		"total_balance":     a.GetTotalBalance(),
		"currency":          a.GetCurrency(),
		"status":            a.GetStatus(),
		"version":           a.GetVersion(),
	}
}
