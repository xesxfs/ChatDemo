package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
)

type WSHandler struct {
	upgrader websocket.Upgrader
	wsconns  map[*WSConn]struct{}
	sync.Mutex
}
type WSConn struct {
	conn      *websocket.Conn
	writeChan chan []byte
	msg       Message
}

type Message interface {
	OnMessage(conn WSConn, data []byte)
}

type Server struct {
	handle *WSHandler
	sync.Mutex
}

func (this *WSConn) Send(b []byte) {
	this.writeChan <- b

}

func (this *Server) Start() {

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
		wsconns: make(map[*WSConn]struct{}),
	}

	httpServer := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10000,
		WriteTimeout:   10000,
		MaxHeaderBytes: 1024,
	}
	this.handle = handler

	httpServer.Serve(ln)
	// httpServer.Shutdown(ctx)
	fmt.Println("close")

}

func NewWSConn(conn *websocket.Conn) *WSConn {
	wscon := &WSConn{
		conn:      conn,
		writeChan: make(chan []byte),
	}

	go func() {

		for b := range wscon.writeChan {
			if b == nil {
				break
			}

			err := wscon.conn.WriteMessage(websocket.TextMessage, b)

			if err != nil {
				break
			}

		}

		wscon.conn.Close()

	}()

	return wscon
}

func NewServer() *Server {
	return &Server{
	// conns: map[*WSConn]struct{}{},
	// conns: make(map[*WSConn]struct{}),
	}
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
	fmt.Printf("%v\n", conn.LocalAddr())

	wsconn := NewWSConn(conn)
	this.Lock()
	this.wsconns[wsconn] = struct{}{}
	this.Unlock()

	for {
		_, data, err := wsconn.conn.ReadMessage()
		if err != nil {
			fmt.Printf("read message err %v\n", err)
			break
		}
		fmt.Printf("read msg %v\n", string(data))
		// fmt.Println(conn)
		wsconn.Send(data)

	}

	this.Lock()
	delete(this.wsconns, wsconn)
	this.Unlock()
	fmt.Printf("client close %v\n", wsconn.conn.LocalAddr())

}

func main() {
	server := NewServer()
	server.Start()

}
