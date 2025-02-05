package logger

import "log"

func InitLogger() {
	log.Println("Logger initialized.")
}

func LogInfo(message string) {
	log.Println("[INFO]:", message)
}

func LogError(message string) {
	log.Println("[ERROR]:", message)
}
