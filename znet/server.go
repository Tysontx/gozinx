package znet

import (
	"gozinx/ziface"
	"fmt"
	"net"
	"gozinx/config"
)

type Server struct {
	// 服务器 ip
	IPVersion string
	Ip string
	// 端口
	Port int
	// 服务器名称
	Name string
	// 当前 Server 由用户绑定的回调 router,也就是 Server 注册的链接对应的处理业务
	// Router ziface.IRouter
	// 多路由
	MsgHandler ziface.IMsgHandler
	// 连接管理
	ConnMgr ziface.IConnMannger
	// server 连接之后自动调用的 Hook 函数
	OnConnStart func(conn ziface.IConnection)
	// server 连接销毁之前调用的 Hook 函数
	OnConnStop func(conn ziface.IConnection)
}

// 初始化 New　方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		IPVersion: "tcp4",
		Ip: config.GlobalObject.Host,
		Port: config.GlobalObject.Port,
		Name: config.GlobalObject.Name,
		MsgHandler: NewMsgHandler(),
		ConnMgr: NewConnMannger(),
	}
	return s
}

// 启动服务器
// 原生 socket 编程
func (s *Server)Start() {
	fmt.Printf("[start] Server Listener at IP:%s, Port:%d, is starting...\n", s.Ip, s.Port)
	// 开启工作池
	s.MsgHandler.StartWorkerPool()
	// 创建套接字
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("znet.ResolveTCPAddr err:", err)
		return
	}
	// 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("znet.ListenTCP err:", err)
		return
	}
	var cid uint32 = 0 // 生成 id 累加器
	go func(){
		// 阻塞等待客户端连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("listener.Accept err:", err)
				continue
			}
			// 判断当前连接数是否大于最大连接数
			if s.GetConnMgr().Len() >= config.GlobalObject.MaxConn {
				fmt.Println("Too Many Connection, MaxConn=", config.GlobalObject.MaxConn)
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			// 运行到这里就已经有客户端和服务端建立了连接
			go dealConn.Start()
			//go func(){
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("conn.Read err:", err)
			//			break
			//		}
			//		fmt.Printf("recv client buf:%s, length:%d\n", string(buf[:cnt]), cnt)
			//		// 回显
			//		_, err = conn.Write(buf[:cnt])
			//		if err != nil {
			//			fmt.Println("conn.Write err:", err)
			//			continue
			//		}
			//	}
			//}()
		}
	}()
}

// 定义一个具体的回显业务
// 针对　type HandleFunc func(*znet.TCPConn, []byte, int) error
//func CallBackBusi(request ziface.IRequest) error {
//	// 回显业务
//	fmt.Println("[conn Handle] CallBack...")
//	conn := request.GetConnection().GetTCPConnection()
//	data := request.GetData()
//	cnt := request.GetDataLen()
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back err")
//		return err
//	}
//	return nil
//}

// 关闭服务器
func (s *Server)Stop() {
	// 服务器关闭时，应将所有连接都关闭
	s.GetConnMgr().ClearConn()
}

// 运行服务器
func (s *Server)Serve() {
	// 启动 server 的监听功能
	s.Start() // 不希望永久阻塞
	select {} // 阻塞
}

// 添加路由
func (s *Server)AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router succ! msgId=", msgId)
}

// 提供一个得到连接管理的方法
func (s *Server)GetConnMgr() ziface.IConnMannger {
	return s.ConnMgr
}

// 注册创建连接之后调用的 Hook 方法
func (s *Server)AddOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
// 注册销毁连接之前调用的 Hook 方法
func (s *Server)AddOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}
// 调用创建链接之后的 HOOK 函数的方法
func (s *Server)CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart...")
		s.OnConnStart(conn)
	}
}
// 调用销毁连接之前的 HOOK 函数的方法
func (s *Server)CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop...")
		s.OnConnStop(conn)
	}
}
