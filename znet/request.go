package znet

import "gozinx/ziface"

type Request struct {
	// 链接信息
	conn ziface.IConnection
	msg ziface.IMessage
}

func NewRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	req := &Request{
		conn: conn,
		msg: msg,
	}
	return req
}

// 得到当前请求的链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetMsg() ziface.IMessage {
	return r.msg
}


