package ziface

type IDataPack interface {
	// 封包方法
	Pack(msg IMessage) ([]byte, error)
	// 解包方法
	UnPack(binaryData []byte) (IMessage, error)
	// 获取头部长度
	GetHandLen() uint32
}
