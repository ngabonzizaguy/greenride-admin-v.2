package services

import (
	"crypto/tls"
	"fmt"
	"greenride/internal/config"
	"net/smtp"
	"sort"
	"strings"
	"sync"
	"time"
)

type EmailAccountService struct {
	accounts     []EmailAccount
	failedCounts map[string]int            // 记录账户失败次数
	lastFailed   map[string]time.Time      // 记录最后失败时间
	sendCounts   map[string][]time.Time    // 记录发送时间戳（用于频率限制）
	dailyCounts  map[string]map[string]int // 记录每日发送次数 [accountName][date] = count
	mutex        sync.RWMutex
}

type EmailAccount struct {
	Config config.EmailAccountConfig
	Auth   smtp.Auth
}

// NewEmailAccountService 创建邮箱账户服务
func NewEmailAccountService(configs []config.EmailAccountConfig) *EmailAccountService {
	service := &EmailAccountService{
		accounts:     make([]EmailAccount, 0),
		failedCounts: make(map[string]int),
		lastFailed:   make(map[string]time.Time),
		sendCounts:   make(map[string][]time.Time),
		dailyCounts:  make(map[string]map[string]int),
	}

	for _, cfg := range configs {
		if cfg.Enabled {
			account := EmailAccount{
				Config: cfg,
				Auth:   smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host),
			}
			service.accounts = append(service.accounts, account)
		}
	}

	// 按优先级排序
	sort.Slice(service.accounts, func(i, j int) bool {
		return service.accounts[i].Config.Priority < service.accounts[j].Config.Priority
	})

	return service
}

func (s *EmailAccountService) ServiceName() string {
	return "account"
}

// SendEmail 发送邮件，自动选择可用的邮箱账户
func (s *EmailAccountService) SendEmail(to string, subject string, body string) error {
	s.mutex.RLock()
	availableAccounts := s.getAvailableAccounts()
	s.mutex.RUnlock()

	if len(availableAccounts) == 0 {
		return fmt.Errorf("no available email accounts")
	}

	var lastErr error
	for _, account := range availableAccounts {
		// 检查发送频率和日限制
		if !s.canSendEmail(account.Config.Name) {
			continue
		}

		err := s.sendEmailWithAccount(account, to, subject, body)
		if err == nil {
			// 发送成功，记录发送时间并重置失败计数
			s.mutex.Lock()
			s.recordSentEmail(account.Config.Name)
			delete(s.failedCounts, account.Config.Name)
			delete(s.lastFailed, account.Config.Name)
			s.mutex.Unlock()
			return nil
		}

		// 记录失败
		s.mutex.Lock()
		s.failedCounts[account.Config.Name]++
		s.lastFailed[account.Config.Name] = time.Now()
		s.mutex.Unlock()

		lastErr = err
	}

	return fmt.Errorf("all email accounts failed, last error: %v", lastErr)
}

// sendEmailWithAccount 使用指定账户发送邮件
func (s *EmailAccountService) sendEmailWithAccount(account EmailAccount, to string, subject string, body string) error {
	// 构建邮件内容
	message := s.buildMessage(account.Config.From, to, subject, body)

	// SMTP 服务器地址
	addr := fmt.Sprintf("%s:%d", account.Config.Host, account.Config.Port)

	if account.Config.SSL {
		return s.sendEmailWithTLS(addr, account.Auth, account.Config.From, []string{to}, message)
	} else {
		return smtp.SendMail(addr, account.Auth, account.Config.From, []string{to}, message)
	}
}

// sendEmailWithTLS 使用TLS发送邮件
func (s *EmailAccountService) sendEmailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 创建TLS配置
	host := strings.Split(addr, ":")[0]
	tlsConfig := &tls.Config{
		ServerName: host,
	}

	// 建立TLS连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return err
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(msg)
	return err
}

// buildMessage 构建邮件消息
func (s *EmailAccountService) buildMessage(from, to, subject, body string) []byte {
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	return []byte(message)
}

