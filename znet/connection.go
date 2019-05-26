package znet

import (
	"net"
	"gozinx/ziface"
	"fmt"
	"io"
	"errors"
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
	// Router ziface.IRouter
	// 多路由
	MsgHandler ziface.IMsgHandler
	// 添加一个 Write 和 Read 的 chan
	msgChan chan []byte
	// 用来 Reader 通知 Writer conn 已经关闭
	writerExitChan chan bool
	// 当前连接属于哪一个 Server 创建
	server ziface.IServer
}

// 初始化连接方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, MsgHandler ziface.IMsgHandler) ziface.IConnection {
	c := &Connection{
		Conn: conn,
		ConnID: connId,
		isClose: false,
		// handleAPI: callback,
		// Router: router,
		MsgHandler: MsgHandler,
		msgChan: make(chan []byte), // 初始化 write 和 Read 之间通讯的 chan
		writerExitChan: make(chan bool),
		server: server,
	}
	// 当连接已经创建成功，则添加到连接管理器中
	server.GetConnMgr().Add(c)
	return c
}

func (c *Connection)StartRead(){
	// 从对端读数据
	fmt.Println("[StartReader Goroutine isStarted...]")
	defer fmt.Println("[StartRead Goroutine isStarted...]connID=", c.ConnID, "Reader is exit, remote add is =", c.GetRemoteAddr().String())
	defer c.Stop() // 关闭 socket

	for {
		//buf := make([]byte, config.GlobalObject.MaxPageageCount)
		//cnt, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err:", err)
		//	break
		//}

		// 按照封包拆包对象，读取两次
		// 创建封包拆包对象
		dp := NewDataPack()
		// 读取消息头部
		headData := make([]byte, dp.GetHandLen())
		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Println("io.ReadFull headData err:", err)
			break
		}
		// 头部 拆包
		var data []byte
		msg, err := dp.UnPack(headData)
		if msg.GetMsgLen() > 0 {
			// 有数据
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				fmt.Println("io.ReadFull body err:", err)
				break
			}
		}
		msg.SetMsgData(data)

		// 将当前一次性得到的对端请求的数据封装成一个 Request
		req := NewRequest(c, msg)

		// 把 req 交给工作池处理
		if config.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}

		// 调用用户传递进来的业务（模板设计模式）
		// go c.MsgHandler.DoMsgHandler(req)
		//go func(){
		//	c.Router.Handle(req)
		//	c.Router.PostHandle(req)
		//	c.Router.PreHandle(req)
		//}()
		//if err := c.handleAPI(req); err != nil {
		//	fmt.Println("ConnId:", c.ConnID, "Handle is err:", err)
		//	break
		//}
	}
}

// 进行写业务
func (c *Connection)StartWrite(){
	fmt.Println("[StartWriter Goroutine isStarted]...")
	defer fmt.Println("[StartWrite Goroutine isStarted...]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err:", err)
				return
			}
		case <-c.writerExitChan:
			return
		}
	}
}

// 启动连接
func (c *Connection)Start() {
	fmt.Println("Connection Starting id=", c.ConnID)
	// 先进行读业务
	go c.StartRead()
	// 进行写业务
	go c.StartWrite()
	// 调用创建连接之后，用户自定义 hook 业务
	c.server.CallOnConnStart(c)
}

// 停止连接
func (c *Connection)Stop() {
	fmt.Println("Conn Stop connId=", c.ConnID)
	if c.isClose == true {
		return
	}
	//调用销毁链接之前用户自定义的Hook函数
	c.server.CallOnConnStop(c)
	// 将当前连接从链接管理器中删除
	c.server.GetConnMgr().Remove(c.ConnID)
	// 关闭原生套接字
	c.Conn.Close()
	// 告知 writer conn 已经关闭
	c.writerExitChan <- true
	close(c.msgChan)
	close(c.writerExitChan)
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
func (c *Connection) Send(msgId uint32, msgData []byte) error {
	if c.isClose == true {
		return errors.New("Connection is closed")
	}
	// 先封包，再发送
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPack(msgId, msgData))
	if err != nil {
		fmt.Println("dp.Pack err:", err)
		return err
	}
	// 将 binaryMsg 发送给对端
	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("send buf err")
	//	return err
	//}

	// 将数据写入 msgChan 中
	c.msgChan <- binaryMsg

	return nil
}
