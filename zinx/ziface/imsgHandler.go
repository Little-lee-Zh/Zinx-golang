package ziface

/*
	消息管理抽象层
*/
type IMsgHandle interface {
	//调度对应的Router消息处理方法
	DoMsgHandler(request IRequest)          //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)    //将消息发送给消息队列处理
}
