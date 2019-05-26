package znet

import (
	"gozinx/ziface"
	"bytes"
	"encoding/binary"
	"fmt"
)

type DataPack struct {
}

// 封包方法
func (dp *DataPack)Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})

	err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()) // 写入 MsgLen
	if err != nil {
		fmt.Println("binary.Write MsgLen err:", err)
		return nil, err
	}

	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()) // 写入 MsgId
	if err != nil {
		fmt.Println("binary.Write MsgId err:", err)
		return nil, err
	}

	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgData()) // 写入 MsgData
	if err != nil {
		fmt.Println("binary.Write MsgData err:", err)
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

// 解包方法
func (dp *DataPack)UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 解包分两步，第一步：读到 MsgLen、MsgId 放到 Message 中；第二步：根据 MsgLen 的长度，再读 MsgData
	messageHead := &Message{}
	// 创建一个读取二进制流的 Reader
	dataBuff := bytes.NewReader(binaryData)
	// 读出 MsgLen
	err := binary.Read(dataBuff, binary.LittleEndian, &messageHead.DataLen)
	if err != nil {
		fmt.Println("binary.Read MsgLen err:", err)
		return nil, err
	}
	// 读出 Id
	err = binary.Read(dataBuff, binary.LittleEndian, &messageHead.Id)
	if err != nil {
		fmt.Println("binary.Read id err:", err)
		return nil, err
	}
	return messageHead, nil
}

// 获取头部长度
func (dp *DataPack)GetHandLen() uint32 {
	return 8
}

// 初始化 new 方法
func NewDataPack() ziface.IDataPack {
	dp := &DataPack{}
	return dp
}
