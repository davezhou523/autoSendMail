package helper

import (
	"crypto/sha256"
	"encoding/base64"
)

func GenerateToken(params, secret string) string {
	data := params + secret
	// 计算SHA-256哈希
	hash := sha256.Sum256([]byte(data))
	// 将哈希结果进行Base64编码
	token := base64.URLEncoding.EncodeToString(hash[:])
	return token
}
