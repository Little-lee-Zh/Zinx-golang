package apis

import (
	"fmt"
	"zinx-Golang/mmo_game_zinx/core"
	"zinx-Golang/mmo_game_zinx/pb"
	"zinx-Golang/zinx/ziface"
	"zinx-Golang/zinx/znet"

	"github.com/golang/protobuf/proto"
)

//玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (*MoveApi) Handle(request ziface.IRequest) {
	//1. 将客户端传来的proto协议解码
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error ", err)
		return
	}

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pid中获取
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		request.GetConnection().Stop()
		return
	}

	fmt.Printf("user pid = %d , move(%f,%f,%f,%f)", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	//3. 根据pid得到player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4. 让player对象发起移动位置信息广播
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
