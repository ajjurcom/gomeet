package common

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// Encryption 使用sha265算法加密密码
func Encryption(password string) string {
	h := sha256.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
