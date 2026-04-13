package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GetCravatar 生成 Cravatar/Gravatar 头像地址
func GetCravatar(email string) string {
	// 1. 处理 Email
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	hash := md5.Sum([]byte(cleanEmail))
	hashHex := hex.EncodeToString(hash[:])

	// Cravatar 原始地址
	cravatarURL := fmt.Sprintf("https://cravatar.cn/avatar/%s?s=200&d=retro", hashHex)

	// 2. 创建一个带有超时的客户端，防止因网络问题阻塞程序
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 3. 发起 HEAD 请求检查响应头
	resp, err := client.Head(cravatarURL)
	if err != nil {
		// 如果请求失败（如网络抖动），降级返回原始 cravatar 地址
		return cravatarURL
	}
	defer resp.Body.Close()

	// 4. 判断是否为默认头像
	// Cravatar 默认头像会在 Header 中包含 avatar-from: default
	if resp.Header.Get("avatar-from") == "default" {
		return fmt.Sprintf("https://avatar.motues.top/avatar?variant=beam&name=%s&colors=FFADAD,FFD6A5,FDFFB6,FF9900,AABBCC", hashHex)
	}

	return cravatarURL
}
