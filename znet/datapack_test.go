package znet

import (
	"testing"
	"fmt"
	"net"
	"io"
)

// 单元测试 DataPack
func TestDataPack(t *testing.T) {
	fmt.Println("start testing datapack...")

	// 模拟一个 server
	// 1、创建监听
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// 2、阻塞等待对端连接
	go func(){
		for {
			conn, err := listener.Accept()
			fmt.Println("listener to client...")
			if err != nil {
				fmt.Println("listener.Accept err:", err)
				return
			}
			go func(conn *net.Conn){
				// 3、读写业务
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHandLen())
					_, err := io.ReadFull(*conn, headData) // 如果缓存没有读满，则会阻塞
					if err != nil {
						fmt.Println("io.ReadFull Head err:", err)
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("dp.Unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msgHead.GetMsgLen())
						_, err := io.ReadFull(*conn, msg.Data)
						if err != nil {
							fmt.Println("io.ReadFull msgData err:", err)
							return
						}
						fmt.Println("result MsgId=", msg.Id, " MsgLen=", msg.DataLen, " MsgData=", string(msg.Data))
					}
				}
			}(&conn)
		}
	}()

	// 模拟一个 client
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("net.Dail err:", err)
		return
	}
	// 封包过程
	dp := NewDataPack()
	msg1 := &Message{
		Id: 1,
		DataLen: 6,
		Data: []byte{'h', 'e', 'l', 'l', 'o', ','},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("dp.Pack err:", err)
		return
	}
	msg2 := &Message{
		Id: 2,
		DataLen: 6,
		Data: []byte{'w', 'o', 'r', 'l', 'd', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("dp.Pack err:", err)
		return
	}
	// 把封包过后的数据追加到一起，造成粘包现象
	sendData1 = append(sendData1, sendData2...)
	// 发送数据
	conn.Write(sendData1)
	select {} // 阻塞不退出
}
