package main

import (
	"LiveLive/dao"
	websocket "LiveLive/kitex_gen/livelive/websocket/websocketservice"
	ws2 "LiveLive/rpc/websocket/ws"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	"net/http"
)

func main() {
	dao.Init()

	//创建并连接上websocket服务器
	hub := ws2.NewHub()
	go func() {
		http.HandleFunc("/ws", ws2.NewHandler(hub))
		log.Println("WebSocket server started on :8060")
		log.Fatal(http.ListenAndServe(":8060", nil))
	}()

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	websocketServiceImpl := new(WebsocketServiceImpl) //impl实现
	websocketServiceImpl.wsHub = hub

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8892")
	svr := websocket.NewServer(websocketServiceImpl,
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "livelive.websocket",
			}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
