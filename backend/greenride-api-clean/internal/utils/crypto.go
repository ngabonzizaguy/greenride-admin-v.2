package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// encrypt encrypts plaintext using the given key with AES.
func Encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts ciphertext using the given key with AES.
func Decrypt(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

// HashPassword 使用bcrypt和salt哈希密码
// 为了避免bcrypt的72字节限制，先使用SHA256预哈希密码+salt
func HashPassword(password, salt string) (string, error) {
	// 先与salt结合
	saltedPassword := password + salt

	// 使用SHA256预哈希以避免bcrypt的72字节限制
	// 这样既保持了安全性，又确保输入长度固定为32字节
	preHash := sha256.Sum256([]byte(saltedPassword))
	preHashHex := hex.EncodeToString(preHash[:])

	// 使用bcrypt哈希预哈希结果
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(preHashHex), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
// 使用与HashPassword相同的预哈希逻辑
func VerifyPassword(password, salt, hashedPassword string) bool {
	// 先与salt结合
	saltedPassword := password + salt

	// 使用SHA256预哈希以保持与HashPassword的一致性
	preHash := sha256.Sum256([]byte(saltedPassword))
	preHashHex := hex.EncodeToString(preHash[:])

	// 使用bcrypt验证预哈希结果
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(preHashHex))
	return err == nil
}

// GenerateSecureHash 生成安全哈希（用于API密钥等）
func GenerateSecureHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// HashString 生成字符串的简单数值哈希（用于一致性分配）
func HashString(input string) int64 {
	hash := sha256.Sum256([]byte(input))
	// 取前8个字节转换为int64
	result := int64(0)
	for i := 0; i < 8 && i < len(hash); i++ {
		result = result*256 + int64(hash[i])
	}
	if result < 0 {
		result = -result
	}
	return result
}

// GenerateAPIKey 生成API密钥
func GenerateAPIKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备用
		return fmt.Sprintf("api_%d_%s", TimeNowMilli(), GenerateShortID())
	}
	return hex.EncodeToString(bytes)
}

// GenerateSecretKey 生成密钥
func GenerateSecretKey() string {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("secret_%d_%s", TimeNowMilli(), GenerateUUID())
	}
	return hex.EncodeToString(bytes)
}
