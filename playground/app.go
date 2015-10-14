package main

import (
	"github.com/tiaotiao/web"
)

type App struct {
	w *web.Web
}

func NewApp() *App {
	a := new(App)
	a.w = web.NewWeb()
	return a
}

func (a *App) Run() error {
	var err error

	RegisterURLs(a.w)

	println("Server Started.")
	err = a.w.ListenAndServe("tcp", ":80")
	return err
}
