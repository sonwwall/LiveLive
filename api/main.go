package main

import (
	"LiveLive/api/handlers"
	"LiveLive/api/rpc"
	"LiveLive/dao"
	"LiveLive/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	middleware.InitJwt()
	dao.Init()

	rpc.InitUserRPCClient()
	rpc.InitCourseRPCClient()

	r := server.Default(server.WithHostPorts("0.0.0.0:8080"))

	r.POST("/user/register", handlers.UserRegister)
	r.POST("/user/login", middleware.JwtMiddleware.LoginHandler)
	auth := r.Group("/auth", middleware.JwtMiddleware.MiddlewareFunc())
	auth.GET("/userinfo", handlers.UserInfo)

	teacher := auth.Group("/teacher")
	teacher.POST("/create_course", handlers.CreateCourse)

	student := auth.Group("/student")
	student.POST("/join_course", handlers.JoinCourse)

	r.Spin()
}
