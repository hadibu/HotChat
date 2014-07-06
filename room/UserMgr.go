package room

import (
	"log"
	"sync"
)

type UserMgr struct {
	user_name_map_ map[string]*User
	user_sid_map_ map[int]*User
	lock_ sync.Mutex
}

func (this *UserMgr) CreateUser(user_name string, sid int) *User {
	user := this.user_name_map_[user_name]
	if user != nil {
		log.Println("用户", user_name, "已登录，顶掉该用户")
		user.Logout()
	}

	user = new(User)
	user.name_ = user_name
	user.sid_ = sid

	if len(this.user_name_map_) == 0 {
		this.user_name_map_ = make(map[string]*User)
	}
	if len(this.user_sid_map_) == 0 {
		this.user_sid_map_ = make(map[int]*User)
	}

	this.lock_.Lock()
	this.user_name_map_[user_name] = user
	this.user_sid_map_[sid] = user
	this.lock_.Unlock()

	log.Println("创建用户", user_name)
	return user
}

func (this *UserMgr) DeleteUser(user *User) {
	this.lock_.Lock()
	delete(this.user_name_map_, user.name_)
	delete(this.user_sid_map_, user.sid_)
	this.lock_.Unlock()

	log.Println("删除用户", user.name_)
}

func (this *UserMgr) GetUserByName(user_name string) *User {
	return this.user_name_map_[user_name]
}

func (this *UserMgr) GetUserBySID(sid int) *User {
	return this.user_sid_map_[sid]
}

// 客户端网络连接断开
func OnClientDisconnect(sid int) {

	user := G_UserMgr.GetUserBySID(sid)
	if user == nil {
		return
	}
	user.Logout()
	log.Println("用户", sid, "断开连接")
}

var G_UserMgr UserMgr
