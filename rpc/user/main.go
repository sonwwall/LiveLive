package main

import (
	"LiveLive/dao"
	user "LiveLive/kitex_gen/livelive/user/userservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"log"
)

func main() {
	dao.Init()

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	userServiceImpl := new(UserServiceImpl) //impl实现

	svr := user.NewServer(userServiceImpl,
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "livelive.user",
			}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
