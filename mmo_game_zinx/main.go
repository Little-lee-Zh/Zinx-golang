package main

import (
	"fmt"
	"zinx-Golang/mmo_game_zinx/apis"
	"zinx-Golang/mmo_game_zinx/core"
	"zinx-Golang/zinx/ziface"
	"zinx-Golang/zinx/znet"
)

//当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)
	//同步当前的PlayerID给客户端，走MsgID:1 消息
	player.SyncPid()
	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()

	//将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该连接绑定一个Pid玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	//同步周边玩家，告知当前玩家已经上线，广播当前玩家的位置
	player.SyncSurrounding()

	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

//当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	//获取当前连接的Pid属性
	pid, _ := conn.GetProperty("pid")

	//根据pid获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线业务
	if pid != nil {
		player.OffOnline()
	}

	fmt.Println("====> Player ", pid, " left =====")

}

func main() {
	//创建Zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnecionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	//启动服务
	s.Serve()
}
