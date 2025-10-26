package main

import (
	"fmt"
	"log"

	"multiagent-chat/internal/database"
	"multiagent-chat/internal/router"
)

func main() {
	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置路由
	routerInstance := router.SetupRouter(database.DB)

	fmt.Println("服务器启动在 http://localhost:8080")
	log.Fatal(routerInstance.Run(":8080"))
}
