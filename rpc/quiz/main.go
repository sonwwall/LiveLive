package main

import (
	"LiveLive/dao"
	quiz "LiveLive/kitex_gen/livelive/quiz/quizservice"
	"LiveLive/ws"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	"net/http"
)

func main() {
	dao.Init()

	hub := ws.NewHub()
	go func() {
		http.HandleFunc("/ws", ws.NewHandler(hub))
		log.Println("WebSocket server started on :8060")
		log.Fatal(http.ListenAndServe(":8060", nil))
	}()
	//NewQuizService(hub)

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	quizServiceImpl := new(QuizServiceImpl) //impl实现
	quizServiceImpl.wsHub = hub

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
