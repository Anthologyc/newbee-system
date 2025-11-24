package main

import (
	"myapp/config"
	"myapp/routes"
)

func main() {
	// 1. 初始化数据库
	config.InitDB()

	// 2. 设置路由
	r := routes.SetupRouter()

	// 3. 启动服务
	r.Run(":8080")
}