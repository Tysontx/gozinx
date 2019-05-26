package znet

import "gozinx/ziface"

type Message struct {
	Id uint32
	DataLen uint32
	Data []byte
}

func NewMsgPack(id uint32, data []byte) ziface.IMessage {
	m := &Message{
		Id: id,
		DataLen: uint32(len(data)),
		Data: data,
	}
	return m
}

// getter
func (m *Message)GetMsgId() uint32 {
	return m.Id
}
func (m *Message)GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message)GetMsgData() []byte {
	return m.Data
}

// setter
func (m *Message)SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message)SetMsgLen(len uint32) {
	m.DataLen = len
}
func (m *Message)SetMsgData(data []byte) {
	m.Data = data
}
