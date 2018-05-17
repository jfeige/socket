package socket

import "errors"

var (
	//错误连接过多
	ErrTooManyErrConns = errors.New("too many listener error on server!")

	//连接数超过了服务器限制
	ErrTooManyConns = errors.New("too many connection on server!")
)
