package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	//模拟的服务器
	//1.创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	//创建一个go承载负责从客户端处理业务
	go func() {
		//2.从客户端读取数据拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("listenner accept err :", err)
			}

			go func(conn net.Conn) {
				//处理客户端的请求
				//拆包的过程
				dp := NewDataPack()
				for {
					//1.第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err")
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2.第二次根据head中datalen再读取data
						msg := msgHead.(*Message) //类型断言
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen长度在读
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err", err)
							return
						}
						//一个完整的消息已经读取完毕
						fmt.Println("Recv MsgID:", msg.Id, "datalen:", msg.DataLen, "data:", string(msg.Data))
					}

				}
			}(conn)
		}
	}()
	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	//创建一个封包对象dp
	dp := NewDataPack()
	//模拟毡包过程，封装两个msg一起发送
	//封装第一个msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}
	//封装第一个msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', 'o', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}
	//粘在一起
	sendData1 = append(sendData1, sendData2...)
	//一次性发送给服务器
	conn.Write(sendData1)
	//客户端阻塞
	select {}
}
