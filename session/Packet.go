package session

import (
	"hash/crc32"
	"encoding/json"
	"log"
	."share"
)

const (
	// 存储Packet内容长度的缓冲区大小
	PACKET_SIZE_BUF_LEN = 6
)

type Packet struct {
	Type_  int16
	Crc32_ uint32
	Data_  string
	Session_    *Session
}

func (this *Packet) SID() int {
	return this.Session_.id_
}

func (this *Packet) GetSession() *Session {
	return this.Session_
}

func (this *Packet) SetType(val int16) {
	this.Type_ = val
}

func (this *Packet) GetType() int16 {
	return this.Type_
}

func (this *Packet) Pack(type_ int16, val string) {
	this.Type_ = type_
	this.Data_ = val
	this.Crc32_ = crc32.ChecksumIEEE([]byte(this.Data_))
}

func (this *Packet) UnPack() (string) {
	if crc32.ChecksumIEEE([]byte(this.Data_)) != this.Crc32_ {
		log.Println("校验失败")
		return ""
	}
	return this.Data_
}

func (this *Packet) ToString() string {
	b, err := json.Marshal(this)
	if err != nil {
		log.Println("error:", err)
	}
	return string(b)
}

func (this *Packet) Size() int {
	b, err := json.Marshal(this)
	if HasError(err) {
		return 0
	}
	return len(b)
}

