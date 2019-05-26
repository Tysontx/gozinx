package ziface

type IServer interface {
	// 启动服务器
	Start()
	// 关闭服务器
	Stop()
	// 运行服务器
	Serve()
	// 添加路由
	AddRouter(msgId uint32, router IRouter)
	// 提供一个得到连接管理的方法
	GetConnMgr() IConnMannger
	// 注册创建连接之后调用的 Hook 方法
	AddOnConnStart(hookFunc func(conn IConnection))
	// 注册销毁连接之前调用的 Hook 方法
	AddOnConnStop(hookFunc func(conn IConnection))
	// 调用创建链接之后的 HOOK 函数的方法
	CallOnConnStart(conn IConnection)
	// 调用销毁连接之前的 HOOK 函数的方法
	CallOnConnStop(conn IConnection)
}
