package handlers

import (
	"myapp/config"
	"myapp/models"
	"net/http"
	"sort"
	"strings"
	"strconv"
	"github.com/gin-gonic/gin"
)

// 🚀 新增：单题错题录入 (给随机/顺序模式用)
func RecordMistake(c *gin.Context) {
	var input struct {
		UserID     uint `json:"user_id"`
		QuestionID uint `json:"question_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	config.DB.Model(&models.Mistake{}).Where("user_id = ? AND question_id = ?", input.UserID, input.QuestionID).Count(&count)
	if count == 0 {
		config.DB.Create(&models.Mistake{UserID: input.UserID, QuestionID: input.QuestionID})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recorded"})
}

// 提交试卷 (支持自定义分值)
func SubmitExam(c *gin.Context) {
	var input struct {
		UserID      uint                `json:"user_id"`
		Answers     map[uint][]string   `json:"answers"`
		// 🚀 新增：接收分值配置
		SingleScore int `json:"single_score"`
		MultiScore  int `json:"multi_score"`
		JudgeScore  int `json:"judge_score"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认分值 (防空值)
	if input.SingleScore == 0 { input.SingleScore = 2 }
	if input.MultiScore == 0 { input.MultiScore = 4 }
	if input.JudgeScore == 0 { input.JudgeScore = 2 }

	var totalScore int = 0
	var correctCount int = 0
	var wrongCount int = 0

	type ResultItem struct {
		QuestionID    uint     `json:"question_id"`
		IsCorrect     bool     `json:"is_correct"`
		UserAnswer    []string `json:"user_answer"`
		CorrectAnswer []string `json:"correct_answer"`
	}
	var results []ResultItem

	for qID, userAns := range input.Answers {
		var q models.Question
		if err := config.DB.First(&q, qID).Error; err != nil { continue }

		isCorrect := false
		questionScore := 0

		// 🚀 使用自定义分值
		switch q.QuestionType {
		case "single_choice":
			questionScore = input.SingleScore
		case "multiple_choice":
			questionScore = input.MultiScore
		case "judgment":
			questionScore = input.JudgeScore
		default:
			questionScore = 2
		}

		userSorted := make([]string, len(userAns))
		copy(userSorted, userAns)
		sort.Strings(userSorted)

		correctSorted := make([]string, len(q.Answer))
		copy(correctSorted, q.Answer)
		sort.Strings(correctSorted)

		if strings.Join(userSorted, ",") == strings.Join(correctSorted, ",") {
			isCorrect = true
		}

		if isCorrect {
			correctCount++
			totalScore += questionScore
		} else {
			wrongCount++
			// 错题录入
			var count int64
			config.DB.Model(&models.Mistake{}).Where("user_id = ? AND question_id = ?", input.UserID, q.ID).Count(&count)
			if count == 0 {
				config.DB.Create(&models.Mistake{UserID: input.UserID, QuestionID: q.ID})
			}
		}

		results = append(results, ResultItem{
			QuestionID:    q.ID,
			IsCorrect:     isCorrect,
			UserAnswer:    userAns,
			CorrectAnswer: q.Answer,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"score":   totalScore,
		"correct": correctCount,
		"wrong":   wrongCount,
		"details": results,
	})
}

// GetMistakes, RemoveMistake 保持不变...
// ... 把之前的代码贴在这里 ...
func GetMistakes(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "")
	var mistakes []models.Mistake
	
	if err := config.DB.Preload("Question").Where("user_id = ?", userID).Find(&mistakes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mistakes"})
		return
	}

	var questions []models.Question
	for _, m := range mistakes {
		if m.Question.ID != 0 {
			questions = append(questions, m.Question)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

func RemoveMistake(c *gin.Context) {
	var input struct {
		UserID     uint `json:"user_id"`
		QuestionID uint `json:"question_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Where("user_id = ? AND question_id = ?", input.UserID, input.QuestionID).Delete(&models.Mistake{})
	c.JSON(http.StatusOK, gin.H{"message": "Removed"})
}

func GenerateExamPaper(c *gin.Context) {
	// 获取参数 (这里用 Query 参数比较方便 GET 请求)
	singleCount, _ := strconv.Atoi(c.DefaultQuery("single", "0"))
	multiCount, _ := strconv.Atoi(c.DefaultQuery("multi", "0"))
	judgeCount, _ := strconv.Atoi(c.DefaultQuery("judge", "0"))

	var questions []models.Question
	var singles, multis, judges []models.Question

	// 分别查询
	if singleCount > 0 {
		config.DB.Where("question_type = ?", "single_choice").Order("RAND()").Limit(singleCount).Find(&singles)
		questions = append(questions, singles...)
	}
	if multiCount > 0 {
		config.DB.Where("question_type = ?", "multiple_choice").Order("RAND()").Limit(multiCount).Find(&multis)
		questions = append(questions, multis...)
	}
	if judgeCount > 0 {
		config.DB.Where("question_type = ?", "judgment").Order("RAND()").Limit(judgeCount).Find(&judges)
		questions = append(questions, judges...)
	}

	c.JSON(http.StatusOK, gin.H{"data": questions})
}