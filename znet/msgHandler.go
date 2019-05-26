package znet

import (
	"gozinx/ziface"
	"fmt"
	"gozinx/config"
)

type MsgHandler struct {
	// 存放路由的 map
	Apis map[uint32] ziface.IRouter
	// 负责 worker 取任务的消息队列，一个 worker 对应一个任务队列
	TaskQueue []chan ziface.IRequest
	// worker 工作池的 worker 数量
	WorkerPoolSize uint32
}

// 调度路由，根据 msgId
func (mh *MsgHandler)DoMsgHandler(request ziface.IRequest) {
	// 从 request 获取 msgId
	msgId := request.GetMsg().GetMsgId()
	router, ok := mh.Apis[msgId]
	if !ok {
		fmt.Println("api msgId=", msgId, " Not found")
		return
	}
	// 根据 msgId 调用对应的 router
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}
// 添加路由到 map 中
func (mh *MsgHandler)AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断 msgId 是否存在
	if _, ok := mh.Apis[msgId]; ok {
		// msgId 已存在
		fmt.Println("repeat Api msgId=", msgId)
		return
	}
	// 添加 msgId 和 router 的对应关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId=", msgId, " succ!")
}

// 启动工作池（在整个 server 服务中，只启动一次）
func (mh *MsgHandler)StartWorkerPool() {
	fmt.Println("WorkerPool Is Start...")
	// 根据 WorkerPoolSize 创建 worker goroutine
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 给 channel 进行开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, config.GlobalObject.MaxWorkerTaskLen)
		// 启动一个 worker，阻塞等待消息从对应的管道中进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}
// 将消息添加到工作池中（将消息发送给对应的消息队列）
func (mh *MsgHandler)SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给 worker，确定当前 request 到底给哪个 worker 来处理
	workId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	// 直接将 request 发送给对应的 worker 的 taskqueues
	mh.TaskQueue[workId] <- request
}

// 一个 worker 真正处理业务函数的 goroutine
func (mh *MsgHandler)startOneWorker(workId int, taskQueue chan ziface.IRequest) {
	fmt.Println("workId=", workId, "is starting...")
	// 不断的从对应的管道等待数据
	for {
		select {
		case msg := <-taskQueue:
			mh.DoMsgHandler(msg)
		}
	}
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: config.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan ziface.IRequest, config.GlobalObject.WorkerPoolSize), // 切片初始化
	}
}
