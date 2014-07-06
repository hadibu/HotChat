package session

import (
	"net"
	."share"
	"log"
	"time"
)

type Session struct {
	id_             int
	recv_bytes_     int
	send_bytes_     int
	conn_           net.Conn
	write_buffer_   chan *Packet

	connect_countdown_ *time.Timer
}

func (this *Session) ID() int {
	return this.id_
}

func (this *Session) RemoteAddr() net.Addr {

	conn, ok := this.conn_.(*net.TCPConn)
	if ok {
		return (*((*net.TCPConn)(conn))).RemoteAddr()
	}

	var dummy  net.Addr
	return dummy
}

func (this *Session) Send(packet *Packet) {
	if packet == nil {
		return
	}
	if len(this.write_buffer_) >= G_Config.WritePacketCountMax {
		log.Println("消息发送失败，超出允许上限", this.id_)
	}

	this.write_buffer_ <- packet
}

func (this *Session) Close() {
	if this.conn_ != nil {
		this.conn_.Close()
		this.conn_ = nil
		ClientDisconnectFunc(this.id_)
		ServerDisconnectFunc(this.id_)
	}
}

