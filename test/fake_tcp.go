package test

import (
	"encoding/hex"
	"fmt"
	"net"
)

const cmdLen = 12
func connHandler(conn net.Conn) {
	//创建消息缓冲区
	buffer := make([]byte, 260)
	for {
		//读取客户端发来的消息放入缓冲区
		n,err := conn.Read(buffer[:cmdLen])
		if err != nil {
			fmt.Println(fmt.Sprintf("read error: %v",err))
			return
		}
		//转化为字符串输出
		clientMsg := hex.EncodeToString(buffer[:n])
		fmt.Println(fmt.Sprintf("收到消息 %v %v",conn.RemoteAddr(),clientMsg))
		//回复客户端消息
		reply := []byte{buffer[0], buffer[1],0x00, 0x00, 0x00, 0x05, 0x00, 0x03 ,0x02 ,0x00, 0x2b, 0xff}
		conn.Write(reply)
		fmt.Println(fmt.Sprintf("回复消息 %v",hex.EncodeToString(reply)))
	}
	conn.Close()
	fmt.Println(fmt.Printf("客户端断开连接 %v",conn.RemoteAddr()))
}
func tcp() {
	addr, err2 := net.ResolveTCPAddr("tcp4", "127.0.0.1:503")
	if err2 != nil {
		panic(err2)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println(fmt.Sprintf("server on %v ", ln.Addr().String()))

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(fmt.Errorf("accept error: %v", err))
		}
		fmt.Println("******accept*********")
		go connHandler(conn)
		fmt.Println("******back*********")
	}
}
