package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx-Golang/zinx/znet"
)

/*
	模拟客户端
*/
func main() {

	fmt.Println("Client1 Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		//发封包message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx Client1 Test Message")))
		if err != nil {
			fmt.Println("pack error:", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		binaryData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(binaryData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetMsgLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
