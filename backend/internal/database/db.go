package database

import (
	"fmt"
	"log"

	"multiagent-chat/internal/model"
	"multiagent-chat/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 数据库全局变量
var DB *gorm.DB

// InitDatabase 初始化数据库
func InitDatabase() error {
	// 使用配置文件中的MySQL连接配置
	dsn := config.AppConfig.GetDSN()

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	// 设置连接池参数（从配置文件读取）
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)

	// 设置数据库编码为UTF-8
	if err := DB.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
		return fmt.Errorf("设置数据库编码失败: %v", err)
	}

	// 显式设置客户端字符集
	if err := DB.Exec("SET character_set_client = utf8mb4").Error; err != nil {
		return fmt.Errorf("设置客户端字符集失败: %v", err)
	}
	if err := DB.Exec("SET character_set_results = utf8mb4").Error; err != nil {
		return fmt.Errorf("设置结果字符集失败: %v", err)
	}
	if err := DB.Exec("SET character_set_connection = utf8mb4").Error; err != nil {
		return fmt.Errorf("设置连接字符集失败: %v", err)
	}

	// 自动迁移表结构
	// 注意：需要在主程序中导入模型包以确保模型被正确注册
	err = DB.AutoMigrate(&model.Agent{}, &model.ChatRoom{}, &model.Message{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Printf("数据库初始化成功（MySQL模式）- 连接到: %s:%d/%s", 
		config.AppConfig.Database.Host, 
		config.AppConfig.Database.Port, 
		config.AppConfig.Database.DBName)
	return nil
}
