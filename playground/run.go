package main

import (
	"github.com/tiaotiao/web"
)

func RunScript(c *web.Context) interface{} {
	var result interface{}
	var err error

	args := struct {
		Lang string `web:"lang,required"`
		Code string `web:"code,required"`
	}{}
	if err = c.Scheme(&args); err != nil {
		return err
	}

	switch args.Lang {
	case "js":
		result, err = RunJavascript(args.Code)
	case "lua":
		result, err = RunLua(args.Code)
	default:
		return web.NewError("language not support", web.StatusBadRequest)
	}

	if err != nil {
		return err
	}

	return result
}
