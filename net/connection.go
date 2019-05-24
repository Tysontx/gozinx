package net

import (
	"net"
	"gozinx/ziface"
	"fmt"
	"gozinx/config"
)

// 具体的 TCP 连接模块
type Connection struct {
	// 当前连接的原生套接字
	Conn *net.TCPConn
	// 连接 id
	ConnID uint32
	// 连接状态
	isClose bool
	// 当前连接所绑定的业务处理方法
	// handleAPI ziface.HandleFunc
	Router ziface.IRouter
}

// 初始化连接方法
func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) ziface.IConnection {
	c := &Connection{
		Conn: conn,
		ConnID: connId,
		isClose: false,
		// handleAPI: callback,
		Router: router,
	}
	return c
}

func (c *Connection)StartRead(){
	// 从对端读数据
	fmt.Println("Reader go id starting...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote add is =", c.GetRemoteAddr().String())
	defer c.Stop() // 关闭 socket

	for {
		buf := make([]byte, config.GlobalObject.MaxPageageCount)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			break
		}
		// 将当前一次性得到的对端请求的数据封装成一个 Request
		req := NewRequest(c, buf, cnt)
		// 调用用户传递进来的业务（模板设计模式）
		go func(){
			c.Router.Handle(req)
			c.Router.PostHandle(req)
			c.Router.PreHandle(req)
		}()
		//if err := c.handleAPI(req); err != nil {
		//	fmt.Println("ConnId:", c.ConnID, "Handle is err:", err)
		//	break
		//}
	}
}

// 启动连接
func (c *Connection)Start() {
	fmt.Println("Connection Starting id=", c.ConnID)
	// 先进行读业务
	go c.StartRead()
}
// 停止连接
func (c *Connection)Stop() {
	fmt.Println("Conn Stop connId=", c.ConnID)
	if c.isClose == true {
		return
	}
	// 关闭原生套接字
	c.Conn.Close()
	c.isClose = true
}
// 获取连接 id
func (c *Connection)GetConnID() uint32 {
	return c.ConnID
}
// 获取 conn 原生套接字
func (c *Connection)GetTCPConnection() *net.TCPConn {
	return c.Conn
}
// 获取远程客户端的 IP 地址
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
// 发送数据给对端
func (c *Connection) Send(data []byte, cnt int) error {
	if _, err := c.Conn.Write(data[:cnt]); err != nil {
		fmt.Println("send buf err")
		return err
	}
	return nil
}
