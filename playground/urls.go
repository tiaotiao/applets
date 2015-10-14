package main

import (
	"github.com/tiaotiao/web"
)

func RegisterURLs(w *web.Web) {
	w.Handle("GET", "/", IndexPage)
	w.Handle("POST", "/api/run", RunScript)
}
