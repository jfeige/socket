package socket

import (
	"net"
	"fmt"
	"io"
	"encoding/json"
	"sync"
	"time"
	"strconv"
)

var(
	cnetwork = "tcp"
	caddress = "127.0.0.1:8090"
	RChan = make(chan []byte)
	CloseChan = make(chan bool)	//连接关闭标示
	AuthKey interface{}		//用于验证连接是否有效，服务器返回
	wg sync.WaitGroup
)

func main(){
	addr,err := net.ResolveTCPAddr(cnetwork,caddress)
	if err != nil{
		fmt.Println(err)
		return
	}
	conn,err := net.DialTCP(cnetwork,nil,addr)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("client begin.....")
	wg.Add(1)
	//启用新线程读取服务器返回值
	go receiveResult(conn)

	wg.Wait()
	//模拟发送数据
	for i := 0; i < 10;i++{
		data := make(map[string]interface{})
		data["auth_key"] = AuthKey
		data["cmd"] = "msg"
		data["info"] = "test" + strconv.Itoa(i)
		ret,_ := json.Marshal(data)
		n,err := conn.Write(Packet(ret))
		fmt.Printf("send sucess:nb.:%d,%d,error:%v\n",i+1,n,err)
		time.Sleep(3*time.Second)
	}


}

func receiveResult(conn *net.TCPConn){

	var buffer []byte
	//处理服务端返回值
	go processRCh()

	for{
		tmpBuf := make([]byte,1024)
		n,err := conn.Read(tmpBuf)
		if err != nil{
			if err == io.EOF{
				//服务器端已关闭连接
				fmt.Println("server hs closed!")
			}else {
				fmt.Printf("other error:%v",err)
			}
			CloseChan <- true
			return
		}
		buffer = UnPacket(append(buffer, tmpBuf[:n]...), RChan)
	}
}

func processRCh(){
	for{
		select{
		case tmp_data :=<- RChan:
			//处理返回结果
			var ret map[string]interface{}
			err := json.Unmarshal(tmp_data,&ret)
			if err != nil{
				fmt.Println(err)
				continue
			}
			fmt.Printf("receive server result:%v\n",ret)
			key,exists := ret["authKey"]
			if exists{
				AuthKey = key
				fmt.Printf("connect sucess!AuthKey:%v\n",AuthKey)
				wg.Done()
			}
		case <- CloseChan:
			return
		}
	}
}