package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	// "time"
)

func handleWrite(conn net.Conn) {
	// hw := "hello world"
	// hw2 := "hello"

	// cmd := packData(200, "hello")

	// binary.BigEndian.PutUint32(cmd, 200)
	// b := []byte(hw)
	// conn.Write(cmd)

	cmd := packData(1, "")
	conn.Write(cmd)

	// c := []byte(hw2)
	// conn.Write(c)

}

func packData(cmd uint32, data string) []byte {
	dataByte := []byte(data)
	cmdByte := make([]byte, len(dataByte)+4)
	binary.BigEndian.PutUint32(cmdByte, cmd)
	copy(cmdByte[4:], dataByte)
	return cmdByte
}

func unPackData(data []byte) (cmd uint32, rdata string) {
	cmd = binary.BigEndian.Uint32(data)
	rdata = string(data[4:])
	return
}

func main() {

	tcpWaite := sync.WaitGroup{}

	for i := 0; i < 1; i++ {
		tcpWaite.Add(1)

		go func() {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {

			}
			// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
			go handleWrite(conn)

			b := make([]byte, 512)
			for {
				cn, err := conn.Read(b)
				if err != nil {
					break
				}
				cmd, data := unPackData(b[:cn])
				fmt.Println(cmd, data)
			}

			tcpWaite.Done()

		}()

	}

	tcpWaite.Wait()
	// time.Sleep(20 * time.Second)

}
