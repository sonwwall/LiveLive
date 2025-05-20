package rpc

import (
	"LiveLive/kitex_gen/livelive/user"
	"LiveLive/kitex_gen/livelive/user/userservice"
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	userClient userservice.Client
)

func InitUserRPCClient() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient("livelive.user", client.WithResolver(r))
	if err != nil {
		panic(err)
	}

	userClient = c
}

func Register(ctx context.Context, req *user.RegisterReq) (*user.RegisterResp, error) {
	return userClient.Register(ctx, req)
}

func Login(ctx context.Context, req *user.LoginReq) (*user.LoginResp, error) {
	return userClient.Login(ctx, req)
}

func UserInfo(ctx context.Context, req *user.UserInfoReq) (*user.UserInfoResp, error) {
	return userClient.UserInfo(ctx, req)
}
