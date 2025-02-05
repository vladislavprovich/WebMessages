package server

import (
	"net/http"
)

// HomeHandler обслуговує головну сторінку
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}
