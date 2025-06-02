package main

import (
	"LiveLive/dao/db"
	dao "LiveLive/dao/rdb"
	"LiveLive/kitex_gen/livelive/base"
	quiz "LiveLive/kitex_gen/livelive/quiz"
	"LiveLive/kitex_gen/livelive/websocket"
	"LiveLive/kitex_gen/livelive/websocket/websocketservice"
	"LiveLive/model"
	"LiveLive/utils"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"log"
	"time"
)

// QuizServiceImpl implements the last service interface defined in the IDL.
type QuizServiceImpl struct {
	wsClient websocketservice.Client
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
	//广播给所有学生
	s.wsClient.BroadcastToCourse(ctx, &websocket.BroadcastToCourseReq{
		CourseId: req.CourseId,
		Data:     data,
	})

	//设置一个过期时间，统计完后自动过期
	duration := 3 * time.Minute
	log.Println(duration)
	key := fmt.Sprintf("choice_answer:%d", questionId)

	//设置定时任务,到时间自动返回统计结果
	time.AfterFunc(time.Unix(req.Deadline, 0).Sub(time.Now()), func() {
		result, _ := s.wsClient.AggregateAnswers(ctx, &websocket.AggregateAnswersReq{
			QuestionId:    int64(questionId),
			CourseId:      req.CourseId,
			CorrectAnswer: req.Answer,
		})
		dao.Redis.Expire(ctx, key, duration) //统计完三分钟后自动清除redis缓存
		answeredChoiceQuestion := &model.AnsweredChoiceQuestion{
			ChoiceQuestionId: questionId,
			Title:            req.Title,
			Options:          datatypes.JSON(optsJson),
			Answer:           req.Answer,
			Accuracy:         result.Accuracy,
		}

		_ = db.AddAnsweredChoiceQuestion(answeredChoiceQuestion)
	})

	res := &quiz.PublishChoiceQuestionResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil

}
