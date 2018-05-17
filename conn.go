package socket

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

//Client对象的一些初始化操作
func (this *server) newConnection(conn net.Conn) (*Client, error) {
	this.Lock()
	defer this.Unlock()
	if len(this.allClients) >= this.maxConns {
		return nil, ErrTooManyConns
	}
	key := "123456"
	authKey := "123456"
	client := &Client{
		Key:      key,
		AuthKey:  authKey,
		Conn:     conn,
		ConnectT: time.Now(),
		RChan:    make(chan []byte),
		WChan:    make(chan []byte),
	}

	this.Clients <- client
	this.allClients[key] = client

	return client, nil
}

//处理客户端连接
func (this *server) handleConnection(client *Client) {
	//这里是鉴权，往客户端写入一个加密的字符串，由客户端去解析。
	//如果5秒没有响应，或者客户端无法解析，则任务该客户端非法，服务端会自动抛弃该连接
	client.Conn.Write([]byte(client.AuthKey))
	client.Conn.SetDeadline(time.Now().Add(5 * time.Second)) //5秒后，没有响应，断开连接

	buffer := make([]byte, 2048)
	for {
		tmpBuf := make([]byte, 1024)
		n, err := client.Conn.Read(tmpBuf)
		if err != nil {
			if err == io.EOF {
				//客户端已关闭连接
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Println("服务器端已断开连接!")
				client.Flag = false
			} else {
				//其他异常
				fmt.Printf("conn.Read() has error:%v\n", err)
			}
			go client.disconnect()
		}
		buffer = UnPacket(append(buffer, tmpBuf[:n]...), client.RChan)
		client.Conn.SetDeadline(time.Now().Add(time.Duration(this.timeout) * time.Second)) //设置超时时间
	}
}

//关闭连接
func (this *Client) disconnect() {
	close(this.RChan)
	close(this.WChan)
	if !this.Flag {
		this.Conn.Close()
	}
	fmt.Println(ser.allClients)
	fmt.Println(this.AuthKey)
	delete(ser.allClients, this.AuthKey)
}
