package main

import (
	"messenger/internal/server"
	"messenger/internal/utils"
)

func main() {
	utils.InitLogger()
	server.InitWebSocket()
	server.StartServer()
}
