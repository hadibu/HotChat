package session

import (
	."share"
	"log"
	"fmt"
	"net"
	"io"
	_"strings"
	"strconv"
	_"encoding/json"
	"time"
)

var HandlePacketFunc func(*Packet) (string, error)
var ClientDisconnectFunc func(int)
var ServerDisconnectFunc func(int)

func Connect(ip string, port int) (*Session , error) {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	var session = new(Session)
	session.id_ = -1
	session.conn_ = (conn)
	session.send_bytes_ = 0
	session.recv_bytes_ = 0
	session.write_buffer_ = make(chan *Packet, G_Config.WritePacketCountMax)

	go handleSession(session)

	return session, nil
}

func StartNetwork(bindIp string, listenPort int) {
	defer RecoverPanicStack()

	log.Println("listen port = ", listenPort)

	s := fmt.Sprintf(":%d", listenPort)
	addr, err := net.ResolveTCPAddr("tcp", bindIp+s)
	listener, err := net.ListenTCP("tcp", addr)
	if HasError(err) {
		return
	}

	var sid int
	sid = 0

	for {
		conn, err := listener.AcceptTCP()
		if HasError(err) {
			continue
		}

		var session Session
		session.id_ = sid
		session.conn_ = conn
		session.send_bytes_ = 0
		session.recv_bytes_ = 0
		session.write_buffer_ = make(chan *Packet, G_Config.WritePacketCountMax)

		conn.SetNoDelay(false)
		conn.SetKeepAlive(true)
		conn.SetLinger(0)

		sid++

		go handleSession(&session)
	}
}

func handleSession(session *Session) {
	log.Println("accept client = ", session.RemoteAddr())

	go func() {
		session.connect_countdown_ = time.NewTimer(time.Duration(G_Config.ConnectTimeout)*time.Second)//

		<-session.connect_countdown_.C

		closeSession(session)
	}()

	go readRoutine(session)
	go writeRoutine(session)
}

func readRoutine(session *Session) {
	defer RecoverPanicStack()
	defer closeSession(session)

	for {
		if session == nil || session.conn_ == nil {
			break
		}

		sizeChunk := make([]byte, PACKET_SIZE_BUF_LEN)
		headerSize, err := io.ReadFull(session.conn_, sizeChunk)
		session.recv_bytes_ += headerSize
		if headerSize == 0 || headerSize != PACKET_SIZE_BUF_LEN || err == io.EOF {
			break
		}

		declareSize, _ := strconv.Atoi(string(sizeChunk))
		dataChunk := make([]byte, declareSize)
		readSize, err := io.ReadFull((session.conn_), dataChunk)
		session.recv_bytes_ += readSize
		if readSize == 0 || readSize != declareSize || err == io.EOF {
			break
		}
		if err != nil {
			log.Println("recv data err = ", err, " client = ", session.RemoteAddr())
			break
		}
		dataBuf := string(dataChunk)

		var packet Packet
		//log.Println("recv dataChunk:", dataBuf)
		JsonDecode(dataBuf, &packet)
		packet.Session_ = session

		log.Println("recv packet:", packet)

		result, err := HandlePacketFunc(&packet)
		if err != nil {
			log.Println("HandlePacket err = ", err)
		} else {
			//log.Println("send ", result)
			if len(result) > 0 {
				packet.Pack(packet.Type_+1, result)
				session.Send(&packet)
			}
		}

		if session.connect_countdown_ != nil {
			session.connect_countdown_.Stop()
			session.connect_countdown_ = nil
		}
	}
}

func writeRoutine(session *Session) {
	defer RecoverPanicStack()
	defer closeSession(session)

	for {
		packet, ok := <-session.write_buffer_
		if !ok || packet == nil {
			log.Println("writeRoutine err")
			break
		}
		if session == nil {
			break
		}
		if session.conn_ == nil {
			break
		}

		writeSizeStr := fmt.Sprintf("%06d", packet.Size())
		write_size, err := io.WriteString(session.conn_, writeSizeStr+packet.ToString())
		if err != nil {
			log.Println("WriteString err", err, " client = ", session.RemoteAddr())
			break
		}
		session.send_bytes_ += write_size
		//log.Println("发送消息", packet)
	}
}

func closeSession(session *Session) {
	if session != nil {
		if session.conn_ != nil {
			log.Println("break session = ", session.RemoteAddr())
		}
		session.Close()
	}
}
