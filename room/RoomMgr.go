package room

import (
	."share"
	"log"
	."router"
)

type RoomMgr struct {
	index_ int
	rooms_ map[int]*Room
}

func (this *RoomMgr) CreateRoom(user_name string, room_name string) *Room {
	room_size := len(this.rooms_)
	if room_size == 0 {
		this.rooms_ = make(map[int]*Room, G_Config.RoomsLimit)
	} else if room_size >= G_Config.RoomsLimit {
		log.Println("房间数超出服务器配置上限，创建房间失败")
		return nil
	}

	var room = new(Room)
	room.id_ = this.index_
	room.creator_ = user_name
	room.name_ = room_name

	this.rooms_[room.id_] = room

	this.index_++

	log.Println("用户", user_name, "创建房间", room_name, room.id_)

	G_Router.PrepareSyncServerInfo()

	return room
}

func (this *RoomMgr) DeleteRoom(room *Room) {
	if room == nil {
		return
	}
	log.Println("删除房间", room.id_)

	delete(this.rooms_, room.id_)
}

func (this *RoomMgr) GetRoom(id int) *Room {
	return this.rooms_[id]
}

func (this *RoomMgr) GetRooms() map[int]*Room {
	return this.rooms_
}

var G_RoomMgr RoomMgr
