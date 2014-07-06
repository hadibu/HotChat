package control

import (
	."session"
	."room"
	"errors"
	"strconv"
)

func OnCreateRoom(packet *Packet) (string, error) {
	user := G_UserMgr.GetUserBySID(packet.SID())
	if user == nil {
		return "", errors.New("创建房间失败，用户未登录")
	}

	room_name := packet.Data_

	room := G_RoomMgr.CreateRoom(user.GetName(), room_name)

	if room != nil {
		return room.ToString(), nil
	} else {
		return "服务器房间爆满，请稍后再试", nil
	}
}

func OnJoinRoom(packet *Packet) (string, error) {
	user := G_UserMgr.GetUserBySID(packet.SID())
	if user == nil {
		return "", errors.New("加入房间失败，用户未登录")
	}

	room_id, _ := strconv.Atoi(packet.Data_)
	room := G_RoomMgr.GetRoom(room_id)
	if room == nil {
		return "", errors.New("加入房间失败，房间不存在"+strconv.Itoa(room_id))
	}

	if !room.AddUser(user) {
		return "进入房间失败，或许房间人数已满？", nil
	}

	users := room.GetUsers()

	result := ""
	for i := range users {
		result = result+users[i].GetName()+"\n"
	}

	return result, nil
}

func OnLeaveRoom(packet *Packet) (string, error) {
	user := G_UserMgr.GetUserBySID(packet.SID())
	if user == nil {
		return "", errors.New("用户未登录")
	}

	if user.GetMyRoom() == nil {
		return "", errors.New("用户未曾进入房间")
	}

	room_id, _,_ := user.GetMyRoom().GetInfo()
	room := G_RoomMgr.GetRoom(room_id)
	if room == nil {
		return "", errors.New("房间不存在"+strconv.Itoa(room_id))
	}

	room.RemoveUser(user)
	
	return "", nil
}

