package socket

import (
	"fmt"
	"os"
	"testing"
	"math/rand"
)

var (
	network  = "tcp"
	address  = "127.0.0.1:8090"
	timeout  = 10
	maxConns = 100
	errCnt   = 10
)

func Test_Server(t *testing.T) {

	ser := NewServer(network, address, timeout, maxConns, errCnt)

	go processClient(ser)

	if err := ser.StartListener(); err != nil {
		fmt.Println(err)
		//错误连接过多，处理服务器上现有的连接，并退出服务
		ser.StopServer()
		os.Exit(2)
	}
}

func processClient(ser *server) {
	for {
		select {
		case client := <-ser.Clients:

			go processRch(client)

			go processWch(client)

		}
	}
}

//处理接收的数据
func processRch(client *Client) {
	for {
		select {
		case tmp_data, ok := <-client.RChan:
			if !ok {
				//通道已关闭
				break
			}
			fmt.Printf("收到客户端请求数据:%s",string(tmp_data))
			//data,err := resolveReceive(tmp_data)
			//具体业务逻辑
			//fmt.Printf("客户端-key:%s,%v,%v",client.Key,data,err)
		}
	}
}

//发送数据到客户端
func processWch(client *Client) {
	for {
		select {
		case tmp_data, ok := <-client.WChan:
			if !ok {
				break
			}
			client.Conn.Write(tmp_data)
		}
	}
}
