package main

import (
	"fmt"
	protocol "zinx-Golang/myDemo/protobufDemo/pb"

	"github.com/golang/protobuf/proto"
)

func main() {
	person := &protocol.Person{
		Name:   "Aceld",
		Age:    16,
		Emails: []string{"https://legacy.gitbook.com/@aceld", "https://github.com/aceld"},
		Phones: []*protocol.PhoneNumber{
			&protocol.PhoneNumber{
				Number: "13113111311",
				Type:   protocol.PhoneType_MOBILE,
			},
			&protocol.PhoneNumber{
				Number: "14141444144",
				Type:   protocol.PhoneType_HOME,
			},
			&protocol.PhoneNumber{
				Number: "19191919191",
				Type:   protocol.PhoneType_WORK,
			},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err:", err)
	}
	fmt.Println("源数据", data)
	newdata := &protocol.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		fmt.Println("unmarshal err:", err)
	}

	fmt.Println(newdata)

}
