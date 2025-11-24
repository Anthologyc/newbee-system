package handlers

import (
	"myapp/config"
	"myapp/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// 获取分类
func GetCategories(c *gin.Context) {
	var categories []string
	if err := config.DB.Model(&models.Question{}).Distinct("category").Pluck("category", &categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// 随机题目 (含筛选)
func GetRandomQuestions(c *gin.Context) {
	countStr := c.DefaultQuery("count", "100")
	category := c.Query("category")
	qType := c.Query("type")
	count, _ := strconv.Atoi(countStr)

	var questions []models.Question
	query := config.DB.Model(&models.Question{})
	if category != "" { query = query.Where("category = ?", category) }
	if qType != "" { query = query.Where("question_type = ?", qType) }
	
	query.Order("RAND()").Limit(count).Find(&questions)
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

// 顺序题目
func GetSequentialQuestions(c *gin.Context) {
	var questions []models.Question
	config.DB.Order("id asc").Find(&questions)
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

// CRUD 操作
func CreateQuestion(c *gin.Context) {
	var input models.Question
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}

func GetQuestions(c *gin.Context) {
	var questions []models.Question
	config.DB.Order("created_at desc").Limit(100).Find(&questions)
	c.JSON(http.StatusOK, gin.H{"data": questions, "total": len(questions)})
}

func GetQuestion(c *gin.Context) {
	var q models.Question
	if err := config.DB.First(&q, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, q)
}

func UpdateQuestion(c *gin.Context) {
	var q models.Question
	if err := config.DB.First(&q, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	var input models.Question
	c.ShouldBindJSON(&input)
	config.DB.Model(&q).Updates(input)
	c.JSON(http.StatusOK, q)
}

func DeleteQuestion(c *gin.Context) {
	config.DB.Delete(&models.Question{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}