package socket

import (
	"bytes"
	"encoding/binary"
)

var (
	Header     = "jfeige"
	HeaderLen  = 6
	ContentLen = 4
)

var (
	commands = []string{"msg", "close", "ping"}
)

//解包
func UnPacket(buffer []byte, dataChan chan []byte) []byte {
	length := len(buffer)
	if length < HeaderLen+ContentLen {
		return buffer
	}
	var i int
	for i = 0; i < len(buffer); i++ {
		if string(buffer[i:i+HeaderLen]) == Header {
			msg_length := BytesToInt(buffer[i+HeaderLen : i+HeaderLen+ContentLen])

			if length < HeaderLen+ContentLen+msg_length {
				break
			}
			dataChan <- buffer[i+HeaderLen+ContentLen : i+HeaderLen+ContentLen+msg_length]
			i += HeaderLen + ContentLen + msg_length - 1
		}
	}
	if i == length {
		return make([]byte, 0)
	}

	return buffer[i:]
}

//打包
func Packet(msg []byte) []byte {
	ret := append(append([]byte(Header), IntToBytes(len(msg))...), msg...)
	return ret

}

func IntToBytes(n int) []byte {
	x := int32(n)
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, x)
	return buffer.Bytes()
}

func BytesToInt(data []byte) int {
	var tmp int32
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &tmp)
	return int(tmp)
}
