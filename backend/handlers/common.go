package handlers

import (
	"myapp/config"
	"myapp/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

// 公告相关
func CreateAnnouncement(c *gin.Context) {
	var input models.Announcement
	c.ShouldBindJSON(&input)
	input.Status = true
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}

func GetAnnouncements(c *gin.Context) {
	var list []models.Announcement
	config.DB.Order("created_at desc").Limit(50).Find(&list)
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func DeleteAnnouncement(c *gin.Context) {
	config.DB.Delete(&models.Announcement{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// 仪表盘
func GetDashboardStats(c *gin.Context) {
	var userCount, qCount int64
	var lastAnno models.Announcement
	config.DB.Model(&models.User{}).Count(&userCount)
	config.DB.Model(&models.Question{}).Count(&qCount)
	config.DB.Order("created_at desc").First(&lastAnno)
	c.JSON(http.StatusOK, gin.H{"user_count": userCount, "question_count": qCount, "announcement": lastAnno})
}

// 用户管理
func GetUsers(c *gin.Context) {
	var users []models.User
	config.DB.Order("created_at desc").Limit(100).Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func DeleteUser(c *gin.Context) {
	config.DB.Delete(&models.User{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}