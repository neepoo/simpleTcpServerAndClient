package process2

import "fmt"

//UserMgr该实例在server端有且只有一个
// 因此将其定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成对userMgr的工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对onlineUsers的添加
func (this *UserMgr) AddOnLineUsers(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 完成对onlineUsers的删除
func (this *UserMgr) DeleteOnLineUsers(userId int) {
	delete(this.onlineUsers, userId)
}

// 返回当前所有在线的用户
func (this *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

// 根据id返回对应的UserProcess
func (this *UserMgr) GetOnlinesUserById(userId int) (up *UserProcess, err error) {

	up, ok := this.onlineUsers[userId]
	if !ok {
		// 说明想找的这个用户当前不在线
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
