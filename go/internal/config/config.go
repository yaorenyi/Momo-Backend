package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 结构体对应你的配置文件字段
type Config struct {
	Port            int    `yaml:"PORT"`
	AllowOrigin     string `yaml:"ALLOW_ORIGIN"`
	ResendAPIKey    string `yaml:"RESEND_API_KEY"`
	ResendFromEmail string `yaml:"RESEND_FROM_EMAIL"`
	EmailAddress    string `yaml:"EMAIL_ADDRESS"`
	AdminName       string `yaml:"ADMIN_NAME"`
	AdminPassword   string `yaml:"ADMIN_PASSWORD"`
}

var GlobalConfig *Config

// DefaultConfig 默认配置内容
func DefaultConfig() *Config {
	return &Config{
		Port:            17171,
		AllowOrigin:     "http://localhost:4321,https://example.com",
		ResendAPIKey:    "re_xxxxxx",
		ResendFromEmail: "<notify@notifications.example.com>",
		EmailAddress:    "me@example.com",
		AdminName:       "admin",
		AdminPassword:   "password",
	}
}

// LoadConfig 加载或初始化配置文件
func LoadConfig() (*Config, error) {
	configPath := "./config/config.yaml"

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// 如果文件不存在，创建并写入默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultCfg := DefaultConfig()
		data, err := yaml.Marshal(defaultCfg)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, err
		}
		GlobalConfig = defaultCfg
		return defaultCfg, nil
	}

	// 读取现有文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	GlobalConfig = &cfg
	return &cfg, nil
}
