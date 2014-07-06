package main

import (
	."control"
	"log"
	"os"
	."room"
	"encoding/json"
	"strings"
	."share"
	."session"
	."router"
	"runtime"
)

func parseConfig(cfgPath string) bool {
	file, err := os.Open(cfgPath)
	if err != nil {
		log.Println("parseConfig err = ", err)
		return false
	}

	defer file.Close()

	bytes := make([]byte, 4096)
	size, err := file.Read(bytes)
	if err != nil || size >= 4096 {
		println("parseConfig err = ", err)
		return false
	}
	data := string(bytes)
	dec := json.NewDecoder(strings.NewReader(data))
	err = dec.Decode(&G_Config)
	if err != nil {
		log.Fatal("parse", data, "\nerr =", err)
		return false
	}
	log.Println("config = ", G_Config)
	return true
}

func initCallback() {
	HandlePacketFunc = HandlePacket
	ClientDisconnectFunc = OnClientDisconnect
	ServerDisconnectFunc = OnServerDisconnect
}

func main() {
	initCallback()
	parseConfig(os.Args[1])
	runtime.GOMAXPROCS(G_Config.CPU)
	InitPacketHandler()


	StartRouter()

	StartNetwork(G_Config.BindIp, G_Config.ListenPort)
}
