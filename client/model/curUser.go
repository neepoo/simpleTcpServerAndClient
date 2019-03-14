package model

import (
	"chapterRoom/common/message"
	"net"
)

// 因为在客户端，我们很多地方会用到curUser，我们将其作为全局变量
type CurUser struct {
	Conn net.Conn
	message.User
}
