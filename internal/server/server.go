package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer запускає HTTP сервер
func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/ws", WebSocketHandler)

	fs := http.FileServer(http.Dir("web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
