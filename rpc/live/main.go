package main

import (
	"LiveLive/dao"
	"LiveLive/kafka/consumer"
	ai "LiveLive/kitex_gen/livelive/ai/aiservice"
	live "LiveLive/kitex_gen/livelive/live/liveservice"
	websocket "LiveLive/kitex_gen/livelive/websocket/websocketservice"
	"github.com/cloudwego/kitex/client"
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

	//连接websocketRPC
	ws_r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	ws_c, err := websocket.NewClient("livelive.websocket", client.WithResolver(ws_r))
	if err != nil {
		panic(err)
	}

	//连接aiRPC
	ai_r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	ai_c, err := ai.NewClient("livelive.ai", client.WithResolver(ai_r))
	if err != nil {
		panic(err)
	}

	liveServiceImpl := new(LiveServiceImpl) //impl实现
	liveServiceImpl.wsClient = ws_c
	liveServiceImpl.aiClient = ai_c

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
