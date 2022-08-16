package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由功能给当前服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)
	GetConnMgr() IConnManager //获取当前server的连接管理器

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(connection IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(connection IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(connection IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(connection IConnection)
}
