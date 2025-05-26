package ws

import (
	"log"
	"net/http"
)

func InitWebsocket() {
	hub := NewHub()

	http.HandleFunc("/ws", NewHandler(hub))

	//go func() {
	//	// 模拟老师发布题目
	//	// 实际可以改为接收 HTTP 请求或调用
	//	for {
	//		msg := []byte(`{"type":"choice_question","data":{"title":"1+1=?","options":["1","2","3","4"]}}`)
	//		hub.BroadcastToCourse(123, msg)
	//		time.Sleep(5 * time.Second)
	//	}
	//}()

	log.Println("WebSocket server started on :8060")
	log.Fatal(http.ListenAndServe(":8060", nil))
}
