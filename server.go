package socket

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var (
	ser *server
)

type server struct {
	sync.Mutex
	network    string
	address    string
	timeout    int
	maxConns   int
	errCnt     int
	maxErrCnt  int
	Listener   net.Listener
	Clients    chan *Client
	allClients map[string]*Client
}

type Client struct {
	Key      string
	AuthKey  string
	Flag     bool //true:连接已关闭;false:连接未关闭
	Conn     net.Conn
	ConnectT time.Time
	RChan    chan []byte
	WChan    chan []byte
}

func NewServer(network, address string, timeout, maxConns, maxErrCnt int) *server {
	ser = &server{
		network:    network,
		address:    address,
		timeout:    timeout,
		maxConns:   maxConns,
		maxErrCnt:  maxErrCnt,
		Clients:    make(chan *Client),
		allClients: make(map[string]*Client),
	}
	return ser
}

func (this *server) StartListener() error {
	listener, err := net.Listen(this.network, this.address)
	if err != nil {
		fmt.Printf("net.Listen has error:%v\n", err)
		os.Exit(1)
	}
	this.Listener = listener
	for {
		conn, err := this.Listener.Accept()
		if err != nil {
			this.errCnt++
			if this.errCnt >= this.maxErrCnt {
				//超过最大错误数量，系统停止服务
				break
			}
			fmt.Printf("listener.Accept has error:%v\n", err)
			continue
		}
		if client, err := this.newConnection(conn); err == nil {
			go this.handleConnection(client)
		} else {
			return err
		}
	}
	return ErrTooManyErrConns
}

//停止服务,处理现有的客户端连接
func (this *server) StopServer() {
	for _,client := range this.allClients{
		client.disconnect()
	}
}
