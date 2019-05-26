package ziface

// 消息管理抽象层
type IMsgHandler interface {
	// 调度路由，根据 msgId
	DoMsgHandler(request IRequest)
	// 添加路由到 map 中
	AddRouter(msgId uint32, router IRouter)
	// 启动工作池
	StartWorkerPool()
	// 将消息添加到工作池中（将消息发送给对应的消息队列）
	SendMsgToTaskQueue(request IRequest)
}
