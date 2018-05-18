package socket

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
	"encoding/json"
)

//Client对象的一些初始化操作
func (this *server) newConnection(conn net.Conn) (*Client, error) {
	this.Lock()
	defer this.Unlock()
	if len(this.allClients) >= this.maxConns {
		return nil, ErrTooManyConns
	}
	key := generateAuthKey()
	authKey := generateAuthKey()	//生成鉴权key,由客户端去解密,再次传递过来
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

//验证客户端是否有效连接
func (this *server) authClient(client *Client){
	data := make(map[string]interface{})
	data["authKey"] = client.AuthKey
	ret,_ := json.Marshal(data)
	client.Conn.Write(Packet(ret))
}


//处理客户端连接
func (this *server) handleConnection(client *Client) {
	//这里是鉴权，往客户端写入一个加密的字符串，由客户端去解析。
	//如果5秒没有响应，或者客户端下次请求时没有带上正确的key，则认为该客户端非法，服务端会自动抛弃该连接
	//client.Conn.Write([]byte(client.AuthKey))
	this.authClient(client)
	client.Conn.SetDeadline(time.Now().Add(5 * time.Second)) //5秒后，没有响应，断开连接

	var buffer []byte
	for {
		tmpBuf := make([]byte, 1024)
		n, err := client.Conn.Read(tmpBuf)
		if err != nil {
			if err == io.EOF {
				//客户端已关闭连接
				fmt.Printf("client:%v is disconnect!\n",client.AuthKey)
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("server has closed!client:%v\n",client.AuthKey)
			} else {
				fmt.Printf("client:%v,conn.Read() has error:%v\n", client.AuthKey,err)
			}
			client.disconnect()
			return
		}
		buffer = UnPacket(append(buffer, tmpBuf[:n]...), client.RChan)
		client.Conn.SetDeadline(time.Now().Add(time.Duration(this.timeout) * time.Second)) //设置超时时间
	}
}





//关闭连接
func (this *Client) disconnect() {
	if !this.Flag {
		close(this.RChan)
		close(this.WChan)
		this.Conn.Close()
		delete(ser.allClients, this.AuthKey)
	}
	this.Flag = true

}
