package main

import (
	"LiveLive/dao"
	ai "LiveLive/kitex_gen/livelive/ai/aiservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	dao.Init()

	////连接websocketRPC
	//ws_r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	//if err != nil {
	//	panic(err)
	//}
	//ws_c, err := websocket.NewClient("livelive.websocket", client.WithResolver(ws_r))
	//if err != nil {
	//	panic(err)
	//}

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	aiServiceImpl := new(AIServiceImpl) //impl实现
	//aiServiceImpl.wsClient = ws_c

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8892")
	svr := ai.NewServer(aiServiceImpl,
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "livelive.ai",
			}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
