package main

import (
	"chapterRoom/server/model"
	"fmt"
	"net"
	"time"
)

// 处理和客户端的通信
//net.Conn引用类型，不要取地址
func process(conn net.Conn) {
	// 这里需要延时关闭conn
	defer conn.Close()

	// 调用主控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("process2,err=", err)
		return
	}

}

// 编写一个函数 完成对userDo的初始化任务
func initUserDao() {
	// 这里的pool本身就是一个全局变量
	model.MyUserDao = model.NewUserDao(pool)
}
func main() {
	// 当服务器启动时，我们就去初始化我们redis的连接池
	initPool("localhost:6379", 16, 0, time.Second*300)

	// 初始化userDao实例
	initUserDao()

	// 监听，等待初始化连接
	fmt.Println("服务器[新的结构]8889端口")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net listen err=", err)
	}
	// 一旦监听成功，等待客户端连接服务器
	for {

		conn, acceptErr := listen.Accept()
		fmt.Println("conn=", conn)
		if acceptErr != nil {
			fmt.Println("listen.Accept ERR=", acceptErr)
		}
		// 一旦连接成功，启动一个协程和客户端保持通信...
		go process(conn)
		// TODO 需要处理不同的消息(ServerProcessMes)
	}
}
