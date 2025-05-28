package main

import (
	"LiveLive/dao"
	"LiveLive/kafka/consumer"
	live "LiveLive/kitex_gen/livelive/live/liveservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	dao.Init()
	go consumer.StartKafkaConsumerChat() //开启一个聊天的kafka进程

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	liveServiceImpl := new(LiveServiceImpl) //impl实现

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8890")
	svr := live.NewServer(liveServiceImpl,
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "livelive.live",
			}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
