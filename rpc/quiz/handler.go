package main

import (
	"LiveLive/dao/db"
	dao "LiveLive/dao/rdb"
	"LiveLive/kitex_gen/livelive/base"
	quiz "LiveLive/kitex_gen/livelive/quiz"
	"LiveLive/model"
	"LiveLive/utils"
	"LiveLive/ws"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"log"
	"time"
)

// QuizServiceImpl implements the last service interface defined in the IDL.
type QuizServiceImpl struct {
	wsHub *ws.WsHub
}

// PublishChoiceQuestion implements the QuizServiceImpl interface.
func (s *QuizServiceImpl) PublishChoiceQuestion(ctx context.Context, req *quiz.PublishChoiceQuestionReq) (resp *quiz.PublishChoiceQuestionResp, err error) {
	optsJson, _ := json.Marshal(req.Options) //四个选项先转化为json类型
	question := &model.ChoiceQuestion{
		CourseID:  req.CourseId,
		TeacherId: req.TeacherId,
		Title:     req.Title,
		Options:   datatypes.JSON(optsJson),
		Answer:    req.Answer,
		Deadline:  utils.TimestampToPtr(req.Deadline),
	}
	err, questionId := db.AddChoiceQuestion(question)
	if err != nil {
		log.Printf("db.AddChoiceQuestion err: %+v", err)
	}

	msg := map[string]interface{}{
		"type": "choice_question",
		"data": map[string]interface{}{
			"question_id": questionId,
			"title":       req.Title,
			"options":     req.Options,
			"deadline":    req.Deadline,
		},
	}
	data, _ := json.Marshal(msg)
	s.wsHub.BroadcastToCourse(req.CourseId, data)

	//设置定时任务,到时间自动返回统计结果
	time.AfterFunc(time.Unix(req.Deadline, 0).Sub(time.Now()), func() {
		AggregateAnswers(int64(questionId), req.CourseId, s.wsHub)
	})

	res := &quiz.PublishChoiceQuestionResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil

}

// AggregateAnswers 答案统计与推送逻辑
func AggregateAnswers(questionID, courseID int64, hub *ws.WsHub) {
	key := fmt.Sprintf("answer:%d", questionID)
	data, err := dao.Redis.HGetAll(context.Background(), key).Result()
	log.Println(data)
	if err != nil {
		log.Println("读取 Redis 答案失败:", err)
		return
	}

	// 统计答案分布
	count := map[string]int{}
	for _, answer := range data {
		count[answer]++
	}
	log.Println(count)

	// 构造结果消息
	resultMsg := map[string]interface{}{
		"type": "answer_result",
		"data": map[string]interface{}{
			"question_id": questionID,
			"summary":     count,
		},
	}
	payload, _ := json.Marshal(resultMsg)

	// 推送给当前课程下所有老师
	for client := range hub.Connections[courseID] {
		if client.Role == 0 {
			client.SendCh <- payload
		}
	}
}
