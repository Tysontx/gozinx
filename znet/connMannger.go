package znet

import (
	"gozinx/ziface"
	"sync"
	"fmt"
	"errors"
)

type ConnMannger struct {
	// 管理全部的连接
	connections map[uint32]ziface.IConnection
	connLock sync.RWMutex
}

func NewConnMannger() ziface.IConnMannger {
	return &ConnMannger{
		connections: make(map[uint32] ziface.IConnection),
	}
}

// 添加连接
func (cm *ConnMannger)Add(connection ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[connection.GetConnID()] = connection
	fmt.Println("Add connId=", connection.GetConnID(), " to mannger succ!")
}
// 删除连接
func (cm *ConnMannger)Remove(connId uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, connId)
	fmt.Println("Delete connId=", connId, " to mannger succ!")
}
// 根据 connId 得到 conn
func (cm *ConnMannger)Get(connId uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	cm.connLock.RUnlock()
	if conn, ok := cm.connections[connId]; ok {
		// 找到了
		return conn, nil
	} else {
		return nil, errors.New("Connection Not Found!")
	}
}
// 得到目前服务器的连接总个数
func (cm *ConnMannger)Len() uint32 {
	return uint32(len(cm.connections))
}
// 清空全部连接方法
func (cm *ConnMannger)ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 遍历删除
	for connID, conn := range cm.connections {
		// 讲全部的 conn 关闭
		conn.Stop()
		// 删除全部
		delete(cm.connections, connID)
	}
	fmt.Println("Clear All Connection Succ! Connnection num= ", cm.Len())
}
