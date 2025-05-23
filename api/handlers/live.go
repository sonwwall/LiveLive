package handlers

import (
	"LiveLive/api/code"
	"LiveLive/api/rpc"
	"LiveLive/kitex_gen/livelive/live"
	"LiveLive/middleware"
	"LiveLive/model"
	"LiveLive/utils/response"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetStreamKey(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}

	var req model.LiveSession
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误",
		})
		return
	}

	req.TeacherID = int64(user.(*model.User).ID)

	result, _ := rpc.GetStreamKey(ctx, &live.GetStreamKeyReq{
		Classname: req.ClassName,
		TeacherId: req.TeacherID,
	})
	if result == nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "内部错误",
		})
		return
	}

	c.JSON(200, response.Response{
		Code: 0,
		Msg:  "ok",
		Data: map[string]string{
			"RtmpUrl":   result.RtmpUrl,
			"StreamKey": result.StreamKey,
		},
	})
}
