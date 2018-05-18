package socket

import (
	"fmt"
	"os"
	"testing"
	"errors"
	"encoding/json"
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

			go client.processRch()

			go client.processWch()

		}
	}
}




//处理接收的数据
func (this *Client)processRch() {
	for {
		select {
		case tmp_data, ok := <-this.RChan:
			if !ok {
				//通道已关闭
				break
			}
			data,err := this.resolveReceive(tmp_data)
			if err != nil{
				this.disconnect()
			}
			//模拟具体业务，给客户端返回值
			fmt.Printf("receive client data:%v\n",data)
			//如需要给客户端返回值，则往各自的WChan通道里，写入返回值即可
			info,ok := data["info"]
			if ok{
				data := make(map[string]interface{})
				data["info"] = info
				ret,_ := json.Marshal(data)
				this.WChan <- Packet(ret)
			}

		}
	}
}


func (this *Client)resolveReceive(tmp_data []byte)(data map[string]interface{},err error){
	data = make(map[string]interface{})
	err = json.Unmarshal(tmp_data,&data)
	if err != nil{
		return nil,err
	}
	//进行一些必需的判断
	auth_key,exists := data["auth_key"]
	if !exists{
		return nil,lockArgsErr("auth_key")
	}
	if auth_key != this.AuthKey{
		return nil,errors.New("wrong key!")
	}
	cmd,exists := data["cmd"]
	if !exists{
		return nil,lockArgsErr("cmd")
	}
	if !InArray(cmd,commands){
		return nil,errors.New("wrong cmd!")
	}
	return
}

//发送数据到客户端
func (this *Client)processWch() {
	for {
		select {
		case tmp_data, ok := <-this.WChan:
			if !ok {
				break
			}
			this.Conn.Write(tmp_data)
		}
	}
}
