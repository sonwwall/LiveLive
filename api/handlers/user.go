package handlers

import (
	"LiveLive/api/code"
	"LiveLive/api/rpc"
	"LiveLive/kitex_gen/livelive/user"
	"LiveLive/model"
	utils2 "LiveLive/utils/md5"

	"LiveLive/utils/response"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

func UserRegister(ctx context.Context, c *app.RequestContext) {
	var req model.User
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  "参数绑定失败：" + err.Error(),
		})
		return
	}
	req.Password = utils2.MD5(req.Password)

	result, _ := rpc.Register(ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Role:     req.Role,
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
		Msg:  "注册成功",
	})
}

func UserLogin(ctx context.Context, c *app.RequestContext) {
	var req model.User
	err := c.Bind(&req)
	if err != nil {
		c.JSON(200, response.Response{
			Code: code.ErrInvalidParams,
			Msg:  err.Error(),
		})
		return
	}
	result, _ := rpc.Login(ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
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
		Msg:  "登录成功",
	})
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	claims := jwt.ExtractClaims(ctx, c)

	username, _ := claims["identity"].(string)
	result, _ := rpc.UserInfo(context.Background(), &user.UserInfoReq{
		Username: username,
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
		Data: user.UserInfoResp{
			Username: result.Username,
			Mobile:   result.Mobile,
			Email:    result.Email,
		},
	})

}
