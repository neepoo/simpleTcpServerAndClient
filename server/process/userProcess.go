package process2

import (
	"chapterRoom/common/message"
	"chapterRoom/server/model"
	"chapterRoom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//字段？
	Conn net.Conn
	// 需要增加一个字段，表示该conn是呢个用户的
	UserId int
}

//编写通知所有在线的用户的方法
// userId该用户通知其他的在线用户，他上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历UserMgr.onlineUsers,然后一个一个的发送，NotifyUserStatusMes
	for id, up := range userMgr.GetAllOnlineUsers() {
		if id == userId {
			// 过滤掉自己
			continue
		}
		// 开始通知【单独写一个方法】
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	//传进来的是在线用户的id
	// 组装我们的NotifyUserStatusMes消息
	// 大的message
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("序列化notifyUserStatusMes时发生了错误,err=", err)
	}
	mes.Data = string(data)
	// 将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化mes时发生了错误,err=", err)
	}
	// 创建Transfer实例发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline出错,err=", err)
	}
}

// 编写一个函数serverProcessLogin函数，处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码
	// 1.先从mes中取出Data，并且直接反序列化为LoginMes
	var loginMes message.LoginMes
	// 先把string-->byte[]
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	fmt.Println("&loginMes=", loginMes)
	if err != nil {
		fmt.Println("反序列化loginMes失败,ERR=", err)
		return
	}

	// 1.先此声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2.声明一个LoginResMes
	var loginResMes message.LoginResMes

	// 我们需要到redis验证
	// 使用model.UserDao到redis验证
	// 取用户，然后判断密码
	_, err = model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROE_USER_NOTEXISTS {

			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROE_USER_PWD {
			loginResMes.Code = 300
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}

	} else {
		loginResMes.Code = 200
		fmt.Println(loginMes, "登陆成功")
		// 这里用户已经登陆成功
		// 把这个用户放入到UserMgr中
		// 将登陆成功的用户的UserId赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnLineUsers(this)
		// 通知其他在线用户，我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		// 将当前用户的id放入loginResMes的UsersId
		// 遍历UserMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
	}

	// 3.序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("序列化失败，err=", err)
		return
	}

	// 4.将data付给resMes
	resMes.Data = string(data)

	// 对resMes序列化，真被发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化resMes失败，err=", err)
		return
	}

	//6. 发送data，将其封装到writePkg函数中
	// 因为使用了分层模式，我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	fmt.Println("在处理注册的函数中，mes=", mes)
	var registerMes message.RegisterMes
	// 先把string-->byte[]
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("反序列化loginMes失败,ERR=", err)
		return
	}
	// 1.先此声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	var registerResMes message.RegisterResMes

	// 到数据库完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	// TODO 上面会抛错，走不到这里
	fmt.Println("&registerMes.User=", &registerMes.User)

	if err != nil {
		if err == model.ERROE_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	// 反序列化res
	// 对resMes序列化，发送
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("序列化registerResMes失败，err=", err)
		return
	}

	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("序列化resMes失败，err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
