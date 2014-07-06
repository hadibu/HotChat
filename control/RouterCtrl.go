package control

import (
	."session"
	"errors"
	."share"
	."router"
	"log"
)

func OnServerRegister(packet *Packet) (string, error) {

	var server = new(Server)
	JsonDecode(packet.Data_, server)
	server.Session_ = packet.GetSession()
	G_Router.AddServer(server)

	result := ""
	str, _ := JsonEncode(G_Router.GetServers())
	result += str+""

	return result, nil
}

func OnServerRegisterBack(packet *Packet) (string, error) {

	var servers map[string]*Server
	JsonDecode(packet.Data_, &servers)
	for i := range servers {
		G_Router.AddServer(servers[i])
		log.Println("notify add server", servers[i])
	}

	return "", nil
}

func OnClientRegister(packet *Packet) (string, error) {

	server := G_Router.GetOneServer()
	if server == nil {
		err := "抱歉，当前整组服务器处于不可用状态"
		return "", errors.New(err)
	}

	return server.Server_addr_, nil
}

func OnServerSyncInfo(packet *Packet) (string, error) {

	var server = new(Server)
	JsonDecode(packet.Data_, server)
	G_Router.SyncServerInfo(server)

	return "", nil
}
