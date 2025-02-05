package main

import (
	"messenger/internal/logger"
	"messenger/internal/server"
)

func main() {
	logger.InitLogger()
	server.InitWebSocket()
	server.StartServer()
}
