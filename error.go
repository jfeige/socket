package socket

import (
	"errors"
	"fmt"
)

var (
	//错误连接过多
	ErrTooManyErrConns = errors.New("error:too many listener error on server!")

	//连接数超过了服务器限制
	ErrTooManyConns = errors.New("error:too many connection on server!")
)


func lockArgsErr(args string)error{

	return fmt.Errorf("error:lock parameter:%s",args)

}