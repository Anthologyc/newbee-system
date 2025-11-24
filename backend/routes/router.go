package routes

import (
	"myapp/handlers"
	"myapp/middleware"
	"github.com/gin-gonic/gin"
	"myapp/ws"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors()) // 使用抽离的中间件

	// 认证
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	api := r.Group("/api")
	{
		// 题目
		api.POST("/questions", handlers.CreateQuestion)
		api.GET("/questions", handlers.GetQuestions)
		api.GET("/questions/:id", handlers.GetQuestion)
		api.PUT("/questions/:id", handlers.UpdateQuestion)
		api.DELETE("/questions/:id", handlers.DeleteQuestion)
		api.GET("/categories", handlers.GetCategories) // 新增

		// 练习模式
		api.GET("/practice/random", handlers.GetRandomQuestions)
		api.GET("/practice/sequential", handlers.GetSequentialQuestions)
		
		// 考试与错题
		api.POST("/exam/submit", handlers.SubmitExam)
		api.GET("/mistakes", handlers.GetMistakes)
		api.POST("/mistakes/remove", handlers.RemoveMistake)
		api.POST("/mistakes/record", handlers.RecordMistake)
		api.GET("/exam/generate", handlers.GenerateExamPaper)
		// 公告
		api.POST("/announcements", handlers.CreateAnnouncement)
		api.GET("/announcements", handlers.GetAnnouncements)
		api.DELETE("/announcements/:id", handlers.DeleteAnnouncement)

		// 管理
		api.GET("/dashboard/stats", handlers.GetDashboardStats)
		api.GET("/users", handlers.GetUsers)
		api.DELETE("/users/:id", handlers.DeleteUser)

		ws.Init()

		r.GET("/ws", func(c *gin.Context) {
        ws.ServeWs(ws.GlobalHub, c)

		api.GET("/rooms", handlers.GetRooms)
		api.POST("/rooms", handlers.CreateRoom)
    })

	}

	return r
}