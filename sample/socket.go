package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "8081", "port")
func main() {
	//启动服务器
	go service()
	time.Sleep(1e9)

	//启动客户端
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecting to " + *host + ":" + *port)
	done := make(chan string)
	go handleWrite(conn, done)
	go handleRead(conn, done)
	fmt.Println(<-done)
	fmt.Println(<-done)
}

func handleWrite(conn net.Conn, done chan string) {
	for i := 4; i > 0; i-- {
		_, e := conn.Write([]byte("hello " + strconv.Itoa(i) + ", "))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	done <- "Sent"
}

func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println("client read:", string(buf[:reqLen]))
	done <- "Read"
}

func service() {
	ln, _ := net.Listen("tcp", ":8081")
	conn, _  := ln.Accept()

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Service read:", string(buf[:n]))

	_, _ = conn.Write([]byte("test123"))
}