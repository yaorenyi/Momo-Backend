package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// GetCravatar 生成 Cravatar/Gravatar 头像地址
func GetCravatar(email string) string {
	// 1. 去除首尾空格并转为小写
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	// 2. 计算 MD5 哈希
	hash := md5.Sum([]byte(cleanEmail))
	hashHex := hex.EncodeToString(hash[:])

	// 3. 返回国内镜像地址 (Cravatar)
	return "https://cravatar.cn/avatar/" + hashHex + "?s=200&d=retro"
}
