package process

import (
	"chapterRoom/client/model"
	"chapterRoom/common/message"
	"fmt"
)

// 客户端要维护的map先搞出来
var onlineUsers = make(map[int]*message.User, 10)
var CurUser model.CurUser // 在用户登录成功后完成对curUser的初始化

// 在客户端显示当前在线用户
func outputOnlineUser() {
	// 遍历onlineUsers即可
	fmt.Println("当前用户在线列表： ")
	for id, user := range onlineUsers {
		// 如果不想显示自己就过滤
		fmt.Println("用户id:", id, "登陆成功\t", "用户状态=", user.UserStatus)
	}
}

// 编写一个方法处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	// 适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 说明原来没有，创建user
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}
