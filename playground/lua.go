package main

import (
	"github.com/Shopify/go-lua"
)

func RunLua(code string) (interface{}, error) {
	var err error
	var result interface{}

	l := lua.NewState()
	err = lua.DoString(l, code)
	if err != nil {
		return nil, err
	}

	l.Global("result")
	result = l.ToValue(-1)

	return result, nil
}
