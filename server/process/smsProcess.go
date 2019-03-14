package process2

import (
	"chapterRoom/common/message"
	"chapterRoom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

//写方法转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	// 先取出内容
	var smsMsg message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMsg)
	if err != nil {
		fmt.Println("反序列化失败，err=", err)
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("err=", err)
	}

	// 遍历当前在线用户，把消息转发出去
	for userId, userConn := range userMgr.GetAllOnlineUsers() {
		// 避免自己给自己发消息
		if userId == smsMsg.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, userConn.Conn)
	}
}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("失败")
	}
}
