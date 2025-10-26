package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 应用配置结构
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	
	Database struct {
		Type         string `yaml:"type"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DBName       string `yaml:"dbname"`
		Charset      string `yaml:"charset"`
		Timezone     string `yaml:"timezone"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	} `yaml:"database"`
	
	API struct {
		APIKey   string `yaml:"api_key"`
		BaseURL  string `yaml:"base_url"`
		Model    string `yaml:"model"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"api"`
	
	Log struct {
		Level   string `yaml:"level"`
		File    string `yaml:"file"`
		Console bool   `yaml:"console"`
	} `yaml:"log"`
}

// AppConfig 全局配置实例
var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 如果没有指定配置文件路径，使用默认路径
	if configPath == "" {
		configPath = "config.yaml"
	}
	
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}
	
	// 解析YAML
	AppConfig = &Config{}
	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}
	
	return nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s&collation=utf8mb4_unicode_ci&interpolateParams=true&readTimeout=10s&writeTimeout=10s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.Charset,
		c.Database.Timezone,
	)
}

// GetServerAddr 获取服务器地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// 兼容性函数，保持向后兼容
func GetAPIKey() string {
	if AppConfig != nil {
		return AppConfig.API.APIKey
	}
	return "sk-2b32164c932c4a6e9e5ce6a519ccecfb" // 默认值
}

func GetBaseURL() string {
	if AppConfig != nil {
		return AppConfig.API.BaseURL
	}
	return "https://dashscope.aliyuncs.com/compatible-mode/v1/" // 默认值
}

func GetModel() string {
	if AppConfig != nil {
		return AppConfig.API.Model
	}
	return "qwen-plus" // 默认值
}

func GetEndpoint() string {
	if AppConfig != nil {
		return AppConfig.API.Endpoint
	}
	return "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation" // 默认值
}
