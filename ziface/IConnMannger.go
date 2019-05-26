package ziface

// 连接管理抽象层
type IConnMannger interface {
	// 添加连接
	Add(connection IConnection)
	// 删除连接
	Remove(connId uint32)
	// 根据 connId 得到 conn
	Get(connId uint32) (IConnection, error)
	// 得到目前服务器的连接总个数
	Len() uint32
	// 清空全部连接方法
	ClearConn()
}
