package main

import (
	_ "bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math"
	"net"
	"net/http"
	"sync"
)

func init() {
	rm = RoomManager{
		clients: make(map[uint64]*Client),
		rooms:   make(map[uint64]*Room),
		genUID:  NewGenerateUID(),
		genRID:  NewGenerateRID(),
		parse:   &Parse{},
	}

}

type WSHandler struct {
	upgrader websocket.Upgrader
	wsconns  map[*WSConn]struct{}
	sync.Mutex
	msg Message
}

type WSConn struct {
	conn      *websocket.Conn
	writeChan chan []byte
}

func (this *WSConn) Send(b []byte) {
	this.writeChan <- b

}

type Message interface {
	OnMessage(client *Client, data []byte)
}

type Mux struct {
}

func (this *Mux) OnMessage(client *Client, data []byte) {
	// fmt.Printf("OnMessage %v\n", string(data))
	rm.DealMsg(client, data)

}

type RoomManager struct {
	clients map[uint64]*Client `客户端列表`
	rooms   map[uint64]*Room   `房间列表`
	genUID  *GenerateUID       `用户ID生成器`
	genRID  *GenerateRID       `房间ID生成器`
	sync.Mutex
	parse Processer
}

type CMD int

const (
	Login CMD = 1 + iota
	Croom
	JHall
)

type LoginHall struct {
	Msg string
}

type RevJRoom struct {
	RId uint64
}

type SendCRoomProto struct {
	Code   int
	RoomNo uint64
	Desc   string
}

func (this *RoomManager) DealMsg(client *Client, data []byte) {

	cmd, jsonData := this.parse.Unmarshal(data)
	fmt.Printf("DealMsg cmd:%d data:%s\n", cmd, jsonData)

	switch cmd {
	case 1:
		var ln LoginHall
		if err := json.Unmarshal([]byte(jsonData), &ln); err != nil {
			fmt.Println("Unmarshal err%v", err)

		}
		// fmt.Printf("Login Hall %v\n", ln)
		this.AddUser(client)
		user := client.user
		lData, err := this.parse.Marshal(cmd, user)
		if err != nil {
			return
		}
		this.SendClientData(user.Id, lData)
	case 2:
		// var revCroom RevCRoomProto
		var code int
		var rid uint64
		var desc string
		if client.user.Id > 0 {
			rid = client.user.Id
			this.JoinRoom(client, rid)
		} else {
			rid = this.CreateRoom(client)
		}

		fmt.Println(rid)
		sendRoomProto := &SendCRoomProto{
			Code:   code,
			RoomNo: rid,
			Desc:   desc,
		}
		sdata, _ := this.parse.Marshal(cmd, sendRoomProto)
		this.SendClientData(client.user.Id, sdata)
	case 3:
		var jroom RevJRoom
		if err := json.Unmarshal([]byte(jsonData), &jroom); err != nil {
			fmt.Println("Unmarshal err ", err)
			return
		}
		this.JoinRoom(client, jroom.RId)
		sdata, _ := this.parse.Marshal(cmd, &jroom)
		this.BroadcastRoom(jroom.RId, sdata)
		// default:

	}

}

func (this *RoomManager) AddUser(client *Client) {
	defer this.Unlock()
	this.Lock()
GenUID:
	uid := this.genUID.Generate()
	if _, ok := this.clients[uid]; ok {
		goto GenUID
	}
	this.clients[uid] = client
	client.user = User{
		Id:   uid,
		Name: fmt.Sprintf("test%d", uid),
	}

	// fmt.Println("add user:", client.user.Name)

}

func (this *RoomManager) CreateRoom(client *Client) uint64 {
	defer this.Unlock()
	this.Lock()
	if client.user.roomNo > 0 {
		return client.user.roomNo
	}
GenRID:
	rid := this.genRID.Generate()
	if _, ok := this.rooms[rid]; ok {
		goto GenRID
	}
	room := NewRoom(this, rid)

	this.rooms[rid] = room
	client.user.roomNo = rid

	room.JoinUser(client.user.Id)

	return rid
}

