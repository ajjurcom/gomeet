package common

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// Encryption 使用sha265算法加密密码
//
// 参数：password 密码
//
// 返回值: string 使用sha265算法加密后的密码，长读64位
func Encryption(password string) string {
	h := sha256.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
