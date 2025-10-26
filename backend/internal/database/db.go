package database

import (
	"fmt"
	"log"

	"multiagent-chat/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 数据库全局变量
var DB *gorm.DB

// InitDatabase 初始化数据库
func InitDatabase() error {
	// MySQL连接配置
	dsn := "root:root@tcp(127.0.0.1:3306)/ai_agents?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci&interpolateParams=true&readTimeout=10s&writeTimeout=10s"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置数据库编码为UTF-8
	if err := DB.Exec("SET NAMES utf8mb4").Error; err != nil {
		return fmt.Errorf("设置数据库编码失败: %v", err)
	}

	// 自动迁移表结构
	// 注意：需要在主程序中导入模型包以确保模型被正确注册
	err = DB.AutoMigrate(&model.Agent{}, &model.ChatRoom{}, &model.Message{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库初始化成功（MySQL模式）")
	return nil
}
