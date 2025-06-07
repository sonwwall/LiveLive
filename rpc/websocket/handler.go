package main

import (
	dao "LiveLive/dao/rdb"
	websocket "LiveLive/kitex_gen/livelive/websocket"
	"LiveLive/rpc/websocket/ws"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// WebsocketServiceImpl implements the last service interface defined in the IDL.
type WebsocketServiceImpl struct {
	wsHub *ws.WsHub
}

// BroadcastToCourse implements the WebsocketServiceImpl interface.
func (s *WebsocketServiceImpl) BroadcastToCourse(ctx context.Context, req *websocket.BroadcastToCourseReq) (resp *websocket.BroadcastToCourseResp, err error) {
	s.wsHub.BroadcastToCourse(req.CourseId, req.Data)
	resp = &websocket.BroadcastToCourseResp{}
	return resp, nil
}

// AggregateAnswers implements the WebsocketServiceImpl interface.
func (s *WebsocketServiceImpl) AggregateAnswers(ctx context.Context, req *websocket.AggregateAnswersReq) (resp *websocket.AggregateAnswersResp, err error) {
	//从redis读出保存的答案
	key := fmt.Sprintf("choice_answer:%d", req.QuestionId)
	data, err := dao.Redis.HGetAll(context.Background(), key).Result()
	log.Println(data)
	if err != nil {
		log.Println("读取 Redis 答案失败:", err.Error())
		return &websocket.AggregateAnswersResp{
			Accuracy: 0,
		}, nil
	}

	// 统计答案分布
	count := map[string]int{}
	sum := 0
	var correctAnswer string
	for _, answer := range data {
		count[answer]++
		sum++
	}
	switch req.CorrectAnswer {
	case 0:
		correctAnswer = "A"
	case 1:
		correctAnswer = "B"
	case 2:
		correctAnswer = "C"
	case 3:
		correctAnswer = "D"
	default:
		correctAnswer = ""
	}
	//统计正确率
	accuracy := (float64(count[correctAnswer]) / float64(sum)) * 100

	// 构造结果消息
	resultMsg := map[string]interface{}{
		"type": "answer_result",
		"data": map[string]interface{}{
			"question_id": req.QuestionId,
			"summary":     count,
			"accuracy":    fmt.Sprintf("%.2f%s", accuracy, "%"),
		},
	}
	payload, _ := json.Marshal(resultMsg)

	// 推送给当前课程的老师
	for client := range s.wsHub.Connections[req.CourseId] {
		if client.Role == 0 {
			client.SendCh <- payload
			break
		}
	}
	return &websocket.AggregateAnswersResp{
		Accuracy: accuracy,
	}, nil
}

// CountRegister implements the WebsocketServiceImpl interface.
func (s *WebsocketServiceImpl) CountRegister(ctx context.Context, req *websocket.CountRegisterReq) (resp *websocket.CountRegisterResp, err error) {
	//从redis读
	key := fmt.Sprintf("sign:course:%d", req.CourseId)
	data, err := dao.Redis.HGetAll(context.Background(), key).Result()
	log.Println(data)
	if err != nil {
		log.Println("读取 Redis 答案失败:", err.Error())
		return &websocket.CountRegisterResp{}, nil
	}
	total := len(data)
	signedIn := 0
	unsignedFields := []string{}

	for field, val := range data {
		if val == "1" {
			signedIn++
		} else {
			unsignedFields = append(unsignedFields, field)
		}
	}
	resultMsg := map[string]interface{}{
		"type": "register_result",
		"data": map[string]interface{}{
			"total":      total,
			"已签到人数": signedIn,
			"未签到人数": total - signedIn,
			"未签到学生": unsignedFields,
		},
	}
	payload, _ := json.Marshal(resultMsg)

	//推送给老师
	for client := range s.wsHub.Connections[req.CourseId] {
		if client.Role == 0 {
			client.SendCh <- payload
			break
		}
	}

	return
}

// AggregateTrueOrFalseAnswers implements the WebsocketServiceImpl interface.
func (s *WebsocketServiceImpl) AggregateTrueOrFalseAnswers(ctx context.Context, req *websocket.AggregateTrueOrFalseAnswersReq) (resp *websocket.AggregateTrueOrFalseAnswersResp, err error) {
	//从redis读出保存的答案
	key := fmt.Sprintf("true_or_false_answer:%d", req.QuestionId)
	data, err := dao.Redis.HGetAll(context.Background(), key).Result()
	log.Println(data)
	if err != nil {
		log.Println("读取 Redis 答案失败:", err.Error())
		return &websocket.AggregateTrueOrFalseAnswersResp{
			Accuracy: 0,
		}, nil
	}

	// 统计答案分布
	count := map[string]int{}
	sum := 0
	var correctAnswer string
	for _, answer := range data {
		count[answer]++
		sum++
	}

	switch req.CorrectAnswer {
	case 0:
		correctAnswer = "false"
	case 1:
		correctAnswer = "true"
	default:
		correctAnswer = ""
	}
	//统计正确率
	accuracy := (float64(count[correctAnswer]) / float64(sum)) * 100

	// 构造结果消息
	resultMsg := map[string]interface{}{
		"type": "answer_result",
		"data": map[string]interface{}{
			"question_id": req.QuestionId,
			"summary":     count,
			"accuracy":    fmt.Sprintf("%.2f%s", accuracy, "%"),
		},
	}
	payload, _ := json.Marshal(resultMsg)

	// 推送给当前课程的老师
	for client := range s.wsHub.Connections[req.CourseId] {
		if client.Role == 0 {
			client.SendCh <- payload
			break
		}
	}
	return &websocket.AggregateTrueOrFalseAnswersResp{
		Accuracy: accuracy,
	}, nil
}
