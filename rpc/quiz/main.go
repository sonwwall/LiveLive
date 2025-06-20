package main

import (
	"LiveLive/dao"
	quiz "LiveLive/kitex_gen/livelive/quiz/quizservice"
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

	//连接websocketRPC
	ws_r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}
	ws_c, err := websocket.NewClient("livelive.websocket", client.WithResolver(ws_r))
	if err != nil {
		panic(err)
	}

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	quizServiceImpl := new(QuizServiceImpl) //impl实现
	quizServiceImpl.wsClient = ws_c

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8891")
	svr := quiz.NewServer(quizServiceImpl,
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "livelive.quiz",
			}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
