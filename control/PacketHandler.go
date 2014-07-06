package control

// 不允许被其他文件import

import (
	."session"
	."share"
	"errors"
)

var PacketHandler map[int]func(*Packet) (string,error)

func InitPacketHandler() {

	PacketHandler = make(map[int]func(*Packet) (string,error))

	PacketHandler[PACKET_TYPE_LOGIN] = OnLogin
	PacketHandler[PACKET_TYPE_LOGOUT] = OnLogout
	PacketHandler[PACKET_TYPE_CREATE_ROOM] = OnCreateRoom
	PacketHandler[PACKET_TYPE_JOIN_ROOM] = OnJoinRoom
	PacketHandler[PACKET_TYPE_LEAVE_ROOM] = OnLeaveRoom
	PacketHandler[PACKET_TYPE_SINGLE_CHAT] = OnSingleChat
	PacketHandler[PACKET_TYPE_MULTI_CHAT] = OnMultiChat
	PacketHandler[PACKET_TYPE_KICK] = OnKick

	PacketHandler[PACKET_TYPE_SERVER_REGISTER] = OnServerRegister
	PacketHandler[PACKET_TYPE_SERVER_REGISTER_BACK] = OnServerRegisterBack
	PacketHandler[PACKET_TYPE_CLIENT_REGISTER] = OnClientRegister
	PacketHandler[PACKET_TYPE_SERVER_SYNC_INFO] = OnServerSyncInfo
}


func HandlePacket(packet *Packet) (string,error) {
	if packet == nil {
		PrintStack()
		return "",errors.New("packet == nil")
	}
	handler := PacketHandler[int(packet.GetType())]
	if handler == nil {
		PrintStack()
		return "",errors.New("handler == nil")
	}
	return handler(packet)
}


