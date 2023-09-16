package ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	go runSvrTest()

	tt := time.After(time.Second * 5)
	<-tt

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws-test", nil)
	if err != nil {

	}
	defer conn.Close()

	_ = conn.WriteMessage(websocket.BinaryMessage, []byte("client-send: hello."))

	_, res, _ := conn.ReadMessage()
	fmt.Println(res)

	select {
	case <-time.After(time.Second * 30):
		return
	}
}

func runSvrTest() {
	var svr = NewServer(context.Background(), ServerConfig{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	svr.Run()
	defer svr.Close()

	http.HandleFunc("/ws-test", func(w http.ResponseWriter, r *http.Request) {
		session, _ := svr.Create(w, r)
		for req := range session.Receive() {

			session.Send([]byte(fmt.Sprintf("server-send: %s", string(req))))
		}
	})

	_ = http.ListenAndServe(":8081", nil)
}
