package znet

import "zinx-Golang/zinx/ziface"

//实现router时先嵌入Baserouter基类，然后根据需要对基类方法进行重写
type BaseRouter struct{}

//BaseRouter的方法都为空，有的Router不希望有Prehandle,PostHandle这两个业务
//所以Router全部继承BaseRouter好处，不需要实现Prehandle,PostHandle
//在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//在处理conn业务的主方法hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

//在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
