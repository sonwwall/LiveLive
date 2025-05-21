package rpc

import (
	"LiveLive/kitex_gen/livelive/course"
	"LiveLive/kitex_gen/livelive/course/courseservice"
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	courseClient courseservice.Client
)

func InitCourseRPCClient() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}

	c, err := courseservice.NewClient("livelive.course", client.WithResolver(r))
	if err != nil {
		panic(err)
	}

	courseClient = c
}

func CreateCourse(ctx context.Context, req *course.CreateCourseReq) (*course.CreateCourseResp, error) {
	return courseClient.CreateCourse(ctx, req)
}

func JoinCourse(ctx context.Context, req *course.JoinCourseReq) (*course.JoinCourseResp, error) {
	return courseClient.JoinCourse(ctx, req)
}

func CreateCourseInvite(ctx context.Context, req *course.CreateCourseInviteReq) (*course.CreateCourseInviteResp, error) {
	return courseClient.CreateCourseInvite(ctx, req)
}
