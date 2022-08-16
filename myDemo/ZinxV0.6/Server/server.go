package main

import (
	"fmt"
	"zinx-Golang/zinx/ziface"
	"zinx-Golang/zinx/znet"
)

//基于Zinx框架来开发的服务器端应用程序

//ping test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再写回ping
	fmt.Println("recv from client:msgID=", request.GetMsgId(), "data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

//Hello zinx自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再写回ping
	fmt.Println("recv from client:msgID=", request.GetMsgId(), "data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello welcome to zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.6]")
	//2.给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloZinxRouter{})
	//3.启动server
	s.Serve()
}
