package ws

import (
	"encoding/json"
	"log"
	"myapp/config"
	"myapp/models"
	"sort"
	"strconv"
	"strings"
)

type RoomConfig struct {
	SingleCount int
	MultiCount  int
	JudgeCount  int
}

type Room struct {
	ID          string
	Hub         *Hub
	Clients     map[*Client]bool
	Broadcast   chan Message
	Status      string
	Config      RoomConfig
	Questions   []models.Question
	Submissions map[uint]*Submission
}

type Submission struct {
	Score   int
	Correct int
	Wrong   int
	Answers map[uint][]string
}

func NewRoom(id string, hub *Hub) *Room {
	return &Room{
		ID:          id,
		Hub:         hub,
		Clients:     make(map[*Client]bool),
		Broadcast:   make(chan Message),
		Status:      "waiting",
		Submissions: make(map[uint]*Submission),
	}
}

func (r *Room) Run() {
	for {
		select {
		case msg := <-r.Broadcast:
			bytes, _ := json.Marshal(msg)
			for client := range r.Clients {
				select {
				case client.Send <- bytes:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

func (r *Room) Join(client *Client) {
	r.Clients[client] = true
	client.RoomID = r.ID
	client.IsReady = false // 刚进来是未准备

	log.Printf("User %d joined room %s", client.UserID, r.ID)

	// 1. 广播有人加入
	r.Broadcast <- Message{
		Type: "user_joined",
		Payload: map[string]interface{}{
			"user_id": client.UserID,
			"count":   len(r.Clients),
		},
	}

	// 2. 同步当前的准备状态 (让新来的人知道谁已经准备了)
	r.BroadcastReadyStatus()
}

func (r *Room) Leave(client *Client) {
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
		client.RoomID = ""
		client.IsReady = false

		r.Broadcast <- Message{
			Type: "user_left",
			Payload: map[string]interface{}{
				"user_id": client.UserID,
				"count":   len(r.Clients),
			},
		}
		r.BroadcastReadyStatus()
	}
}

func (r *Room) HandleMessage(sender *Client, msg Message) {
	switch msg.Type {
	case "ready":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if status, ok := payload["status"].(bool); ok {
				sender.IsReady = status
				log.Printf("User %d ready status: %v", sender.UserID, status)
				r.BroadcastReadyStatus()
				r.CheckAndAutoStart()
			}
		}

	case "start_game":
		r.StartGame()

	case "progress_update":
		// 转发进度，务必带上 UserID
		r.Broadcast <- Message{
			Type:    "progress_update",
			UserID:  sender.UserID,
			Payload: msg.Payload,
		}

	case "submit_exam":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if answersRaw, ok := payload["answers"].(map[string]interface{}); ok {
				answers := make(map[uint][]string)
				for k, v := range answersRaw {
					qID, _ := strconv.Atoi(k)
					if arr, ok := v.([]interface{}); ok {
						strs := make([]string, len(arr))
						for i, s := range arr {
							strs[i] = s.(string)
						}
						answers[uint(qID)] = strs
					}
				}
				r.HandleSubmit(sender, answers)
			}
		}
	}
}

func (r *Room) CheckAndAutoStart() {
	// 至少要有2个人才能开始
	if len(r.Clients) < 2 {
		return
	}
	
	allReady := true
	for client := range r.Clients {
		if !client.IsReady {
			allReady = false
			break
		}
	}
	
	if allReady {
		log.Printf("Room %s: All players ready, auto starting game...", r.ID)
		// 稍微延迟一点点体验更好（可选），这里直接调用
		r.StartGame()
	}
}

func (r *Room) BroadcastReadyStatus() {
	// map key JSON只能是string
	statusMap := make(map[string]bool)
	for c := range r.Clients {
		statusMap[strconv.Itoa(int(c.UserID))] = c.IsReady
	}
	r.Broadcast <- Message{
		Type:    "ready_update",
		Payload: statusMap,
	}
}

func (r *Room) StartGame() {
	if r.Status == "playing" {
		return
	}
	
	sCount := r.Config.SingleCount
	mCount := r.Config.MultiCount
	jCount := r.Config.JudgeCount
	if sCount == 0 && mCount == 0 && jCount == 0 {
		sCount = 10
	}

	var questions []models.Question
	var singles, multis, judges []models.Question

	if sCount > 0 {
		config.DB.Where("question_type = ?", "single_choice").Order("RAND()").Limit(sCount).Find(&singles)
		questions = append(questions, singles...)
	}
	if mCount > 0 {
		config.DB.Where("question_type = ?", "multiple_choice").Order("RAND()").Limit(mCount).Find(&multis)
		questions = append(questions, multis...)
	}
	if jCount > 0 {
		config.DB.Where("question_type = ?", "judgment").Order("RAND()").Limit(jCount).Find(&judges)
		questions = append(questions, judges...)
	}

	r.Questions = questions
	r.Status = "playing"
	r.Submissions = make(map[uint]*Submission) // 重置提交

	r.Broadcast <- Message{
		Type: "game_start",
		Payload: map[string]interface{}{
			"questions": questions,
		},
	}
}

func (r *Room) HandleSubmit(client *Client, answers map[uint][]string) {
	score := 0
	correct := 0
	wrong := 0

	for _, q := range r.Questions {
		userAns, ok := answers[q.ID]
		if !ok {
			wrong++
			continue
		}
		userSorted := make([]string, len(userAns))
		copy(userSorted, userAns)
		sort.Strings(userSorted)

		correctSorted := make([]string, len(q.Answer))
		copy(correctSorted, q.Answer)
		sort.Strings(correctSorted)

		if strings.Join(userSorted, ",") == strings.Join(correctSorted, ",") {
			correct++
			if q.QuestionType == "multiple_choice" {
				score += 4
			} else {
				score += 2
			}
		} else {
			wrong++
		}
	}

	r.Submissions[client.UserID] = &Submission{
		Score:   score,
		Correct: correct,
		Wrong:   wrong,
		Answers: answers,
	}

	r.Broadcast <- Message{
		Type: "user_finished",
		Payload: map[string]interface{}{
			"user_id": client.UserID,
		},
	}
	r.CheckGameOver()
}

func (r *Room) CheckGameOver() {
	// 所有人交卷后结束
	if len(r.Submissions) >= len(r.Clients) && len(r.Clients) > 0 {
		r.Status = "finished"
		report := make(map[uint]*Submission)
		for uid, sub := range r.Submissions {
			report[uid] = sub
		}
		r.Broadcast <- Message{
			Type: "game_over",
			Payload: map[string]interface{}{
				"report": report,
			},
		}
	}
}