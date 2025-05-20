package handlers

import (
	"LiveLive/api/code"
	"LiveLive/middleware"
	"LiveLive/model"
	"LiveLive/utils/response"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func CreateCourse(ctx context.Context, c *app.RequestContext) {

	user, _ := c.Get(middleware.IdentityKey)
	//log.Println("user:", user.(*model.User).Username)
	//log.Println("role:", user.(*model.User).Role)
	//
	if user.(*model.User).Role != 1 {
		c.JSON(200, response.Response{
			Code: code.ErrNoPermission,
			Msg:  "抱歉，您无权访问",
		})
	}
}
