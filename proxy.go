package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func proxy_down(client, server net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := server.Read(buf)
		if err != nil {
			fmt.Printf("[DOWN]Failure to read:%s\n", err.Error())
			return
		}
		_, err = client.Write(buf)
		if err != nil {
			fmt.Printf("[DOWN]Failure to write: %s\n", err.Error())
			return
		}
		fmt.Printf("[DOWN]%s\n", buf)
	}
}

func proxy_up(client net.Conn) {
	defer client.Close()
	hostport, err := bufio.NewReader(client).ReadString('\n')
	if err != nil {
		fmt.Printf("[INIT]Failure to read:%s\n", err.Error())
		return
	}
	hostport = strings.TrimSpace(hostport)
	fmt.Printf("[INIT]Proxy to =%s=\n", hostport)
	server, err := net.Dial("tcp", hostport)
	if err != nil {
		fmt.Printf("[INIT]Failure to connect %s:%s\n", hostport, err.Error())
		return
	}
	defer server.Close()

	go proxy_down(client, server)

	for {
		buf := make([]byte, 1024)
		_, err := client.Read(buf)
		if err != nil {
			fmt.Printf("[UP]Failure to read:%s\n", err.Error())
			return
		}
		_, err = server.Write(buf)
		if err != nil {
			fmt.Printf("[UP]Failure to write: %s\n", err.Error())
			return
		}
		fmt.Printf("[UP]%s\n", buf)
	}
}

func main() {
	fmt.Printf("Server is ready...\n")
	l, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Printf("Failure to listen: %s\n", err.Error())
	}

	for {
		if c, err := l.Accept(); err == nil {
			go proxy_up(c)
		}
	}
}
