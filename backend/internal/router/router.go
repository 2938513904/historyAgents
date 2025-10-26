package router

import (
	"net/http"
	"strings"

	"multiagent-chat/internal/controller"
	"multiagent-chat/internal/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 设置UTF-8编码中间件
	router.Use(func(c *gin.Context) {
		// 只为API路由设置JSON内容类型
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
		c.Next()
	})

	// 配置静态文件服务
	router.Static("/static", "../frontend/static")
	router.LoadHTMLFiles("../frontend/index.html", "../frontend/test-connection.html", "../frontend/dashscope-test.html")

	// 设置根路径路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/test-connection.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test-connection.html", nil)
	})

	router.GET("/dashscope-test.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashscope-test.html", nil)
	})

	// WebSocket路由
	wsManager := websocket.NewManager()
	wsManager.Start()

	// 创建控制器实例
	agentController := controller.NewAgentController(db)
	chatRoomController := controller.NewChatRoomController(db, wsManager)
	apiController := controller.NewAPIController()

	// API路由
	api := router.Group("/api")
	{
		// 智能体管理
		api.GET("/agents", agentController.GetAgents)
		api.POST("/agents", agentController.CreateAgent)
		api.PUT("/agents/:id", agentController.UpdateAgent)
		api.DELETE("/agents/:id", agentController.DeleteAgent)
		api.DELETE("/agents-empty", agentController.DeleteEmptyAgents)

		// 聊天室管理
		api.GET("/chatrooms", chatRoomController.GetChatRooms)
		api.POST("/chatrooms", chatRoomController.CreateChatRoom)
		api.GET("/chatrooms/:id", chatRoomController.GetChatRoom)
		api.POST("/chatrooms/:id/start", chatRoomController.StartChatRoom)
		api.POST("/chatrooms/:id/stop", chatRoomController.StopChatRoom)
		api.DELETE("/chatrooms/:id", chatRoomController.DeleteChatRoom)
		api.DELETE("/chatrooms-empty", chatRoomController.DeleteEmptyChatRooms)

		// 连接测试（使用硬编码配置）
		api.POST("/test-connection", apiController.TestAPIConnection)
		api.POST("/test-dashscope", apiController.TestDashScopeAPI)
	}

	router.GET("/api/ws/:roomId", wsManager.HandleWebSocket)

	return router
}
