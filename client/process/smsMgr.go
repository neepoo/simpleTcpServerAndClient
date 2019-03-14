package process

import (
	"chapterRoom/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) {
	//显示即可
	// mes.type一定是smstype
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("饭序列化失败")
	}
	// 显示

	info := fmt.Sprintf("用户id:\t%d,对大家说:\t%s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
}
