package router

import (
	."session"
	."share"
	"log"
	"strconv"
)

type Server struct {
	// server_addr_格式 127.0.0.1:6000
	Server_addr_ string
	Room_count_  int
	User_count_  int

	Session_ *Session
}

type Router struct {
	servers_ map[string]*Server
	session_ *Session

	master_addr_ string
}

func (this *Router) GetServers() map[string]*Server {
	return this.servers_
}

func (this *Router) startAsClient() bool {
	session, err := Connect(G_Config.RouterIp, G_Config.RouterPort)
	if session != nil {
		this.session_ = session

		var server Server
		server.Server_addr_ = G_Config.BindIp+":"+strconv.Itoa(G_Config.ListenPort)
		server.Room_count_ = 0
		server.User_count_ = 0

		var packet Packet
		data, _ := JsonEncode(server)
		packet.Pack(PACKET_TYPE_SERVER_REGISTER, data)
		this.send(&packet)

		return true
	}
	if err != nil {
		log.Println("Connect Error = ", err)
	}
	return false
}

func (this *Router) startAsServer() bool {
	log.Println("路由服务器启动")
	StartNetwork(G_Config.RouterIp, G_Config.RouterPort)
	return true
}

func (this *Router) send(packet *Packet) {
	if this.session_ != nil {
		this.session_.Send(packet)
	}
}

func (this *Router) AddServer(server *Server) bool {
	if server == nil {
		return false
	}
	this.servers_[server.Server_addr_] = server
	log.Println("add server", server)

	return true
}

func (this *Router) RemoveServer(server *Server) bool {
	if server == nil {
		return false
	}
	delete(this.servers_, server.Server_addr_)
	log.Println("remove server", server)

	return true
}

func (this *Router) PrepareSyncServerInfo() bool {
	var packet Packet
	data, _ := JsonEncode(this)
	packet.Pack(PACKET_TYPE_SERVER_SYNC_INFO, data)
	this.send(&packet)

	return true
}

func (this *Router) SyncServerInfo(server *Server) bool {
	if server == nil {
		return false
	}
	this.servers_[server.Server_addr_] = server
	log.Println("upsert server", server)

	return true
}

// 做个简单的LVS
func (this *Router) GetOneServer() *Server {
	min := G_Config.RoomsLimit
	var target *Server
	for i := range this.servers_ {
		server := this.servers_[i]
		if server == nil {
			continue
		}
		if server.Room_count_ < min {
			min = server.Room_count_
			target = server
		}
	}
	return target
}

var G_Router Router

func StartRouter() bool {
	G_Router.servers_ = make(map[string]*Server)
	G_Router.master_addr_ = G_Config.RouterIp+":"+strconv.Itoa(G_Config.RouterPort)
	if G_Router.startAsClient() {
		return true
	}
	return G_Router.startAsServer()
}

// 服务端网络连接断开
func OnServerDisconnect(sid int) {
	for i := range G_Router.servers_ {
		server := G_Router.servers_[i]
		if server == nil {
			continue
		}
		if server.Session_.ID() == sid {
			G_Router.RemoveServer(server)
			log.Println("服务器", sid, "断开连接")

			if server.Server_addr_ == G_Router.master_addr_ {
				log.Println("路由服务器宕机，尝试重连。")
				go func() {
					for {
						if G_Router.startAsClient() {
							log.Println("重连路由服务器成功。")
							break
						}
					}
				}()
			}

			break
		}
	}
	G_Router.RemoveServer(nil)
}
