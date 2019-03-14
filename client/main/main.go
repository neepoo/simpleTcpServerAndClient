package main

import (
	"chapterRoom/client/process"
	"fmt"
)

// 定义两个变量，一个表示用户id， 一个表示密码
var userId int
var userPwd string
var userName string

func main() {

	// 接受用户的选择
	var key int
	for true {
		fmt.Println("---------------欢迎登陆多人聊天室---------")
		fmt.Println("---------------1. 登陆聊天室---------")
		fmt.Println("---------------2. 注册用户---------")
		fmt.Println("---------------3. 退出系统---------")
		fmt.Println("---------------请选择(1-3)---------")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户的ID：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			// 1.创建一个UserProcess实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(昵称)：")
			fmt.Scanf("%s\n", &userName)

			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")

		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}

}
