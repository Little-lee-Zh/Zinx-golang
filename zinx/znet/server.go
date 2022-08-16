package znet

import (
	"fmt"
	"net"
	"zinx-Golang/zinx/utils"
	"zinx-Golang/zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//当前的Server的消息管理模块，用来绑定MsgId和对应的处理逻辑
	MsgHandler ziface.IMsgHandle
	//该server连接管理器
	ConnMgr ziface.IConnManager

	//该server创建之后自动调用Hook函数OnConnStart
	OnConnStart func(conn ziface.IConnection)
	//该server销毁之前自动调用Hook函数OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[zinx] server Name : %s, listenner at IP : %s, port : %d is starting\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[zinx] Version : %s, maxconn : %d, maxPacketsize : %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)

	go func() {
		//0.开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()
		//1.获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr err :", err)
			return
		}
		//2.监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen", s.IPVersion, "err: ", err)
			return
		}
		fmt.Println("start Zinx server, ", s.Name, "success, now Listening...")
		var cid uint32
		cid = 0
		//3.阻塞的等待客户端链，处理客户端链接业务（读写）
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接个数的判断，如果超过最大连接，那么关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many Connections maxconn:", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//将处理新连接得业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前得链接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//将一些服务器的资源状态或者一些已经开辟的链接信息进行停止或者回收
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

//路由功能给当前服务注册一个路由方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router succ!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(Iconnection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(connection)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(connection)
	}
}
