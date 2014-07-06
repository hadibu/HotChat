package room

import (
	"log"
	."session"
	."share"
)

type User struct {
	sid_  int
name_ string
	room_ *Room
	session_ *Session
}

func (this *User) SetName(val string) {
	this.name_ = val
}

func (this *User) GetName() string {
	return this.name_
}

func (this *User) GetMyRoom() *Room {
	return this.room_
}

func (this *User) SetSession(val *Session) {
	this.session_ = val
}

func (this *User) SendMsg(packet *Packet) {
	if this.session_ != nil {
		this.session_.Send(packet)
	}
}

func (this *User) Login() (string, error) {
	result := make([]string, 0)
	rooms := G_RoomMgr.GetRooms()
	if len(rooms) > 0 {
		index := 1
		for i := range rooms {
			result = append(result, rooms[i].ToString())
			index++
		}
	} else {
		result = append(result, "")
	}

	log.Println("用户", this.name_, "登入")

	return JsonEncode(result)
}

func (this *User) Logout() {
	if this.session_ == nil {
		return
	}

	if this.room_ != nil {
		this.room_.RemoveUser(this)
	}
	G_UserMgr.DeleteUser(this)

	this.session_.Close()
	this.session_ = nil

	log.Println("用户", this.name_, "登出")
}


func (this *User) CreateRoom(room_name string) {
	room_id := G_RoomMgr.CreateRoom(this.name_, room_name)
	log.Println("用户", this.name_, "创建房间", room_id)
}

func (this *User) JoinRoom(room_id int) {
	room := G_RoomMgr.GetRoom(room_id)
	room.AddUser(this)
	log.Println("用户", this.name_, "加入房间", room_id)
}

func (this *User) LeaveRoom() {}

func (this *User) SingleChat(user *User, msg string) {
	if user == nil {
		return
	}
	var packet Packet
	packet.Pack(PACKET_TYPE_SINGLE_CHAT+1, this.name_+":"+msg)
	user.SendMsg(&packet)
}

func (this *User) MultiChat(msg string) {
	if this.room_ == nil {
		return
	}
	var packet Packet
	packet.Pack(PACKET_TYPE_MULTI_CHAT+1, this.name_+":"+msg)
	users := this.room_.GetUsers()
	for i := range users {
		user := users[i]
		if user == nil {
			continue
		}
		user.SendMsg(&packet)
	}
}

func (this *User) Kick(name string) {
	user := G_UserMgr.GetUserByName(name)
	this.room_.RemoveUser(user)
}
