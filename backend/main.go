package main

import (
	"log"
	"net/http"

	"multiagent-chat/internal/config"
	"multiagent-chat/internal/database"
	"multiagent-chat/internal/router"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	log.Println("配置文件加载成功")

	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置路由
	r := router.SetupRouter(database.DB)

	// 启动服务器
	serverAddr := config.AppConfig.GetServerAddr()
	log.Printf("服务器启动在: http://%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
