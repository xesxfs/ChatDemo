package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type WSHandler struct {
	upgrader websocket.Upgrader
}

func (this *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error")
		return
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message err %v", err)
			break
		}
		fmt.Println("read msg :" + string(data))
		// fmt.Println(conn)

	}

}

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen err")
		return
	}

	handler := &WSHandler{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 1000 * 10,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10000,
		WriteTimeout:   10000,
		MaxHeaderBytes: 1024,
	}

	httpServer.Serve(ln)
	fmt.Println("close")

}
