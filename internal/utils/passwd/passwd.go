package passwd

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12

// 接收明文密码，返回加密后的哈希值
func Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("密码不能为空")
	}

	// 使用bcrypt生成密码哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	return string(hashedBytes), nil
}

// 比较明文密码和存储的哈希值
func Verify(password, hashedPassword string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}

	if hashedPassword == "" {
		return fmt.Errorf("哈希密码不能为空")
	}

	// 使用bcrypt比较密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("密码验证失败: %w", err)
	}

	return nil
}
