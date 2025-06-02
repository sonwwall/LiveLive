package handlers

import (
	"LiveLive/api/code"
	"LiveLive/api/rpc"
	"LiveLive/kitex_gen/livelive/live"
	"LiveLive/middleware"
	"LiveLive/model"
	"LiveLive/utils/response"
	"context"
	"fmt"
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
			Msg:  "参数错误" + err.Error(),
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

func WatchLive(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 2 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}

	type watchReq struct {
		Classname   string `json:"classname,required" form:"classname,required"`
		TeacherName string `json:"teacher_name,required" form:"teacher_name,required"`
	}

	var req watchReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误：" + err.Error(),
		})
		return
	}

	result, _ := rpc.WatchLive(ctx, &live.WatchLiveReq{
		TeacherName: req.TeacherName,
		Classname:   req.Classname,
		StudentId:   int64(user.(*model.User).ID),
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
	//  手动设置目标请求路径
	//c.Request.URI().SetRequestURI(result.Addr) // 例如 "/live/teacher_3_course_2_teacher01_class01.flv"
	//c.Request.SetRequestURI(result.Addr)
	//log.Printf("url:http://127.0.0.1:7001%s", c.Request.RequestURI())

	//// 创建代理（target 是 base host）
	//proxy, err := reverseproxy.NewSingleHostReverseProxy("http://127.0.0.1:7001")
	//if err != nil {
	//	c.JSON(500, map[string]string{"error": "proxy error"})
	//	return
	//}
	//
	////  代理转发（此时 JoinURLPath 拼接就是对的）
	//proxy.ServeHTTP(ctx, c)
	c.JSON(200, response.Response{

		Code: 0,
		Msg:  "ok",
		Data: map[string]string{
			"url": fmt.Sprintf("http://127.0.0.1%s", result.Addr),
		},
	})
}

func PublishRegister(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}
	type PublishRegisterReq struct {
		Classname string `json:"classname,required" form:"classname,required"`
	}
	var req PublishRegisterReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误:" + err.Error(),
		})
		return
	}
	result, _ := rpc.PublishRegister(ctx, &live.PublishRegisterReq{
		Classname:   req.Classname,
		TeacherId:   int64(user.(*model.User).ID),
		TeacherName: user.(*model.User).Username,
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