func (this *RoomManager) JoinRoom(client *Client, rid uint64) bool {
	defer this.Unlock()
	this.Lock()
	if room, ok := this.rooms[rid]; ok {
		client.user.roomNo = rid
		room.JoinUser(client.user.Id)

		return true

	}
	return false

}

func (this *RoomManager) SendClientData(uid uint64, data []byte) {

	if client, ok := this.clients[uid]; ok {
		client.conn.Send(data)

	} else {
		fmt.Printf("不存在的用户ID%d\n", uid)
	}

}

func (this *RoomManager) BroadcastRoom(rid uint64, data []byte) {
	if room, ok := this.rooms[rid]; ok {
		for _, uid := range room.users {
			if uid > 0 {
				this.SendClientData(uid, data)
			}

		}

	}

}

type Room struct {
	roomNo uint64
	sync.Mutex
	users [4]uint64
	rm    *RoomManager
}

func NewRoom(rm *RoomManager, rid uint64) *Room {
	return &Room{
		// users: [4]int,
		roomNo: rid,
		rm:     rm,
	}

}

func (this *Room) JoinUser(userId uint64) bool {
	defer this.Unlock()
	this.Lock()
	for key, value := range this.users {
		if value == 0 {
			this.users[key] = userId
			return true
		}
	}
	// this.users = append(this.users, userId)
	return false
}

func (this *Room) RoomLogic(uid uint64, cmd uint32, jsonData string) {
	switch cmd {

	}

}

func (this *Room) ExitUser(userId uint64) bool {
	defer this.Unlock()
	this.Lock()
	for key, value := range this.users {
		if value == userId {
			this.users[key] = 0
			return true

		}

	}
	return false
}

type Client struct {
	user User
	conn *WSConn
}

func NewClient(conn *WSConn) *Client {
	return &Client{conn: conn}

}

type User struct {
	Id     uint64
	Name   string
	roomNo uint64
}

type Server struct {
	handle *WSHandler
	sync.Mutex
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
		msg:     &Mux{},
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

			err := wscon.conn.WriteMessage(websocket.BinaryMessage, b)

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
	client := NewClient(wsconn)
	this.Lock()
	this.wsconns[wsconn] = struct{}{}
	this.Unlock()

	for {
		_, data, err := wsconn.conn.ReadMessage()
		if err != nil {
			fmt.Printf("read message err %v\n", err)
			break
		}
		// fmt.Println(conn)
		// wsconn.Send(data)
		this.msg.OnMessage(client, data)
	}

	this.Lock()
	delete(this.wsconns, wsconn)
	this.Unlock()
	fmt.Printf("client close %v\n", wsconn.conn.LocalAddr())

}

var rm RoomManager

func main() {
	server := NewServer()
	server.Start()
}

type GenerateUID struct {
	cuid uint64
	sync.Mutex
}

func (this *GenerateUID) Generate() uint64 {
	defer this.Unlock()
	this.Lock()
	if this.cuid >= math.MaxInt64 {
		this.cuid = 0
	}
	this.cuid++
	return this.cuid
}

func NewGenerateUID() *GenerateUID {
	return &GenerateUID{
		cuid: 0,
	}

}

type GenerateRID struct {
	crid uint64
	sync.Mutex
}

func (this *GenerateRID) Generate() uint64 {
	defer this.Unlock()
	this.Lock()
	if this.crid >= math.MaxInt64 {
		this.crid = 0
	}
	this.crid++
	return this.crid
}

func NewGenerateRID() *GenerateRID {
	return &GenerateRID{
		crid: 0,
	}
}

type Processer interface {
	Unmarshal(data []byte) (cmd uint32, json string)
	Marshal(cmd uint32, data interface{}) ([]byte, error)
}
type Parse struct {
}

func (this *Parse) Unmarshal(data []byte) (cmd uint32, json string) {
	cmd = binary.LittleEndian.Uint32(data)
	json = string(data[4:])
	fmt.Println(json)
	return
}

func (this *Parse) Marshal(cmd uint32, data interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(data)
	dataByte := []byte(jsonData)
	cmdByte := make([]byte, len(dataByte)+4)
	binary.BigEndian.PutUint32(cmdByte, cmd)
	copy(cmdByte[4:], dataByte)

	return cmdByte, err
}
