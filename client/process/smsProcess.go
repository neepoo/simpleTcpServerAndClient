package process

import (
	"chapterRoom/client/utils"
	"chapterRoom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

// 群发消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	// 1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	// 创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	// 将smsMes序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json marshal smsMes err,err=", err)
		return
	}
	mes.Data = string(data)
	// 序列化mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json marshal mes err,err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("send mes err=", err)
		return
	}
	return
}
