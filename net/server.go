package net

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
	Router ziface.IRouter
}

// 初始化 New　方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		IPVersion: "tcp4",
		Ip: config.GlobalObject.Host,
		Port: config.GlobalObject.Port,
		Name: config.GlobalObject.Name,
		Router: nil,
	}
	return s
}

// 启动服务器
// 原生 socket 编程
func (s *Server)Start() {
	fmt.Printf("[start] Server Listener at IP:%s, Port:%d, is starting...\n", s.Ip, s.Port)
	// 创建套接字
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.ResolveTCPAddr err:", err)
		return
	}
	// 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("net.ListenTCP err:", err)
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
			dealConn := NewConnection(conn, cid, s.Router)
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
// 针对　type HandleFunc func(*net.TCPConn, []byte, int) error
func CallBackBusi(request ziface.IRequest) error {
	// 回显业务
	fmt.Println("[conn Handle] CallBack...")
	conn := request.GetConnection().GetTCPConnection()
	data := request.GetData()
	cnt := request.GetDataLen()
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back err")
		return err
	}
	return nil
}

// 关闭服务器
func (s *Server)Stop() {
	// 将服务器资源进行回收
}

// 运行服务器
func (s *Server)Serve() {
	// 启动 server 的监听功能
	s.Start() // 不希望永久阻塞
	select {} // 阻塞
}

// 添加路由
func (s *Server)AddRouter(router ziface.IRouter) {
		s.Router = router
}
