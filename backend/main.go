package main

import (
	"fmt"
	"myapp/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("your_super_secret_key") // 生产环境请放入环境变量

// 初始化数据库
func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ Failed to connect to database")
	}
	// 🚀 修改这里：添加 &models.Question{} 到自动迁移列表
	db.AutoMigrate(&models.User{}, &models.Question{})
	fmt.Println("✅ Database migrated: Users & Questions")
}

// 注册接口
func register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	// 默认角色为 user
	if input.Role == "" {
		input.Role = "user"
	}

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// 登录接口
func login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"role":  user.Role,
	})
}

// 创建题目
func createQuestion(c *gin.Context) {
	var input models.Question
	// 绑定 JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 写入数据库
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}
	c.JSON(http.StatusOK, input)
}

// 获取题目列表 (支持 ?category=xx&type=xx 筛选)
func getQuestions(c *gin.Context) {
	var questions []models.Question
	
	category := c.Query("category")
	qType := c.Query("type")
	
	query := db.Model(&models.Question{})
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if qType != "" {
		query = query.Where("question_type = ?", qType)
	}

	// 按创建时间倒序，限制返回 100 条（生产环境建议分页）
	query.Order("created_at desc").Limit(100).Find(&questions)
	
	c.JSON(http.StatusOK, gin.H{"data": questions, "total": len(questions)})
}

// 更新题目
func updateQuestion(c *gin.Context) {
	id := c.Param("id")
	var question models.Question
	
	if err := db.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var input models.Question
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Updates 会自动处理 JSON 序列化更新
	db.Model(&question).Updates(input)
	c.JSON(http.StatusOK, question)
}

// 删除题目
func deleteQuestion(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&models.Question{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func main() {
	initDB()
	r := gin.Default()
	
	// 允许跨域 (简单版)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/register", register)
	r.POST("/login", login)

// 🚀 新增：注册 /api/questions 路由组
	api := r.Group("/api")
	{
		api.POST("/questions", createQuestion)
		api.GET("/questions", getQuestions)
		api.PUT("/questions/:id", updateQuestion)
		api.DELETE("/questions/:id", deleteQuestion)
	}

	r.Run(":8080")
}