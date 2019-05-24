package ziface

// 抽象 IRequest 一次性请求的数据封装
type IRequest interface {
	// 得到当前请求的链接
	GetConnection() IConnection
	// 得到链接的数据
	GetData() []byte
	// 得到链接的长度
	GetDataLen() int
}
