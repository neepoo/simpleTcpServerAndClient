package process

import (
	"chapterRoom/client/utils"
	"chapterRoom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// 显示登陆成功的页面
// 保持通讯
// 读取服务器消息，显示在页面上

// 显示登陆成功后的界面
func ShowMenu() {
	fmt.Println("-------恭喜xxx登陆成功--------")
	fmt.Println("-------1. 显示在线用户列表--------")
	fmt.Println("-------2. 发送消息--------")
	fmt.Println("-------3. 信息列表--------")
	fmt.Println("-------4. 退出系统--------")
	fmt.Println("-------请选择(1-4)------")
	var key int
	var content string

	// 因为我们总会使用到SmsProcess实例
	// 所以定义在外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("请输入你想对大家说的话")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入有误")
		fmt.Printf("key=%T,key=%v", key, key)
	}

}

// 在main中拿到连接，保持连接
func processServerMes(conn net.Conn) {

	// 创建transfer实例,不停的读取服务器发送的消息
	tf := utils.Transfer{
		Conn: conn,
	}
	for {
		mes, err := tf.ReadPkg()

		if err != nil {
			fmt.Println("f.ReadPkg err=", err)
			return
		}
		// 读取消息，又是下一步逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType: // 有人上线了
			//1. 取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			if err != nil {
				fmt.Println("反序列化通知用户上线失败")
			}
			//2. 把这个人加入到客户端维护的map[int]User中
			updateUserStatus(&notifyUserStatusMes)

		case message.SmsMesType: //有人群发消息，需要转发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("该类型，我暂时无法识别")
		}
	}
}
