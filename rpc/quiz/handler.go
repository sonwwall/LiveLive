package main

import (
	"LiveLive/dao/db"
	"LiveLive/kitex_gen/livelive/base"
	quiz "LiveLive/kitex_gen/livelive/quiz"
	"LiveLive/model"
	"LiveLive/utils"
	"LiveLive/ws"
	"context"
	"encoding/json"
	"gorm.io/datatypes"
	"log"
)

// QuizServiceImpl implements the last service interface defined in the IDL.
type QuizServiceImpl struct {
	wsHub *ws.WsHub
}

//func NewQuizService(hub *ws.WsHub) *QuizServiceImpl {
//	return &QuizServiceImpl{wsHub: hub}
//}

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
	err = db.AddChoiceQuestion(question)
	if err != nil {
		log.Printf("db.AddChoiceQuestion err: %+v", err)
	}

	msg := map[string]interface{}{
		"type": "choice_question",
		"data": map[string]interface{}{
			"title":    req.Title,
			"options":  req.Options,
			"deadline": req.Deadline,
		},
	}
	data, _ := json.Marshal(msg)
	s.wsHub.BroadcastToCourse(req.CourseId, data)

	res := &quiz.PublishChoiceQuestionResp{
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
	}
	return res, nil

}
