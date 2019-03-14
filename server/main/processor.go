package main

import (
	"chapterRoom/common/message"
	"chapterRoom/server/process"
	"chapterRoom/server/utils"
	"fmt"
	"io"
	"net"
)

// 创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

//依据mes.Type调用不同的函数
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	// 先看看是否能够接受到客户端发送过来的群聊消息
	fmt.Println("mes=", mes)
	switch mes.Type {
	case message.LoginMesType:
		//处理登录逻辑
		//创建一个UserProcess实例
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 处理注册
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("不是 注册，登陆无法处理")
	}
	return err
}

// 会生成mes
func (this *Processor) process2() (err error) {
	// 读客户端发送的信息
	for {
		// 这里我们将读取数据包，直接封装为一个函数readPkg(),返回message,err
		// 创建tramfer实例，完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，退出该协程")
				return err
			} else {
				fmt.Println("readPkg err", err)
				return err
			}

		}
		//fmt.Println("mes=", mes)
		// 判断登陆是否合法，返回LoginResMes
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}

	}
}
