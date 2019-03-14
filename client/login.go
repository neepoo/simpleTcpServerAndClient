package main

import (
	"chapterRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 写一个函数，完成登陆校验
func login(userId int, userPwd string) (err error) {
	// 下一步就要开始定协议...
	//fmt.Printf("userId=%d, userPwd=%s\n", userId, userPwd)
	//1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 2. 准备数据发送给server
	var mes message.Message
	// 确定发送类型
	mes.Type = message.LoginMesType

	// 3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4.将loginMes序列化给到mes.data
	dataUnit, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5.将dataUnit赋给mse.Data字段
	mes.Data = string(dataUnit)

	// 6.将mse序列化
	dataSend, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 7.到这个时候，dataSend就是我们要发送的消息
	// 7.1但是我们事先发送长度，在发送dataSend
	// 把dataSend长度->一个byte切片
	var pkgLen uint32
	pkgLen = uint32(len(dataSend))
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buffer[0:4])
	if n != 4 || err != nil {
		fmt.Println("长度发送失败，err=", err)
		return
	}
	fmt.Printf("客户端发送数据长度成功, datalen=%d,内容=%s", pkgLen, string(dataSend))

	// 发送消息本身
	_, err = conn.Write(dataSend)
	if err != nil {
		fmt.Println("conn.Write(dataSend)，err=", err)
		return
	}
	// 还需要处理server发过来的消息
	mes, err = readPkg(conn)
	if err != nil {
		fmt.Println(" readPkg(conn) fail, err=", err)
	}
	// 将mes的Data部分反序列化为LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登陆成功")
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	fmt.Println("loginResMes=", loginResMes)
	return
}

// TODO 320P
