package main

import (
	"LiveLive/dao"
	course "LiveLive/kitex_gen/livelive/course/courseservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	dao.Init()

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	courseServiceImpl := new(CourseServiceImpl) //impl实现

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8889")
	svr := course.NewServer(courseServiceImpl,
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "livelive.course",
			}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
