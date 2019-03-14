package model

import (
	"chapterRoom/common/message"
	_ "chapterRoom/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// 服务器启动时，初始化一个userDao实例
// 做成全局变量，在需要和redis操作时，直接使用就可以
var (
	MyUserDao *UserDao
)

// 定义一个UserDao结构体
// 完成对User结构体的各种操作
// 操作redis
type UserDao struct {
	// 直接把连接池作为字段赋给它
	pool *redis.Pool
}

//使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// UserDao应该有的方法
// 1.依据用户id返回一个User实例+err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过给定的id去redis查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROE_USER_NOTEXISTS
		}
		return
	}
	// 这里需要把res 反序列化为User实例
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("有用户，但是 反序列化时失败了, err=", err)
		return
	}
	return
}

// 完成登陆的校验Login
//1. Login完成对用户的验证
//2. 如果用户的id和pwd都正确，则返回一个user实例
//3. 如果其一有错，返回对应的错误信息

func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// 先从UserDao的连接池取一个连接
	conn := this.pool.Get()
	defer conn.Close()

	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// 这时证明得到了用户，但是密码不确定
	if user.UserPwd != userPwd {
		err = ERROE_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	// 先从UserDao的连接池取一个连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		fmt.Println("该用户存在，err=", err)
		err = ERROE_USER_EXISTS
		return
	}
	fmt.Println("user=none? ", user, " ww")
	// 这时说明该用户还没注册过，则可以注册
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	// 入库
	_, err = conn.Do("Hset", "users", fmt.Sprintf("%d", user.UserId), string(data))

	if err != nil {
		fmt.Println("保存注册用户出错,err=", err)
		return
	}
	return
}
