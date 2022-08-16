package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	//"zinx-Golang/zinx/utils"
	"zinx-Golang/zinx/utils"
	"zinx-Golang/zinx/ziface"
)

//链接模块
type Connection struct {
	//当前conn隶属于哪个server
	TcpServer ziface.IServer
	//当前链接的socket TCP套接字
	Conn net.TCPConn
	//链接的ID
	ConnID uint32
	//当前链接的状态
	isClosed bool
	// //当前链接所绑定的处理业务方法API
	// handleAPI ziface.HandleFunc
	//告知当前链接已经退出的channel
	ExitChan chan bool

	//无缓冲管道，用于读写goroutine之间的通信
	msgChan chan []byte
	//消息的管理masgID和对应的处理业务
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

//初始化链接模块得 方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       *conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//链接得读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("[Reader is exit]", "connID = ", c.ConnID, "remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端得数据到buf中
		// buf := make([]byte, utils.GlobalObject.MaxPacketSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }
		//创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流8字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head err:", err)
			break
		}

		//拆包得到msgID和msglen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}
		//根据datalen再次读取，放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data err:", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由中找到注册绑定的Conn对应的router调用
			//根据绑定好的MsgID找到对应处理的api业务
			go c.MsgHandler.DoMsgHandler(&req)
		}
		// //执行注册的路由方法
		// go func(request ziface.IRequest) {
		// 	c.Router.PreHandle(request)
		// 	c.Router.Handle(request)
		// 	c.Router.PostHandle(request)
		// }(&req)
	}
}

//写消息Goroutine，用户将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())
	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan: //如果可读，代表reader已经退出，此时writer也要
			//conn已经关闭
			return
		}
	}
}

//启动链接让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start.. ConnID = ", c.ConnID)
	//启动从当前链接得读数据业务
	go c.StartReader()
	//启动从当前链接得写数据业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接后需要调用的处理业务 执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止链接结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop.. ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	//调用开发者注册的 销毁链接之前 需要执行的业务
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()
	//告知writer关闭（）read出错，调用stop，再告知writer
	c.ExitChan <- true

	//将当前连接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return &c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供sendMsg方法，将我们要发给客户端的数据，先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包MsgDatalen/msgId/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id :", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端(直接)，现在需要先发给管道
	// _, err = c.Conn.Write(binaryMsg)
	// if err != nil {
	// 	fmt.Println("Write masg id:", msgId, "error:", err)
	// 	return errors.New("conn Write error")
	// }
	c.msgChan <- binaryMsg
	return nil
}

//设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加一个属性
	c.property[key] = value
}

//获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	//读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//删除属性
	delete(c.property, key)
}
