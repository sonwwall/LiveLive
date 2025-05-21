package handlers

import (
	"LiveLive/api/code"
	"LiveLive/api/rpc"
	"LiveLive/kitex_gen/livelive/course"
	"LiveLive/middleware"
	"LiveLive/model"
	"LiveLive/utils"
	"LiveLive/utils/response"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

func CreateCourse(ctx context.Context, c *app.RequestContext) {

	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}

	var req model.Course //不写指针
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误：" + err.Error(),
		})
		return
	}

	result, _ := rpc.CreateCourse(ctx, &course.CreateCourseReq{
		Classname:   req.Classname,
		Description: req.Description,
		TeacherId:   int64(user.(*model.User).Model.ID),
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
		Msg:  "创建成功",
	})

}

func CreateCourseInvite(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}

	var req model.CourseInvite
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误" + err.Error(),
		})
		return
	}

	//从指针转换为值类型
	var maxUsage int64
	if req.MaxUsage != nil {
		maxUsage = *req.MaxUsage
	}

	result, _ := rpc.CreateCourseInvite(ctx, &course.CreateCourseInviteReq{
		Classname: req.Classname,
		MaxUsage:  maxUsage,
		ExpiredAt: utils.PtrToTimestamp(req.ExpiredAt),
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
		Data: fmt.Sprintf("邀请码为：%s", result.InviteCode),
	})

}

func JoinCourse(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(middleware.IdentityKey)
	if user.(*model.User).Role != 2 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
		return
	}
	var req model.CourseMember
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数错误：" + err.Error(),
		})
		return
	}
	result, _ := rpc.JoinCourse(ctx, &course.JoinCourseReq{
		Classname: req.Classname,
		StudentId: int64(user.(*model.User).Model.ID),
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
		Msg:  "加入成功",
	})
}
