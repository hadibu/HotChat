package room

import (
	"log"
	."share"
	"encoding/json"
	"time"
	"sync"
)

type Room struct {
	id_      int
	name_    string
	users_ map[string]*User
	creator_ string
	admin_   string
	close_countdown_ *time.Timer

	lock_ sync.Mutex
}

func (this *Room) GetInfo() (int, string, int ) {
	return this.id_, this.name_, len(this.users_)
}

func (this *Room) GetUsers() map[string]*User {
	return this.users_
}

func (this *Room) ToString() string {
	type _Room struct {
		Id_      int
		Name_    string
	}
	var room _Room
	room.Id_ = this.id_
	room.Name_ = this.name_
	b, err := json.Marshal(&room)
	if err != nil {
		log.Println("error:", err)
	}
	return string(b)
}

func (this *Room) CreateUser(user_name string) User {
	var user User
	user.name_ = user_name
	return user
}

func (this *Room) AddUser(user *User) bool {
	if user == nil {
		return false
	}

	this.lock_.Lock()
	user_count := len(this.users_)
	if user_count >= G_Config.UsersLimit {
		log.Println("用户登录失败，房间用户达到上限", user.name_)
		this.lock_.Unlock()
		return false
	}

	user.room_ = this
	if user_count == 0 {
		this.users_ = make(map[string]*User, G_Config.UsersLimit)
	}
	this.users_[user.name_] = user
	this.lock_.Unlock()

	if this.admin_ == "" {
		this.admin_ = user.name_
		log.Println("用户", user.name_, "被成为管理员")
	}

	if this.close_countdown_ != nil {
		this.close_countdown_.Stop()
		this.close_countdown_ = nil
		log.Println("取消房间关闭倒计时")
	}

	log.Println("用户", user.name_, "进入房间", this.id_)

	return true
}

func (this *Room) RemoveUser(user *User) {
	log.Println("用户", user.name_, "离开房间", this.id_)

	if this.admin_ == user.name_ {
		log.Println("管理员", user.name_, "离开房间")
	}

	this.lock_.Lock()
	delete(this.users_, user.name_)
	this.lock_.Unlock()

	if len(this.users_) == 0 {
		go func() {
			this.close_countdown_ = time.NewTimer(time.Duration(G_Config.RoomKeepTime)*time.Second)

			<-this.close_countdown_.C

			G_RoomMgr.DeleteRoom(this)
		}()
	}
}

func (this *Room) GetUser(user_name string) *User {
	return this.users_[user_name]
}
