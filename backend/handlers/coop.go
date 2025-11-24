package handlers

import (
	"myapp/ws"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 获取房间列表
func GetRooms(c *gin.Context) {
	if ws.GlobalHub == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hub not initialized"})
		return
	}
	list := ws.GlobalHub.GetRoomList()
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// 创建房间 (支持参数配置)
func CreateRoom(c *gin.Context) {
	// 尝试获取 UserID (从 query 参数获取，如果没有则默认为 0)
	userIDStr := c.Query("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	// 1. 解析 JSON 参数
	var input struct {
		Single int `json:"single"`
		Multi  int `json:"multi"`
		Judge  int `json:"judge"`
	}
	
	// 绑定参数 (如果绑定失败，例如用户没传，也不应该直接报错，而是用默认值)
	if err := c.ShouldBindJSON(&input); err != nil {
		// 给一个默认值
		input.Single = 10
		input.Multi = 0
		input.Judge = 0
	}

	// 2. 生成房间 ID
	roomID := uuid.New().String()[:8]
	
	// 3. 创建房间 (注意：CreateRoom 需要两个参数: id, hostID)
	room := ws.GlobalHub.CreateRoom(roomID, uint(userID))
	
	// 4. 🚀 关键：把配置存入 Room，供 StartGame 使用
	// 注意：需要确保 Room 结构体有 Config 字段，且 Config 字段是公开的
	room.Config = ws.RoomConfig{
		SingleCount: input.Single,
		MultiCount:  input.Multi,
		JudgeCount:  input.Judge,
	}
	
	c.JSON(http.StatusOK, gin.H{
		"room_id": roomID,
		"message": "Room created successfully",
		"config":  room.Config,
	})
}