// getAvailableAccounts 获取可用的邮箱账户（排除最近失败太多次的）
func (s *EmailAccountService) getAvailableAccounts() []EmailAccount {
	now := time.Now()
	available := make([]EmailAccount, 0)

	for _, account := range s.accounts {
		name := account.Config.Name

		// 检查失败次数和时间
		failCount := s.failedCounts[name]
		if failCount >= 3 {
			// 如果失败超过3次，需要等待一段时间
			if lastFail, exists := s.lastFailed[name]; exists {
				// 等待时间随失败次数递增：5分钟、15分钟、30分钟
				waitMinutes := 5 * failCount
				if waitMinutes > 30 {
					waitMinutes = 30
				}
				if now.Sub(lastFail) < time.Duration(waitMinutes)*time.Minute {
					continue // 跳过这个账户
				}
			}
		}

		available = append(available, account)
	}

	return available
}

// GetAccountStatus 获取账户状态信息
func (s *EmailAccountService) GetAccountStatus() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	now := time.Now()
	today := now.Format("2006-01-02")

	status := make(map[string]interface{})
	for _, account := range s.accounts {
		name := account.Config.Name

		// 统计最近100秒内的发送次数
		recentSendCount := 0
		if sendTimes, exists := s.sendCounts[name]; exists {
			cutoff := now.Add(-100 * time.Second)
			for _, sendTime := range sendTimes {
				if sendTime.After(cutoff) {
					recentSendCount++
				}
			}
		}

		// 统计今日发送次数
		todaySendCount := 0
		if s.dailyCounts[name] != nil {
			todaySendCount = s.dailyCounts[name][today]
		}

		accountStatus := map[string]interface{}{
			"name":            name,
			"host":            account.Config.Host,
			"from":            account.Config.From,
			"priority":        account.Config.Priority,
			"enabled":         account.Config.Enabled,
			"failCount":       s.failedCounts[name],
			"recentSendCount": recentSendCount,      // 最近100秒发送次数
			"todaySendCount":  todaySendCount,       // 今日发送次数
			"canSend":         s.canSendEmail(name), // 当前是否可以发送
		}

		if lastFail, exists := s.lastFailed[name]; exists {
			accountStatus["lastFailed"] = lastFail.Format("2006-01-02 15:04:05")
		}

		status[name] = accountStatus
	}

	return status
}

// ResetAccountStatus 重置指定账户的状态
func (s *EmailAccountService) ResetAccountStatus(accountName string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.failedCounts, accountName)
	delete(s.lastFailed, accountName)
	delete(s.sendCounts, accountName)
	if s.dailyCounts[accountName] != nil {
		delete(s.dailyCounts, accountName)
	}
}

// canSendEmail 检查是否可以发送邮件（频率和日限制检查）
func (s *EmailAccountService) canSendEmail(accountName string) bool {
	now := time.Now()

	// 检查频率限制：200 emails/100 seconds
	if sendTimes, exists := s.sendCounts[accountName]; exists {
		// 清除100秒前的记录
		var recentSends []time.Time
		cutoff := now.Add(-100 * time.Second)
		for _, sendTime := range sendTimes {
			if sendTime.After(cutoff) {
				recentSends = append(recentSends, sendTime)
			}
		}
		s.sendCounts[accountName] = recentSends

		// 检查是否超过频率限制
		if len(recentSends) >= 200 {
			return false
		}
	}

	// 检查日限制：100 emails/day
	today := now.Format("2006-01-02")
	if s.dailyCounts[accountName] == nil {
		s.dailyCounts[accountName] = make(map[string]int)
	}

	if s.dailyCounts[accountName][today] >= 100 {
		return false
	}

	return true
}

// recordSentEmail 记录发送的邮件
func (s *EmailAccountService) recordSentEmail(accountName string) {
	now := time.Now()

	// 记录发送时间（用于频率限制）
	if s.sendCounts[accountName] == nil {
		s.sendCounts[accountName] = make([]time.Time, 0)
	}
	s.sendCounts[accountName] = append(s.sendCounts[accountName], now)

	// 记录日发送次数
	today := now.Format("2006-01-02")
	if s.dailyCounts[accountName] == nil {
		s.dailyCounts[accountName] = make(map[string]int)
	}
	s.dailyCounts[accountName][today]++
}
