package handlers

import (
	"LiveLive/api/code"
	"LiveLive/api/rpc"
	"LiveLive/kitex_gen/livelive/quiz"
	"LiveLive/middleware"
	"LiveLive/model"
	"LiveLive/utils/response"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"time"
)

func PublishChoiceQuestion(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}

	type PublishChoiceQuestionReq struct {
		CourseID int64    `json:"course_id,required" form:"course_id,required"` //由前端直接传入
		Title    string   `json:"title,required" form:"title,required"`
		Options  []string `json:"options,required" form:"options,required"`
		Answer   int8     `json:"answer,required" form:"answer,required"`
		Deadline uint64   `json:"deadline,required" form:"deadline,required"` //秒
	}

	var req PublishChoiceQuestionReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误：" + err.Error(),
		})
		return
	}

	result, _ := rpc.PublishChoiceQuestion(ctx, &quiz.PublishChoiceQuestionReq{
		TeacherId: int64(user.(*model.User).ID),
		CourseId:  req.CourseID,
		Title:     req.Title,
		Options:   req.Options,
		Answer:    req.Answer,
		Deadline:  time.Now().Add(time.Duration(req.Deadline) * time.Second).Unix(),
	})

	if result == nil {
		c.JSON(200, response.Response{
			Code: -1,
			Msg:  "内部错误",
		})
		return
	}
	if result.BaseResp.Code != 0 {
		c.JSON(200, response.Response{
			Code: result.BaseResp.Code,
			Msg:  result.BaseResp.Msg,
		})
	}
	c.JSON(200, response.Response{
		Code: 0,
		Msg:  "ok",
	})
}

func PublishTrueOrFalseQuestion(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}
	type PublishTrueOrFalseQuestionReq struct {
		CourseID int64  `json:"course_id,required" form:"course_id,required"` //由前端直接传入
		Title    string `json:"title,required" form:"title,required"`
		Answer   int8   `json:"answer,required" form:"answer,required"`
		Deadline uint64 `json:"deadline,required" form:"deadline,required"` //秒
	}

	var req PublishTrueOrFalseQuestionReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误:" + err.Error(),
		})
		return
	}

	result, _ := rpc.PublishTrueOrFalseQuestion(ctx, &quiz.PublishTrueOrFalseQuestionReq{
		CourseId:  req.CourseID,
		Title:     req.Title,
		TeacherId: int64(user.(*model.User).ID),
		Deadline:  time.Now().Add(time.Duration(req.Deadline) * time.Second).Unix(),
		Answer:    req.Answer,
	})
	if result == nil {
		c.JSON(200, response.Response{
			Code: -1,
			Msg:  "内部错误",
		})
		return
	}
	if result.BaseResp.Code != 0 {
		c.JSON(200, response.Response{
			Code: result.BaseResp.Code,
			Msg:  result.BaseResp.Msg,
		})
		return
	}
	c.JSON(200, response.Response{
		Code: 0,
		Msg:  "ok",
	})

}
