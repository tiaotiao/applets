package main

import (
	"github.com/robertkrimen/otto"
)

func RunJavascript(code string) (interface{}, error) {
	vm := otto.New()

	result, err := vm.Run(code)
	if err != nil {
		return nil, err
	}

	result, err = vm.Get("result")
	if err != nil {
		return nil, err
	}

	return result, nil
}
