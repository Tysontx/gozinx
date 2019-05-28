package ziface

import "net"

// 抽象连接层
type IConnection interface {
	// 启动连接
	Start()
	// 停止连接
	Stop()
	// 获取连接 id
	GetConnID() uint32
	// 获取 conn 原生套接字
	GetTCPConnection() *net.TCPConn
	// 获取远程客户端的 IP 地址
	GetRemoteAddr() net.Addr
	// 发送数据给对端
	Send(msgId uint32, msgData []byte) error
	// 设置属性
	SetProperty(key string, value interface{})
	// 获取属性
	GetProperty(key string) (interface{}, error)
	// 删除属性
	RemoveProperty(key string)
}

// 业务处理方法　抽象定义
// type HandleFunc func(request IRequest) error

//func CallBackBusi(request ziface.IRequest) error {
