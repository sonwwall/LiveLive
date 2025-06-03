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
	rpc.InitLiveRPCClient()
	rpc.InitQuizRPCClient()

	r := server.Default(server.WithHostPorts("0.0.0.0:8080"))

	r.POST("/user/register", handlers.UserRegister)
	r.POST("/user/login", middleware.JwtMiddleware.LoginHandler)
	auth := r.Group("/auth", middleware.JwtMiddleware.MiddlewareFunc())
	auth.GET("/userinfo", handlers.UserInfo)

	teacher := auth.Group("/teacher")
	teacher.POST("/create_course", handlers.CreateCourse)
	teacher.POST("/create_courseInvite", handlers.CreateCourseInvite)
	teacher.GET("/streamKey", handlers.GetStreamKey)
	teacher.POST("/choice_question", handlers.PublishChoiceQuestion)
	teacher.POST("/publish_register", handlers.PublishRegister)
	teacher.POST("/start_recording", handlers.StartRecording)
	teacher.POST("/stop_recording", handlers.StopRecording)

	student := auth.Group("/student")
	student.POST("/join_course", handlers.JoinCourse)
	student.GET("/watch_live", handlers.WatchLive)

	r.Spin()
}
