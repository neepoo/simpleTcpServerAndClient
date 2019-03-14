package message

// 消息类型常量
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// 这里定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息的具体数据
}

// 序列化之后，放到message结构体的Data字段
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户id
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

// server返回的
type LoginResMes struct {
	Code  int    `json:"code"`  // 返回状态码
	Error string `json:"error"` // 返回错误信息
	// 还少一个描述在线用户的字段
	UsersId []int // 保存用户id的切片
}

type RegisterMes struct {
	User User `json:"user"` // 类型就是User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回状态码,400表示该用户已经占用，200表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// 为了配合服务器端推送用户状态变化的类型
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` // 用户Id
	Status int `json:"status"` // 用户状态
}

// 增加一个发送消息的SmsMessage
type SmsMes struct {
	Content string `json:"content"`
	User           // 匿名结构体，继承
}
