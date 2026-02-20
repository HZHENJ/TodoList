package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRandomString 生成指定长度的随机 Hex 字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// GenerateDefaultUsername 生成默认用户名
func GenerateDefaultUsername() string {
	return fmt.Sprintf("USER_%s", GenerateRandomString(16))
}
