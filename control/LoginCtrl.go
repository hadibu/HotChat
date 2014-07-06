package control

import (
	."session"
	"log"
	."room"
)

func OnLogin(packet *Packet) (string, error) {
	name := packet.Data_
	log.Println("user ", name, " login")

	user := G_UserMgr.CreateUser(name,packet.SID())
	user.SetSession(packet.GetSession())

	return user.Login()


}

func OnLogout(packet *Packet) (string, error) {
	log.Println("OnLogout packet:", packet)
	return "",nil
}

