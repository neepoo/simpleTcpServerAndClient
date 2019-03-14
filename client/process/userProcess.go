package process

import (
	"chapterRoom/client/utils"
	"chapterRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	// 暂时不需要任何字段

}

// 注册方法
func (this *UserProcess) Register(userId int, userPwd string, userName string) {
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
	mes.Type = message.RegisterMesType

	// 3.创建一个RegisterMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4.先把registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//5.
	mes.Data = string(data)

	// 6.将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 7.
	tf := utils.Transfer{
		Conn: conn,
	}

	//8.发送数据给服务器端
	err = tf.WritePkg([]byte(data))
	if err != nil {
		fmt.Println("when register write error err=", err)
	}

	// 接受信息
	mes, err = tf.ReadPkg() // mes就是
	if err != nil {
		fmt.Println(" readPkg(conn) fail, err=", err)
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功， 可以登陆")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}

// 关联一个登陆的方法
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
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
	//"{\"userId\":100,\"userPwd\":\"123456\",\"username\":\"scoot\"}"
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
	fmt.Printf("客户端发送数据长度成功, datalen=%d,内容=%s\n\n", pkgLen, string(dataSend))

	// 发送消息本身
	_, err = conn.Write(dataSend)
	if err != nil {
		fmt.Println("conn.Write(dataSend)，err=", err)
		return
	}

	// 还需要处理server发过来的消息
	// 创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn, // 自身带的
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println(" readPkg(conn) fail, err=", err)
	}
	// 将mes的Data部分反序列化为LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// 初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		// 进入二级菜单
		// 显示我们的登陆成功的菜单
		// 这里我们还需要再客户端启动一个协程
		// 这个协程保持和服务器端的通讯。
		// 如果服务器有数据推送给我们，则接收并显示在客户端终端

		// 现在可以显示当前在线用户列表
		// 遍历loginResMes.UsersId
		fmt.Println("当前在线用户列表如下: ")
		for _, v := range loginResMes.UsersId {
			if v == userId {
				continue
			}
			fmt.Printf("用户id:\t%v\n", v)
			// 完成客户端的onlineUsers初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		go processServerMes(conn)
		for {
			// 循环显示菜单
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
