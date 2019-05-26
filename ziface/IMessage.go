package ziface

// 将请求的一个消息封装到 message 中
type IMessage interface {
	// getter
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	// setter
	SetMsgId(id uint32)
	SetMsgLen(len uint32)
	SetMsgData(data []byte)
}
