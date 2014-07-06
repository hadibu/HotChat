package share

type Config struct {
	CPU                 int
	BindIp              string
	ListenPort          int
	RouterIp            string
	RouterPort          int
	ConnectTimeout      int
	LogPath             string
	RoomsLimit          int
	UsersLimit          int
	ReadPacketCountMax  int
	WritePacketCountMax int
	RoomKeepTime        int
}

var G_Config Config
