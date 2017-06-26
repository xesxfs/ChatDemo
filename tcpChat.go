package main

import (
	"encoding/binary"
	"fmt"
	_ "io"
	"net"
	// "sync"
	_ "time"
)

var funMap map[string]func(data interface{})

func packData(cmd uint32, data string) []byte {
	dataByte := []byte(data)
	cmdByte := make([]byte, len(dataByte)+4)
	binary.BigEndian.PutUint32(cmdByte, cmd)
	copy(cmdByte[4:], dataByte)
	return cmdByte
}

func unPackData(data []byte) (cmd uint32, rdata string) {
	cmd = binary.BigEndian.Uint32(data)
	// binary.BigEndian.Uint32(b)
	rdata = string(data[4:])
	return
}

func handleConn(tcpCon *TcpConn) {
	b := make([]byte, 512)
	for {
		cn, err := tcpCon.conn.Read(b)
		if err != nil {
			// fmt.Println(err.Error(), cn)
			// fmt.Println("close conn1")
			// var closeCmd uint32 = 500
			// strconv.Itoa(closeCmd)
			// var closeBuf = make([]byte, 4)
			// binary.BigEndian.PutUint32(closeBuf, closeCmd)
			closeData := packData(500, "")
			tcpCon.msg.OnMessage(tcpCon, closeData)
			break
		}
		// result := fmt.Sprintf("%d:%s", tcpCon.id, string(b[:cn]))
		tcpCon.msg.OnMessage(tcpCon, b[:cn])

		// fmt.Println(result)
		// tcpCon.Write([]byte(result))
	}

}

func init() {
	tabMgr = NewTableManager()

}

// func decode(data string) {
// 	cmd := data
// 	if f, ok := funMap[cmd]; ok {
// 		f(data)
// 	}
// }

// func Register(cmd string, f func(data interface{})) {
// 	// fun,ok= funMap(cmd)
// 	funMap[cmd] = f
// }

type TableManager struct {
	tables   []int
	tablesNo []string
	users    map[int]*TcpConn
}

func (this *TableManager) OnMessage(tcpConn *TcpConn, data []byte) {
	cmd, rdata := unPackData(data)
	fmt.Println(cmd, rdata)
	switch cmd {
	case 1:
		this.users[tcpConn.id] = tcpConn
		this.Broadcast(1, "Join")
	case 2:
	case 3:
	case 4:
	case 5:
	case 500:
		if _, ok := this.users[tcpConn.id]; ok {
			delete(this.users, tcpConn.id)
		}
		this.Broadcast(500, "Exit")
	}
	// fmt.Println(data)
	// tcpConn.Write(data)

}

func (this *TableManager) Broadcast(cmd uint32, data string) {
	sendData := packData(cmd, data)
	for _, value := range this.users {
		value.Write(sendData)
	}

}

func NewTableManager() *TableManager {
	tm := &TableManager{}
	tm.users = make(map[int]*TcpConn)
	return tm
}

type TcpConn struct {
	conn      net.Conn
	writeChan chan []byte
	id        int
	msg       Message
}

type Message interface {
	OnMessage(tcpConn *TcpConn, data []byte)
}

func (this *TcpConn) Write(b []byte) {
	this.writeChan <- b
}

func NewTcpConn(conn net.Conn, id int) *TcpConn {
	tcpConn := &TcpConn{}
	tcpConn.conn = conn
	tcpConn.writeChan = make(chan []byte, 512)
	tcpConn.id = id
	tcpConn.msg = tabMgr
	go handleWirte(tcpConn)
	return tcpConn
}

func handleWirte(tcpConn *TcpConn) {
	for b := range tcpConn.writeChan {
		tcpConn.conn.Write(b)
	}

}

var conns map[*TcpConn]struct{}
var tabMgr *TableManager

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen error")

	}

	conns = map[*TcpConn]struct{}{}
	iCount := 0

	// laddr := ln.Addr()
	// fmt.Println(laddr.String())
	for {

		conn, err := ln.Accept()
		if err != nil {
			continue

		}
		iCount++
		// addr := conn.RemoteAddr()
		// fmt.Println(addr.Network(), addr.String())
		// go handleWirte(conn)
		go func() {
			tcpCon := NewTcpConn(conn, iCount)
			conns[tcpCon] = struct{}{}
			handleConn(tcpCon)
			delete(conns, tcpCon)
			fmt.Println("close", tcpCon.id)

		}()

	}

}
