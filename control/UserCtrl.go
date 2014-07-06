package control

import (
	."session"
	."room"
	"errors"
	"strconv"
	"log"
	"strings"
)

func OnSingleChat(packet *Packet) (string, error) {
	type Msg struct {
		UserName string
		Data     string
	}
	var msg Msg
	arr := strings.Split(packet.Data_, ":")
	msg.UserName = arr[0]
	msg.Data = arr[1]

	user := G_UserMgr.GetUserBySID(packet.SID())
	if user == nil {
		return "", errors.New("用户未登录")
	}

	if user.GetMyRoom() == nil {
		return "", errors.New("用户未曾进入房间")
	}

	be_chatted_user := G_UserMgr.GetUserByName(msg.UserName)
	if be_chatted_user == nil {
		return "", errors.New("被聊天用户未登录")
	}

	if be_chatted_user.GetMyRoom() == nil || be_chatted_user.GetMyRoom() != user.GetMyRoom() {
		return "", errors.New("被聊天用户房间信息不匹配")
	}

	room_id, _, _ := user.GetMyRoom().GetInfo()
	room := G_RoomMgr.GetRoom(room_id)
	if room == nil {
		return "", errors.New("房间不存在"+strconv.Itoa(room_id))
	}

	log.Println("用户发送：", msg.Data, "给用户：", msg.UserName)

	user.SingleChat(be_chatted_user, msg.Data)
	return "",nil
}

func OnMultiChat(packet *Packet) (string, error) {
	type Msg struct {
		UserName string
		Data     string
	}
	var msg Msg
	arr := strings.Split(packet.Data_, ":")
	msg.UserName = arr[0]
	msg.Data = arr[1]

	user := G_UserMgr.GetUserBySID(packet.SID())
	if user == nil {
		return "", errors.New("用户未登录")
	}

	if user.GetMyRoom() == nil {
		return "", errors.New("用户未曾进入房间")
	}

	if msg.UserName != "all" {
		return "", errors.New("非法")
	}

	room_id, _, _ := user.GetMyRoom().GetInfo()
	room := G_RoomMgr.GetRoom(room_id)
	if room == nil {
		return "", errors.New("房间不存在"+strconv.Itoa(room_id))
	}

	log.Println("用户发送：", msg.Data, "给用户：", msg.UserName)

	user.MultiChat(msg.Data)
	return "",nil
}

func OnKick(packet *Packet) (string, error) {

	arr := strings.Split(packet.Data_, ":")
	cmd := arr[0]
	be_kicked_user_name := arr[1]

	if cmd != "kick" {
		return "", errors.New("指令非法")
	}

	user := G_UserMgr.GetUserBySID(packet.SID())
	if user == nil {
		return "", errors.New("用户未登录")
	}

	if user.GetMyRoom() == nil {
		return "", errors.New("用户未曾进入房间")
	}

	be_kicked_user := G_UserMgr.GetUserByName(be_kicked_user_name)
	if be_kicked_user == nil {
		return "", errors.New("被踢用户未登录"+be_kicked_user_name)
	}

	if be_kicked_user.GetMyRoom() == nil || be_kicked_user.GetMyRoom() != user.GetMyRoom() {
		return "", errors.New("被踢用户房间信息不匹配")
	}

	room_id, _, _ := user.GetMyRoom().GetInfo()
	room := G_RoomMgr.GetRoom(room_id)
	if room == nil {
		return "", errors.New("房间不存在"+strconv.Itoa(room_id))
	}

	log.Println("用户：", user.GetName(), "踢掉用户：", be_kicked_user_name)

	user.Kick(be_kicked_user_name)
	return "", nil
}
